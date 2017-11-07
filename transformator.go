package main

import (
	"bytes"
	"text/template"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"os"
	"strings"
	"io"
	"log"
	"net/http"
	"path"
)

const TEMPLATE_FILE_EXT string = ".tmpl"
const TEMPLATE_PKGNAME string = "PKGNAME"

var funcMap template.FuncMap = template.FuncMap{
	"ToUpper":   strings.ToUpper,
	"ToLower":   strings.ToLower,
	"HasPrefix": strings.HasPrefix,
}

//ohejbak pro sablonu. potrebuju promenne
func (s *Service) SetTemplateVariable(key string, val interface{}) string {
	s.TemplateVariableMap[key] = val
	return "" //musi neco vracet aby to slo ze sablony zavolat
}
func (s *Service) GetTemplateVariable(key string) interface{} {
	return s.TemplateVariableMap[key]
}

func TransformFile(templName string, tmplFilename string, destFilename string, data *Service) (error) {
	var err error
	//ohejbak pro sablonu. potrebuju promenne
	//vynuluje pomocne uloziste "promenych"
	data.TemplateVariableMap = map[string]interface{}{}

	//t:= template.Must(template.ParseFiles(tmplFilename)).Funcs(funcMap)

	t := template.New(templName).Funcs(funcMap)
	t, err = t.ParseFiles(tmplFilename)
	if err != nil {
		return errors.WithStack(err)
	}

	var b bytes.Buffer

	err = t.Execute(&b, data)
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(destFilename, b.Bytes(), 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func TransformReader(templName string, tmpl io.Reader, destFilename string, data *Service) (error) {
	var err error
	//ohejbak pro sablonu. potrebuju promenne
	//vynuluje pomocne uloziste "promenych"
	data.TemplateVariableMap = map[string]interface{}{}

	//t:= template.Must(template.ParseFiles(tmplFilename)).Funcs(funcMap)

	buf := new(bytes.Buffer)
	buf.ReadFrom(tmpl)
	tmplstr := buf.String()

	t := template.New(templName).Funcs(funcMap)
	t, err = t.Parse(tmplstr)
	if err != nil {
		return errors.WithStack(err)
	}

	var b bytes.Buffer

	err = t.Execute(&b, data)
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(destFilename, b.Bytes(), 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func cpFile(src, dst string) (error) {
	in, err := os.Open(src)
	if err != nil {
		return errors.WithStack(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return errors.WithStack(err)
	}
	err = out.Sync()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func cpReader(in io.Reader, dst string) (error) {
	out, err := os.Create(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return errors.WithStack(err)
	}
	err = out.Sync()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func TransformDir(skelDir string, destDir string, data *Service) (error) {
	//prevede jmeno java baliku na adresar (javax.xml.datatype -> javax/xml/datatype)
	pkgDir := filepath.Join(strings.Split(data.JavaPackage, ".")...)

	err := filepath.Walk(skelDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		relPath, err := filepath.Rel(skelDir, srcPath)
		if err != nil {
			return errors.WithStack(err)
		}

		//pro sablonu odmazu priponu
		if filepath.Ext(srcPath) == TEMPLATE_FILE_EXT {
			relPath = relPath[:len(relPath)-len(TEMPLATE_FILE_EXT)]
		}
		relPath = strings.Replace(relPath, TEMPLATE_PKGNAME, pkgDir, -1)
		destPath := filepath.Join(destDir, relPath)
		if f.IsDir() {
			//adresar, pokud neexistuje, tak ho zalozi
			log.Println("md", destPath)
			os.MkdirAll(destPath, os.ModePerm);
		} else if filepath.Ext(srcPath) == TEMPLATE_FILE_EXT {
			//sablona
			log.Println("trans", destPath)
			err = TransformFile(filepath.Base(srcPath), srcPath, destPath, data)
			if err != nil {
				return err;
			}
		} else {
			//soubor, jen se prekopiruje
			log.Println("cp", destPath)
			err = cpFile(srcPath, destPath)
			if err != nil {
				return err;
			}

		}
		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type walkFunc func(fs http.FileSystem, path string, file os.FileInfo) error

func walkVirtual(fs http.FileSystem, dir string, walkFn walkFunc) (error) {
	dirFile, err := fs.Open(dir);
	if err != nil {
		return errors.WithStack(err)
	}
	df, err := dirFile.Stat()
	if err != nil {
		return errors.WithStack(err)
	}
	//volani funkce pro adresar
	err = walkFn(fs, dir, df)
	if err != nil {
		return err
	}
	files, err := dirFile.Readdir(-1)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, f := range files {
		path := path.Join(dir, f.Name())
		if f.IsDir() {
			err = walkVirtual(fs, path, walkFn)
			if err != nil {
				return err
			}
		} else {
			//volani funkce pro soubor
			err = walkFn(fs, path, f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func TransformVirtualDir(fs http.FileSystem, skelDir string, destDir string, data *Service) (error) {
	//prevede jmeno java baliku na adresar (javax.xml.datatype -> javax/xml/datatype)
	pkgDir := filepath.Join(strings.Split(data.JavaPackage, ".")...)

	err := walkVirtual(fs, skelDir, func(fs http.FileSystem, srcPath string, f os.FileInfo) error {
		relPath, err := filepath.Rel(skelDir, srcPath)
		if err != nil {
			return errors.WithStack(err)
		}

		//pro sablonu odmazu priponu
		if filepath.Ext(srcPath) == TEMPLATE_FILE_EXT {
			relPath = relPath[:len(relPath)-len(TEMPLATE_FILE_EXT)]
		}
		relPath = strings.Replace(relPath, TEMPLATE_PKGNAME, pkgDir, -1)
		destPath := filepath.Join(destDir, relPath)
		if f.IsDir() {
			//adresar, pokud neexistuje, tak ho zalozi
			log.Println("md", destPath)
			os.MkdirAll(destPath, os.ModePerm);
		} else if filepath.Ext(srcPath) == TEMPLATE_FILE_EXT {
			//sablona
			log.Println("trans", destPath)
			file, err := fs.Open(srcPath)
			if err != nil {
				return errors.WithStack(err)
			}
			err = TransformReader(filepath.Base(srcPath), file, destPath, data)
			if err != nil {
				return err;
			}
		} else {
			//soubor, jen se prekopiruje
			log.Println("cp", destPath)
			file, err := fs.Open(srcPath)
			if err != nil {
				return errors.WithStack(err)
			}
			err = cpReader(file, destPath)
			if err != nil {
				return err;
			}

		}
		return nil
	})
	if err != nil {
		return err
	}

	/*
		err:=filepath.Walk(skelDir,
		if err!=nil{
			return errors.WithStack(err)
		}
	*/
	return nil
}
