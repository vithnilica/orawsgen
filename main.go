package main

//go:generate go-bindata -prefix ./ -pkg main -o ./bindata.go ./templates/...

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-oci8"
	"flag"
	"os"
	"log"
	"io/ioutil"
	"path/filepath"
	"github.com/elazarl/go-bindata-assetfs"
	"encoding/json"
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
var tmpl *string = flag.String("tmpl", "soap", "Šablona generovaného projektu (soap, wsa, rs, rs-swager)")
var tmplExport *bool = flag.Bool("expt", false, "Exportuje použitou šablonu")
var confExport *bool = flag.Bool("expconf", false, "Exportuje použitou konfiguraci")
var defExport *bool = flag.Bool("expdef", false, "Exportuje výchozí nastavení")
var defConfFile *string= flag.String("def", "", "Náhrada výchozího nastavení")

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

	if *logEnabled != true {
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

	data, err := OrclServiceConfig(db, *searchPkgName, *appName, *appVer, *nameSpace, *javaPackage, *javaDS, *defConfFile)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}

	destDir := filepath.Join(*dir, *appName)

	os.MkdirAll(destDir, os.ModePerm);

	if *confExport{
		fmt.Println("export konfigurace do",filepath.Join(destDir, "orawsgen-conf.json"))
		j, err := json.MarshalIndent(data,"", "  ")
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
		err=ioutil.WriteFile(filepath.Join(destDir, "orawsgen-conf.json"),j,0644)
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
	}

	if *defExport{
		fmt.Println("export výchozího nastavení do",filepath.Join(destDir, "orawsgen-default.json"))
		j, err := json.MarshalIndent(GetDefaultDataTypeMap(),"", "  ")
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
		err=ioutil.WriteFile(filepath.Join(destDir, "orawsgen-default.json"),j,0644)
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
	}

	if *tmplDir != "" {
		err = TransformDir(*tmplDir, destDir, data)
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
	} else {
		fs := assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}
		err = TransformVirtualDir(&fs, filepath.Join("templates", *tmpl), destDir, data)
		if err != nil {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
		if *tmplExport {
			err = RestoreAssets(filepath.Join(destDir, "orawsgen-templ"), filepath.Join("templates", *tmpl))
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("hotovo")
}
