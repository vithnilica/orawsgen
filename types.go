package orawsgen



type Service struct {
	MavenAppName string
	MavenAppVer string
	WsdlNameSpace       string //url: http://oracle.generated/
	WsdlAppName         string //vh_test_pkg
	JavaPackage         string //balicek testws4
	JavaDS				string
	Methods             []Method
	DataTypeMap         map[string]DataType
	StructDataTypes     []StructDataType
	ListDataTypes       []ListDataType
	ArrayDataTypes      []ArrayDataType
	TemplateVariableMap map[string]interface{}
}

type Method struct {
	//polozky nactene z databaza
	DbName string
	DbPkgOwner string
	DbPkgName string
	DbSubprogramId int64
	IsFunction       bool   //metoda je funkce
	HasInputParam    bool   //metoda ma parametry
	HasOutputParam   bool   //metoda ma vystupni parametry
	//
	ResultDataTypeId string //typ navratove hodnoty u funkce (nil kdyz nejde o funkci)
	//polozky ktere se "dopocitaji"
	JdbcStmt         string //volani procedurt/funkce vcetne pojmenovanych parametru napr. "{call vh_test_pkg.p1(:p1,:p2,:p3)}"
	//WsdlPkgName string
	WsdlMethodName   string //jmeno funkce: p1 (prevedeno na male pismena, bez podrzitek)
	JavaWSMethodName string //jmeno funkce: p1 (male pismeno na zacatku, bez podtrzitek, pismeno za podtrzitkem velke)
	JavaWSClassName  string //odvozeno od jmena funkce, jen s velkym pismenem: P1
	//parametry
	Params           []Param
}


type UnsupportedMethod struct {
	//polozky nactene z databaza
	DbName string
	DbPkgOwner string
	DbPkgName string
	DbSubprogramId int64
	//
	Reason string
}


type GenneratedConfig struct{
	Methods            []Method
	UnsupportedMethods []UnsupportedMethod
	DataTypeMap        map[string]DataType
	StructDataTypes    []StructDataType
	ListDataTypes      []ListDataType
	ArrayDataTypes     []ArrayDataType
}


type Param struct {
	//polozky nactene z databaza
	DbParamName string
	IsIn          bool
	IsOut         bool
	//
	DataTypeId    string
	//polozky ktere se "dopocitaji"
	JavaParamName string //_return
	JdbcStmtParam string
	WsdlParamName string
}






type DataType struct {
	DataTypeId            string
	JavaClass             string
	JdbcSetParamMethod    string //jmeno funkce pro predani parametru
	JdbcRegOutParamMethod string //jmeno funkce pro registraci vystupniho parametru
	JdbcGetOutParamMethod string
	JdbcGetDbObject string
	JdbcGetWsObject string
}

type StructItem struct{
	WsdlItemName string
	JavaItemName string
	ItemDataTypeId string
}


type StructDataType struct {
	DataTypeId string
	SqlTypeName string //VH_TEST_T1
	WsdlTypeName   string //jmeno typu: VhTestT1 (prevedeno na male pismena, prvni pismeno velke)
	JavaClassName  string //StructVhTestT1
	Items []StructItem
}

type ListDataType struct {
	DataTypeId string
	SqlTypeName string //VH_TEST_A1
	JavaClassName  string //ListVhTestA1 (veskutecnosti to bude pouzity jenom pro nazev metod)
	ItemDataTypeId string
}

type ArrayDataType struct {
	DataTypeId string
	SqlTypeName string //VH_TEST_A1
	WsdlTypeName   string //jmeno typu: VhTestA1 (prevedeno na male pismena, prvni pismeno velke)
	JavaClassName  string //ArrayVhTestA1
	ItemDataTypeId string
}





func GetDefaultDataTypeMap()(map[string]DataType){
	var defaultDataTypeMap = map[string]DataType{
		"default:number": {
			"default:number",
			"java.math.BigDecimal",
			"DefaultTypesMethods.setNumberParam",
			"DefaultTypesMethods.registerNumberOutParam",
			"DefaultTypesMethods.getNumberOutParam",
			"DefaultTypesMethods.getDbObjectNumber",
			"DefaultTypesMethods.getWsObjectNumber",
		},
		"default:varchar2": {
			"default:varchar2",
			"String",
			"DefaultTypesMethods.setVarchar2Param",
			"DefaultTypesMethods.registerVarchar2OutParam",
			"DefaultTypesMethods.getVarchar2OutParam",
			"DefaultTypesMethods.getDbObjectVarchar2",
			"DefaultTypesMethods.getWsObjectVarchar2",
		},
		"default:date": {
			"default:date",
			"javax.xml.datatype.XMLGregorianCalendar",
			"DefaultTypesMethods.setDateParam",
			"DefaultTypesMethods.registerDateOutParam",
			"DefaultTypesMethods.getDateOutParam",
			"DefaultTypesMethods.getDbObjectDate",
			"DefaultTypesMethods.getWsObjectDate",
		},
		"default:xmltype": {
			"default:xmltype",
			"orawsgen.XmlAny",
			"DefaultTypesMethods.setXmlTypeParam",
			"DefaultTypesMethods.registerXmlTypeOutParam",
			"DefaultTypesMethods.getXmlTypeOutParam",
			"DefaultTypesMethods.getDbObjectXmlType",
			"DefaultTypesMethods.getWsObjectXmlType",
		},
		"default:clob": {
			"default:clob",
			"String",
			"DefaultTypesMethods.setClobParam",
			"DefaultTypesMethods.registerClobOutParam",
			"DefaultTypesMethods.getClobOutParam",
			"DefaultTypesMethods.getDbObjectClob",
			"DefaultTypesMethods.getWsObjectClob",
		},
		"default:blob": {
			"default:blob",
			"byte[]",
			"DefaultTypesMethods.setBlobParam",
			"DefaultTypesMethods.registerBlobOutParam",
			"DefaultTypesMethods.getBlobOutParam",
			"DefaultTypesMethods.getDbObjectBlob",
			"DefaultTypesMethods.getWsObjectBlob",
		},
	}

	return defaultDataTypeMap
}