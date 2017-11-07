package main

import (
	"database/sql"
	"github.com/pkg/errors"
	"strings"
	"strconv"
	"log"
	"fmt"
)

type UnsupportedError struct {
	Reason string
}

func (e *UnsupportedError) Error() string {
	return e.Reason
}

//retezec prevede na male pismena, odmaze podtrzika a pismena ktera nasledovala za podtrzitkem zmeni na velka
func deUnderscore(str string, firstCap bool) (string) {
	str = strings.ToLower(str)
	var str2 string = ""
	var u bool = firstCap
	for _, c := range str {
		if (c == '_') {
			u = true
		} else {
			if u == true {
				u = false
				str2 = str2 + strings.ToUpper(string(c))
			} else {
				str2 = str2 + string(c)
			}
		}
	}
	return str2
}

func getParamDataType(dataType string, typeOwner string, typeName string, genneratedConfig *GenneratedConfig, tx *sql.Tx) (dataTypeId string, err error) {
	if dataType == "NUMBER" {
		dataTypeId = "default:number"
	} else if dataType == "VARCHAR2" {
		dataTypeId = "default:varchar2"
	} else if dataType == "DATE" {
		dataTypeId = "default:date"
	} else if dataType == "CLOB" {
		dataTypeId = "default:clob"
	} else if dataType == "BLOB" {
		dataTypeId = "default:blob"
	} else if (dataType == "OPAQUE/XMLTYPE" || dataType == "UNDEFINED") && typeName == "XMLTYPE" {
		dataTypeId = "default:xmltype"
	} else if dataType == "OBJECT" {
		dataTypeId, err = addStruct(typeOwner, typeName, genneratedConfig, tx)
		if err != nil {
			return "", err
		}
	} else if dataType == "TABLE" {
		dataTypeId, err = addArray(true, typeOwner, typeName, genneratedConfig, tx)
		if err != nil {
			return "", err
		}
	} else {
		//vrati ze je to nepodporovany
		return "", &UnsupportedError{Reason: "metoda obsahuje nepodporovany typ parametru (" + dataType + ", " + typeOwner + ", " + typeName + ")"}
	}
	return dataTypeId, nil
}

func getItemDataType(typeOwner string, typeName string, typeCode string, genneratedConfig *GenneratedConfig, tx *sql.Tx) (dataTypeId string, err error) {
	if typeName == "NUMBER" {
		dataTypeId = "default:number"
	} else if typeName == "VARCHAR2" {
		dataTypeId = "default:varchar2"
	} else if typeName == "DATE" {
		dataTypeId = "default:date"
	} else if typeName == "CLOB" {
		dataTypeId = "default:clob"
	} else if typeName == "BLOB" {
		dataTypeId = "default:blob"
	} else if typeName == "XMLTYPE" {
		dataTypeId = "default:xmltype"
	} else if typeCode == "OBJECT" {
		dataTypeId, err = addStruct(typeOwner, typeName, genneratedConfig, tx)
		if err != nil {
			return "", err
		}
	} else if typeCode == "COLLECTION" {
		dataTypeId, err = addArray(false, typeOwner, typeName, genneratedConfig, tx)
		if err != nil {
			return "", err
		}
	} else {
		//vrati ze je to nepodporovany
		return "", &UnsupportedError{Reason: "nepodporovany typ (" + typeOwner + ", " + typeName + ", " + typeCode + ")"}
	}
	return dataTypeId, nil
}

//najde skutecne jmeno baliku (zatim pocitam se stejnym schematem)
func findPkg(searchPkgName string, tx *sql.Tx) (pkgOwner string, pkgName string, err error) {
	log.Println("findPkg", searchPkgName)
	err = tx.QueryRow(`
		select o.object_name, o.owner from all_objects o
		where upper(o.object_name)=upper(:1) and o.object_type in ('PACKAGE') and o.owner=user`, searchPkgName).Scan(&pkgName, &pkgOwner)
	if err != nil {
		return "", "", errors.Wrap(err, "balik nenalezen")
	}
	return pkgOwner, pkgName, nil
}

func addPkgMethod(methodName string, pkgOwner string, pkgName string, subprogramId int64, genneratedConfig *GenneratedConfig, tx *sql.Tx) (err error) {
	log.Println("addMethod", pkgName, subprogramId)
	var method Method
	var jdbcStmtParamStr string
	var paramNum int = 0

	method.DbName = methodName;
	method.DbPkgOwner = pkgOwner
	method.DbPkgName = pkgName
	method.DbSubprogramId = subprogramId

	//TODO test na duplicitu

	//vychozi hodnoty ktere se budou menit podle parametru
	method.IsFunction = false
	method.HasInputParam = false
	method.HasOutputParam = false
	method.Params = make([]Param, 0)

	rows, err := tx.Query(`
		select aa.argument_name, aa.position, aa.data_level, aa.in_out, aa.data_type, aa.type_owner, aa.type_name, aa.type_subname, aa.type_link
		from all_arguments aa
		where aa.package_name=:1 and aa.owner=:2 and aa.subprogram_id=:3 and aa.data_level=0
		order by aa.sequence`, pkgName, pkgOwner, subprogramId)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		var param Param
		var argName sql.NullString
		var position int
		var dataLevel int
		var inOut string
		var dataType sql.NullString
		var typeOwner sql.NullString
		var typeName sql.NullString
		var typeSubname sql.NullString
		var typeLink sql.NullString
		var isFunctionRestult bool

		err = rows.Scan(&argName, &position, &dataLevel, &inOut, &dataType, &typeOwner, &typeName, &typeSubname, &typeLink)
		if err != nil {
			return errors.WithStack(err)
		}

		param.IsIn = false
		param.IsOut = false
		isFunctionRestult = false

		//test jestli nejde o "prazdny radek" pro proceduru bez parametru
		if position == 1 && dataLevel == 0 && inOut == "IN" && argName.String == "" && dataType.String == "" && typeOwner.String == "" && typeName.String == "" {
			//je to procedura bez parametru, tak to preskakuju
			continue
		}

		if position == 0 {
			//navratova hodnota funkce
			isFunctionRestult = true
			method.IsFunction = true
		} else {
			isFunctionRestult = false
			param.DbParamName = argName.String
			if inOut == "IN" {
				param.IsIn = true
				param.IsOut = false
				method.HasInputParam = true
			} else if inOut == "OUT" {
				param.IsIn = false
				param.IsOut = true
				method.HasOutputParam = true
			} else {
				param.IsIn = true
				param.IsOut = true
				method.HasInputParam = true
			}
			paramNum++
			if paramNum == 1 {
				jdbcStmtParamStr = ":p1"
			} else {
				jdbcStmtParamStr = jdbcStmtParamStr + ", :p" + strconv.Itoa(paramNum)
			}

			//param.JavaParamName = "p" + strconv.Itoa(paramNum)
			param.JavaParamName = deUnderscore(param.DbParamName, false)
			param.JdbcStmtParam = "p" + strconv.Itoa(paramNum)
			param.WsdlParamName = deUnderscore(param.DbParamName, false)

			if param.IsOut {
				if param.IsIn {
					param.WsdlOutParamName = strings.ToLower(param.WsdlParamName) + "Inout"
				} else {
					param.WsdlOutParamName = strings.ToLower(param.WsdlParamName) + "Out"
				}
			}
		}

		var dataTypeId string

		dataTypeId, err = getParamDataType(dataType.String, typeOwner.String, typeName.String, genneratedConfig, tx)
		if err != nil {
			switch e := err.(type) {
			case *UnsupportedError:
				genneratedConfig.UnsupportedMethods = append(genneratedConfig.UnsupportedMethods, UnsupportedMethod{DbName: methodName, DbPkgOwner: pkgOwner, DbPkgName: pkgName, DbSubprogramId: subprogramId, Reason: e.Reason})
				return nil
			default:
				return err
			}
		}

		if isFunctionRestult {
			method.ResultDataTypeId = dataTypeId
		} else {
			param.DataTypeId = dataTypeId
			method.Params = append(method.Params, param)
		}

	}

	if (method.IsFunction) {
		method.JdbcStmt = "{call :r1:=" + pkgOwner + "." + pkgName + "." + methodName + "(" + jdbcStmtParamStr + ")}" //volani procedurt/funkce vcetne pojmenovanych parametru napr. "{call :r1=vh_test_pkg.p1(:p1,:p2,:p3)}"
	} else {
		method.JdbcStmt = "{call " + pkgOwner + "." + pkgName + "." + methodName + "(" + jdbcStmtParamStr + ")}" //volani procedurt/funkce vcetne pojmenovanych parametru napr. "{call vh_test_pkg.p1(:p1,:p2,:p3)}"
	}

	//method.WsdlMethodName = strings.Replace(strings.ToLower(methodName), "_", "", -1) //jmeno funkce: p1 (prevedeno na male pismena, bez podrzitek)
	method.WsdlMethodName = deUnderscore(methodName, false)   //jmeno funkce: p1
	method.JavaWSMethodName = deUnderscore(methodName, false) //jmeno funkce: p1 (male pismeno na zacatku, bez podtrzitek, pismeno za podtrzitkem velke)
	method.JavaWSClassName = deUnderscore(methodName, true)   //odvozeno od jmena funkce, jen s velkym pismenem: P1
	method.WsdlPkgName = strings.ToLower(pkgName)

	//ohjebak pro outparametry "setrideny" nejdrive podle abecedy a pak pomoci hashmapy v jave
	var parr15 [][2]string
	if method.IsFunction {
		parr15 = append(parr15, [2]string{"return", "ret"})
	}
	for i := range method.Params {
		if method.Params[i].IsOut {
			if method.Params[i].IsIn {
				parr15 = append(parr15, [2]string{strings.ToLower(method.Params[i].WsdlParamName) + "Inout", method.Params[i].JavaParamName})
			} else {
				parr15 = append(parr15, [2]string{strings.ToLower(method.Params[i].WsdlParamName) + "Out", method.Params[i].JavaParamName})
			}
		}
	}
	method.JavaWSAOutParams = MixerJava15(parr15)

	genneratedConfig.Methods = append(genneratedConfig.Methods, method)
	return nil

}

//nacte vsechny metody k baliku
func addPkgMethods(pkgOwner string, pkgName string, genneratedConfig *GenneratedConfig, tx *sql.Tx) (err error) {
	log.Println("addPkgMethods", pkgName)
	rows, err := tx.Query(`
		select p.procedure_name,p.subprogram_id,p.aggregate,p.pipelined
		from all_procedures p
		where p.object_name=:1 and p.owner=:2 and p.object_type='PACKAGE' and p.procedure_name is not null
		order by p.subprogram_id`, pkgName, pkgOwner)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		var methodName string
		var subprogramId int64
		var aggregate string
		var pipelined string

		rows.Scan(&methodName, &subprogramId, &aggregate, &pipelined)

		if aggregate == "YES" {
			genneratedConfig.UnsupportedMethods = append(genneratedConfig.UnsupportedMethods, UnsupportedMethod{DbName: methodName, DbPkgOwner: pkgOwner, DbPkgName: pkgName, DbSubprogramId: subprogramId, Reason: "metoda je agregacni"})
			continue
		}
		if pipelined == "YES" {
			genneratedConfig.UnsupportedMethods = append(genneratedConfig.UnsupportedMethods, UnsupportedMethod{DbName: methodName, DbPkgOwner: pkgOwner, DbPkgName: pkgName, DbSubprogramId: subprogramId, Reason: "metoda je pipelined"})
			continue
		}

		//hlida duplicitni jmeno (nepovedlo se zjistit podle jakeho pravidla wsa generuje jmeno)
		dupl := false
		for i := range genneratedConfig.Methods {
			if (genneratedConfig.Methods[i].DbName == methodName) {
				dupl = true
				break
			}
		}
		if dupl {
			genneratedConfig.UnsupportedMethods = append(genneratedConfig.UnsupportedMethods, UnsupportedMethod{DbName: methodName, DbPkgOwner: pkgOwner, DbPkgName: pkgName, DbSubprogramId: subprogramId, Reason: "duplicitni jmeno metody"})
			continue
		}

		err = addPkgMethod(methodName, pkgOwner, pkgName, subprogramId, genneratedConfig, tx)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func addArray(isParam bool, typeOwner string, typeName string, genneratedConfig *GenneratedConfig, tx *sql.Tx) (arrayDataTypeId string, err error) {
	log.Println("addArray", typeName)
	if isParam {
		arrayDataTypeId = "generated:list:" + strings.ToLower(typeName)
	} else {
		arrayDataTypeId = "generated:array:" + strings.ToLower(typeName)
	}

	//overi si ze uz nebylo generovano, pokud ano vratit predchozi vysledek a do databaze nechodi
	_, ok := genneratedConfig.DataTypeMap[arrayDataTypeId]
	if ok {
		//TODO zkontrolovat ze to neni duplicita (stejne jmeno, jinej vlastnik)?
		return arrayDataTypeId, nil
	}

	var dataType DataType
	dataType.DataTypeId = arrayDataTypeId

	var typeName4Java string

	if isParam {
		typeName4Java = "List" + deUnderscore(typeName, true)
		//dataType.JavaClass - tady se doplni List<Typ> az bude znamo
		dataType.JdbcSetParamMethod = "GeneratedTypesMethods.set" + typeName4Java + "Param"            //jmeno funkce pro predani parametru
		dataType.JdbcRegOutParamMethod = "GeneratedTypesMethods.register" + typeName4Java + "OutParam" //jmeno funkce pro registraci vystupniho parametru
		dataType.JdbcGetOutParamMethod = "GeneratedTypesMethods.get" + typeName4Java + "OutParam"
		//dataType.JdbcGetDbObject string
		//dataType.JdbcGetWsObject string
	} else {
		typeName4Java = "Array" + deUnderscore(typeName, true)
		dataType.JavaClass = "GeneratedWSTypes." + typeName4Java
		//dataType.JdbcSetParamMethod    string //jmeno funkce pro predani parametru
		//dataType.JdbcRegOutParamMethod string //jmeno funkce pro registraci vystupniho parametru
		//dataType.JdbcGetOutParamMethod string
		dataType.JdbcGetDbObject = "GeneratedTypesMethods.getDbObject" + typeName4Java + ""
		dataType.JdbcGetWsObject = "GeneratedTypesMethods.getWsObject" + typeName4Java + ""
	}

	rows, err := tx.Query(`
		select ct.elem_type_owner, ct.elem_type_name, t.typecode from all_coll_types ct, all_types t
		where ct.coll_type='TABLE' and ct.type_name=:1 and ct.owner=:2 and ct.elem_type_owner=t.owner(+) and ct.elem_type_name=t.type_name(+) `, typeName, typeOwner)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer rows.Close()

	//test jestli typ existuje
	if !rows.Next() {
		return "", errors.New("typ nenalezen v pohledu all_coll_types " + typeOwner + ", " + typeName)
	}

	var elemTypeOwner sql.NullString
	var elemTypeName sql.NullString
	var elemTypeCode sql.NullString

	rows.Scan(&elemTypeOwner, &elemTypeName, &elemTypeCode)

	var dataTypeId string

	dataTypeId, err = getItemDataType(elemTypeOwner.String, elemTypeName.String, elemTypeCode.String, genneratedConfig, tx)
	if err != nil {
		return "", err
	}

	if isParam {
		var listDataType ListDataType
		listDataType.DataTypeId = arrayDataTypeId
		listDataType.SqlTypeName = typeOwner + "." + typeName
		listDataType.JavaClassName = typeName4Java //ListVhTestA1 (veskutecnosti to bude pouzity jenom pro nazev metod)
		listDataType.ItemDataTypeId = dataTypeId
		genneratedConfig.ListDataTypes = append(genneratedConfig.ListDataTypes, listDataType)
	} else {
		var arrayDataType ArrayDataType
		arrayDataType.DataTypeId = arrayDataTypeId
		arrayDataType.SqlTypeName = typeOwner + "." + typeName
		arrayDataType.WsdlTypeName = deUnderscore(typeName, true) //jmeno typu: VhTestA1 (prevedeno na male pismena, prvni pismeno velke)
		arrayDataType.JavaClassName = typeName4Java               //ArrayVhTestA1
		arrayDataType.ItemDataTypeId = dataTypeId
		genneratedConfig.ArrayDataTypes = append(genneratedConfig.ArrayDataTypes, arrayDataType)
	}

	genneratedConfig.DataTypeMap[arrayDataTypeId] = dataType

	return arrayDataTypeId, nil
}

func addStruct(typeOwner string, typeName string, genneratedConfig *GenneratedConfig, tx *sql.Tx) (structDataTypeId string, err error) {
	log.Println("addStruct", typeName)
	structDataTypeId = "generated:struct:" + strings.ToLower(typeName)

	//overi si ze uz nebylo generovano, pokud ano vratit predchozi vysledek a do databaze nechodi
	_, ok := genneratedConfig.DataTypeMap[structDataTypeId]
	if ok {
		//TODO zkontrolovat ze to neni duplicita (stejne jmeno, jinej vlastnik)?
		return structDataTypeId, nil
	}

	var dataType DataType
	dataType.DataTypeId = structDataTypeId

	var typeName4Java string
	typeName4Java = "Struct" + deUnderscore(typeName, true)

	dataType.JavaClass = "GeneratedWSTypes." + typeName4Java
	dataType.JdbcSetParamMethod = "GeneratedTypesMethods.set" + typeName4Java + "Param"            //jmeno funkce pro predani parametru
	dataType.JdbcRegOutParamMethod = "GeneratedTypesMethods.register" + typeName4Java + "OutParam" //jmeno funkce pro registraci vystupniho parametru
	dataType.JdbcGetOutParamMethod = "GeneratedTypesMethods.get" + typeName4Java + "OutParam"
	dataType.JdbcGetDbObject = "GeneratedTypesMethods.getDbObject" + typeName4Java + ""
	dataType.JdbcGetWsObject = "GeneratedTypesMethods.getWsObject" + typeName4Java + ""

	var structDataType StructDataType
	structDataType.DataTypeId = structDataTypeId
	structDataType.SqlTypeName = typeOwner + "." + typeName    //VH_TEST_T1
	structDataType.WsdlTypeName = deUnderscore(typeName, true) //jmeno typu: VhTestT1 (prevedeno na male pismena, prvni pismeno velke)
	structDataType.JavaClassName = typeName4Java               //StructVhTestT1
	structDataType.Items = make([]StructItem, 0)

	rows, err := tx.Query(`
		select a.attr_name, a.attr_type_owner, a.attr_type_name, t.typecode
		from all_type_attrs a, all_types t where a.type_name=:1 and a.owner=:2 and a.attr_type_owner=t.owner(+) and a.attr_type_name=t.type_name(+)
		order by a.attr_no`, typeName, typeOwner)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		var attrName sql.NullString
		var attrTypeOwner sql.NullString
		var attrTypeName sql.NullString
		var attrTypeCode sql.NullString

		rows.Scan(&attrName, &attrTypeOwner, &attrTypeName, &attrTypeCode)

		var dataTypeId string

		dataTypeId, err = getItemDataType(attrTypeOwner.String, attrTypeName.String, attrTypeCode.String, genneratedConfig, tx)
		if err != nil {
			return "", err
		}

		var item StructItem
		item.ItemDataTypeId = dataTypeId
		item.WsdlItemName = deUnderscore(attrName.String, false) //jmeno polozky prevedene na male pisemena?
		item.JavaItemName = deUnderscore(attrName.String, false) //TODO mozna predelat na p1, p2, p3 ...
		structDataType.Items = append(structDataType.Items, item)
	}

	//kontrola ze typ existuje
	if (len(structDataType.Items) == 0) {
		return "", errors.New("typ nenalezen v pohledu all_type_attrs " + typeOwner + ", " + typeName)
	}

	//ohjebak pro polozky "setrideny" nejdrive podle abecedy a pak pomoci hashmapy v jave
	var iarr15 [][2]string
	for i := range structDataType.Items {
		iarr15 = append(iarr15, [2]string{structDataType.Items[i].WsdlItemName, structDataType.Items[i].JavaItemName})
	}
	structDataType.JavaWSAItems = MixerJava15(iarr15)

	genneratedConfig.StructDataTypes = append(genneratedConfig.StructDataTypes, structDataType)
	genneratedConfig.DataTypeMap[structDataTypeId] = dataType

	return structDataTypeId, nil
}

func orclGenerateConfig(searchPkgName string, tx *sql.Tx) (genneratedConfig *GenneratedConfig, err error) {
	var pkgOwner string
	var pkgName string

	genneratedConfig = &GenneratedConfig{}
	genneratedConfig.DataTypeMap = make(map[string]DataType, 0)

	//najde skutecne jmeno baliku (zatim pocitam se stejnym schematem)
	pkgOwner, pkgName, err = findPkg(searchPkgName, tx)
	if err != nil {
		//fmt.Printf("%+v\n",err)
		return nil, err
	}

	//nacte vsechny metody k baliku
	err = addPkgMethods(pkgOwner, pkgName, genneratedConfig, tx)
	if err != nil {
		//fmt.Printf("%+v\n",err)
		return nil, err
	}

	return genneratedConfig, nil

}

//doplni javovy typ List<Typ> z toho co pri generovani nebylo znamo
func orclFinalizeGeneratedList(dataTypeMap *map[string]DataType, listDataTypes *[]ListDataType) {
	for _, listDataType := range *listDataTypes {
		dt := (*dataTypeMap)[listDataType.DataTypeId]
		dt.JavaClass = "List<" + (*dataTypeMap)[listDataType.ItemDataTypeId].JavaClass + ">"
		(*dataTypeMap)[listDataType.DataTypeId] = dt;
	}
}

func OrclServiceConfig(db *sql.DB, searchPkgName string, appName string, appVer string, nameSpace string, javaPackage string, javaDS string) (*Service, error) {
	var tx *sql.Tx
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	genneratedConfig, err := orclGenerateConfig(searchPkgName, tx)
	if err != nil {
		return nil, err
	}

	var data = Service{
		MavenAppName:     appName,
		MavenAppVer:      appVer,
		WsdlNameSpace:    nameSpace,
		WsdlAppName:      appName,
		WsdlPortTypeName: strings.ToLower(searchPkgName),
		JavaPackage:      javaPackage,
		JavaDS:           javaDS,
	}

	if data.DataTypeMap == nil {
		data.DataTypeMap = map[string]DataType{};
	}

	//pripojeni vygenerovanych hodnot
	data.Methods = append(data.Methods, genneratedConfig.Methods...)
	for k, v := range genneratedConfig.DataTypeMap {
		data.DataTypeMap[k] = v
	}
	data.StructDataTypes = append(data.StructDataTypes, genneratedConfig.StructDataTypes...)
	data.ArrayDataTypes = append(data.ArrayDataTypes, genneratedConfig.ArrayDataTypes...)
	data.ListDataTypes = append(data.ListDataTypes, genneratedConfig.ListDataTypes...)

	//pripojeni vychozich hodnot
	for k, v := range GetDefaultDataTypeMap() {
		data.DataTypeMap[k] = v
	}

	//doplneni polozek ktere pri generovani nebyly dostupne
	orclFinalizeGeneratedList(&data.DataTypeMap, &data.ListDataTypes)

	//vypise netpodporovane metody
	if len(genneratedConfig.UnsupportedMethods) > 0 {
		fmt.Println("metody ktere se nepouziji")
		for _, unsupportedMethod := range genneratedConfig.UnsupportedMethods {
			fmt.Println(unsupportedMethod.DbPkgName+"."+unsupportedMethod.DbName, "cislo metody:", unsupportedMethod.DbSubprogramId, "duvod:", unsupportedMethod.Reason)
		}
	}

	return &data, nil
}
