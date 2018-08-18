package main

import (
	"io"
	"log"
	"os"
	"text/template"

	"github.com/pdk/crudgen/crudlib"
)

// TemplateData provides data to the template.
type TemplateData struct {
	PackageName   string
	TableName     string
	Style         crudlib.BindStyle
	Select        string
	IncludeUpdate bool
	Struct
}

// PrintTemplate will write out the new .go source
func PrintTemplate(w io.Writer, packageName string, tableName string, s Struct, style crudlib.BindStyle, selectName string, noUpdate bool) {

	templatePath := os.Getenv("GOPATH") + "/src/github.com/pdk/crudgen/templates/crud.go-template"

	template, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("%s", err)
	}

	templateData := TemplateData{
		PackageName:   packageName,
		TableName:     tableName,
		Style:         style,
		Struct:        s,
		Select:        selectName,
		IncludeUpdate: !noUpdate,
	}

	err = template.Execute(w, templateData)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}
