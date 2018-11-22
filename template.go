package main

import (
	"encoding/json"
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

// NewTemplateData constructs a TemplateData for template execution or dumping.
func NewTemplateData(packageName string, tableName string, s Struct, style crudlib.BindStyle, selectName string, noUpdate bool) TemplateData {

	return TemplateData{
		PackageName:   packageName,
		TableName:     tableName,
		Style:         style,
		Struct:        s,
		Select:        selectName,
		IncludeUpdate: !noUpdate,
	}
}

// DumpJSON writes out JSON info required for custom templates.
func DumpJSON(w io.Writer, data TemplateData) {

	d := map[string]interface{}{}
	s := data.Struct
	b := data.Style

	d["PackageName"] = data.PackageName
	d["BindStyle"] = data.Style.String()
	d["TableName"] = data.TableName
	d["IncludeUpdate"] = data.IncludeUpdate
	d["StructName"] = s.Name

	d["InsertStatement"] = b.Rebind(s.InsertStatement(data.TableName))
	d["InsertBindVars"] = s.InsertBindVars()
	d["InsertColumnNames"] = s.insertColumnNames()
	d["InsertFieldNames"] = s.insertFieldNames()

	d["UpdateStatement"] = b.Rebind(s.UpdateStatement(data.TableName))
	d["UpdateBindVars"] = s.UpdateBindVars()
	d["UpdateKeyColumnNames"] = s.keyColumnNames()
	d["UpdateKeyFieldNames"] = s.keyFieldNames()
	d["UpdateValueColumnNames"] = s.valueColumnNames()
	d["UpdateValueFieldNames"] = s.valueFieldNames()

	d["DeleteStatement"] = b.Rebind(s.DeleteStatement(data.TableName))
	d["DeleteBindVars"] = s.DeleteBindVars()
	d["DeleteKeyColumnNames"] = s.keyColumnNames()
	d["DeleteKeyFieldNames"] = s.keyFieldNames()

	d["SelectStatement"] = s.SelectStatement(data.TableName)
	d["SelectColumnNames"] = s.selectColumnNames()
	d["SelectFieldNames"] = s.selectFieldNames()

	d["CreateTimestampColumnNames"] = s.createTimestampColumnNames()
	d["CreateTimestampFieldNames"] = s.createTimestampFieldNames()

	d["UpdateTimestampColumnNames"] = s.updateTimestampColumnNames()
	d["UpdateTimestampFieldNames"] = s.updateTimestampFieldNames()

	d["HasAutoIncrColumn"] = s.HasAutoIncrColumn()
	d["AutoIncrColumnName"] = s.AutoIncrColumnName()
	d["AutoIncrFieldName"] = s.AutoIncrFieldName()

	result, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		log.Fatalf("cannot make JSON: %s", err)
	}

	_, err = w.Write(result)

	if err != nil {
		log.Fatalf("cannot write result: %s", err)
	}
}

// PrintTemplate will write out the new .go source
func PrintTemplate(w io.Writer, templateData TemplateData) {

	templatePath := os.Getenv("GOPATH") + "/src/github.com/pdk/crudgen/templates/crud.go-template"

	template, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = template.Execute(w, templateData)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}
