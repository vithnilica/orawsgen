package main

import (
	"github.com/vithnilica/orawsgen"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-oci8"
)


func main() {
	//log.SetFlags(0)
	//log.SetOutput(ioutil.Discard)


	conStr := "test/password@192.168.0.20:1521/xe"
	//conStr:="customer/customer@test152"


	searchPkgName := "vh_test_pkg"
	nameSpace:="http://oracle.generated/"
	appName:="vh_test_pkg_app"
	javaPackage:="vh.test.pkg"


/*
	searchPkgName := "cz_pk2ap_ws4"
	nameSpace:="http://oracle.generated/"
	appName:="ap_ws4"
	javaPackage:="apws4"
*/

/*
	searchPkgName := "cz_ws_moa2"
	nameSpace:="http://opus_pk_moa2.ws.allianz.cz/"
	appName:="opus_pk_moa"
	javaPackage:="moa"
*/



	db, err := sql.Open("oci8", conStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//db.SetMaxOpenConns(1)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	data,err:=orawsgen.OrclServiceConfig(db, searchPkgName, appName, nameSpace, javaPackage)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}




	err=orawsgen.TransformDir("templates/wsa/","out/pokus1", data)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}

	fmt.Println("hotovo")
}
