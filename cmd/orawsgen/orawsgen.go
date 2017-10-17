package main

//go:generate go-bindata -prefix ../../ -pkg main -o ./bindata.go ../../templates/...

import (
	"github.com/vithnilica/orawsgen"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-oci8"
	"flag"
	"os"
	"log"
	"io/ioutil"
	"path/filepath"
	"github.com/elazarl/go-bindata-assetfs"
)


var conStr *string = flag.String("c", "", "Přihlašovací údaje do databáze (např. user/password@db123)")
var searchPkgName *string = flag.String("pkg", "", "Jméno balíku (např. cz_ws_moa2)")
var nameSpace *string = flag.String("ns", "http://oracle.generated/", "Namespace webové služby")
var appName *string = flag.String("app", "", "Jméno webové služby (např. opus_pk_moa2)")
var appVer *string = flag.String("appver", "1.0.0", "Verze webové služby")
var javaPackage *string = flag.String("javapkg", "generated", "Jméno balíku v javě")
var javaDS *string = flag.String("ds", "java:/OracleDS", "JNDI datového zdroje")
var logEnabled *bool = flag.Bool("log", false, "Zapne logování")
var dir *string = flag.String("dir", getwd(), "Pracovní adresář")
var tmplDir *string = flag.String("tdir", "", "Adresář s šablonou generovaného projektu")
var tmpl *string = flag.String("tmpl", "wsa", "Šablona generovaného projektu")
var tmplExport *bool = flag.Bool("texp", false, "Exportuje použitou šablonu")

func getwd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}


func main() {
	flag.Parse()
	if conStr == nil || *conStr == "" {
		fmt.Println("Nejsou vyplněny přihlašovací údaje do databáze")
		flag.Usage()
		os.Exit(1)
	}
	if searchPkgName == nil || *searchPkgName == "" {
		fmt.Println("Není vyplněno jméno balíku")
		flag.Usage()
		os.Exit(1)
	}

	if appName == nil || *appName == "" {
		fmt.Println("Není vyplněno jméno webové služby")
		flag.Usage()
		os.Exit(1)
	}

	if *logEnabled!=true{
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	db, err := sql.Open("oci8", *conStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	data,err:=orawsgen.OrclServiceConfig(db, *searchPkgName, *appName, *appVer, *nameSpace, *javaPackage, *javaDS)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}

	destDir:=filepath.Join(*dir,*appName)

	if *tmplDir!=""{
		err=orawsgen.TransformDir(*tmplDir,destDir, data)
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
	}else{
		fs:=assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}
		err=orawsgen.TransformVirtualDir(&fs,filepath.Join("templates",*tmpl),destDir, data)
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
		if *tmplExport{
			err=RestoreAssets(filepath.Join(destDir,"template"), filepath.Join("templates",*tmpl))
			if err != nil {
				panic(err)
			}
		}
	}


	fmt.Println("hotovo")
}
