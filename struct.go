package main

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pdk/crudgen/crudlib"
)

// Values used to identify particular columns for insert/update/delete
// operations.
const (
	AutoIncr        = "autoincr"
	Key             = "key"
	CreateTimestamp = "create_timestamp"
	UpdateTimestamp = "update_timestamp"
)

// Struct is a struct declaration we found in the source file.
type Struct struct {
	Name   string
	Fields []Field
	Comps  []Composition
}

// Field is what we found in a struct declaration.
type Field struct {
	Name    string
	DBTag   string
	CrudTag string
	Prefix  string
}

// Composition is another struct that is composed into the struct.
type Composition struct {
	Name     string
	Composed Struct
}

// Compose adds another struct to be composed.
func (s *Struct) Compose(name string, comp Struct) {
	s.Comps = append(s.Comps, Composition{
		Name:     name,
		Composed: comp,
	})
}

// AppendField adds a field name/tag pair to a Struct.
func (s *Struct) AppendField(name, dbTag, crudTag string) {
	s.Fields = append(s.Fields, Field{
		Name:    name,
		DBTag:   dbTag,
		CrudTag: crudTag,
	})
}

// NewStruct returns a new Struct struct, which is used to identify found
// structs.
func NewStruct(name string) *Struct {

	s := Struct{
		Name: name,
	}

	return &s
}

// AllFields recursively builds a list of all the fields in the struct.
func (s Struct) AllFields() []Field {
	all := make([]Field, 0)
	skipNames := make(map[string]bool)

	for _, c := range s.Comps {
		for _, fld := range c.Composed.AllFields() {
			fld.Prefix = c.Name + "." + fld.Prefix
			all = append(all, fld)
		}
		skipNames[c.Name] = true
	}

	for _, f := range s.Fields {
		if !skipNames[f.Name] {
			all = append(all, f)
		}
	}

	return all
}

// ColumnName returns either the DBTag value or snake-case of the Name.
func (f Field) ColumnName() string {
	if f.DBTag != "" {
		return f.DBTag
	}

	return strcase.ToSnake(f.Name)
}

// FieldName returns Prefix+Name.
func (f Field) FieldName() string {
	return f.Prefix + f.Name
}

// pickFields filters a slice of Fields
func pickFields(fields []Field, picker func(Field) bool) []Field {
	l := make([]Field, 0)
	for _, f := range fields {
		if picker(f) {
			l = append(l, f)
		}
	}

	return l
}

// columnNames gets the ColumnNames (database names) of a slice of Fields.
func columnNames(fields []Field) []string {
	v := []string{}
	for _, f := range fields {
		v = append(v, f.ColumnName())
	}

	return v
}

// fieldNames gets the FieldNames (go name) of a slice of Fields.
func fieldNames(fields []Field) []string {
	v := []string{}
	for _, f := range fields {
		v = append(v, f.FieldName())
	}

	return v
}

func isAutoIncr(f Field) bool {
	return f.CrudTag == AutoIncr
}

func notAutoIncr(f Field) bool {
	return !isAutoIncr(f)
}

func (s Struct) insertColumnNames() []string {
	return columnNames(pickFields(s.AllFields(), notAutoIncr))
}

func (s Struct) insertFieldNames() []string {
	return fieldNames(pickFields(s.AllFields(), notAutoIncr))
}

func isKey(f Field) bool {
	return f.CrudTag == Key || f.CrudTag == AutoIncr
}

func notKey(f Field) bool {
	return !isKey(f)
}

func (s Struct) valueColumnNames() []string {
	return columnNames(pickFields(s.AllFields(), notKey))
}

func (s Struct) valueFieldNames() []string {
	return fieldNames(pickFields(s.AllFields(), notKey))
}

func (s Struct) keyColumnNames() []string {
	return columnNames(pickFields(s.AllFields(), isKey))
}

func (s Struct) keyFieldNames() []string {
	return fieldNames(pickFields(s.AllFields(), isKey))
}

func (s Struct) selectColumnNames() []string {
	return columnNames(s.AllFields())
}

func (s Struct) selectFieldNames() []string {
	return fieldNames(s.AllFields())
}

// AutoIncrColumnName returns the name of the column marked with crud tag
// AutoIncr.
func (s Struct) AutoIncrColumnName() string {
	f := columnNames(pickFields(s.AllFields(), isAutoIncr))
	if len(f) == 0 {
		return ""
	}

	return f[0]
}

// AutoIncrFieldName returns the name of the column marked with crud tag
// AutoIncr.
func (s Struct) AutoIncrFieldName() string {
	f := fieldNames(pickFields(s.AllFields(), isAutoIncr))
	if len(f) == 0 {
		return ""
	}

	return f[0]
}

// HasAutoIncrColumn returns true/false indicating if one of the columns is
// autoincrmenting (has crud tag = AutoIncr).
func (s Struct) HasAutoIncrColumn() bool {
	return s.AutoIncrColumnName() != ""
}

// InsertStatement produces the insert statement for a particular struct.
func (s Struct) InsertStatement(tableName string) string {

	idName := s.AutoIncrColumnName()

	if idName == "" {
		return crudlib.InsertStatement(tableName, s.insertColumnNames())
	}

	return crudlib.InsertReturningStatement(tableName, s.insertColumnNames(), idName)
}

// InsertBindVars returns field names suitable to include in template insert
// execution.
func (s Struct) InsertBindVars() string {
	return "r." + strings.Join(s.insertFieldNames(), ", r.")
}

// UpdateStatement returns the sql UPDATE sql for the struct.
func (s Struct) UpdateStatement(tableName string) string {
	return crudlib.UpdateStatement(tableName, s.keyColumnNames(), s.valueColumnNames())
}

// UpdateBindVars returns field names suitable to include in the template update
// execution.
func (s Struct) UpdateBindVars() string {
	return "r." + strings.Join(append(s.valueFieldNames(), s.keyFieldNames()...), ", r.")
}

// DeleteStatement returns the sql DELETE statement for the struct.
func (s Struct) DeleteStatement(tableName string) string {
	return crudlib.DeleteStatement(tableName, s.keyColumnNames())
}

// DeleteBindVars returns the field names required for a delete.
func (s Struct) DeleteBindVars() string {
	return "r." + strings.Join(s.keyFieldNames(), ", r.")
}

// SelectStatement returns the standard select statement for the struct.
func (s Struct) SelectStatement(tableName string) string {
	return crudlib.SelectStatement(tableName, s.selectColumnNames(), "")
}

// ScanVars produces an expression suitable for rows.Scan().
func (s Struct) ScanVars(structVarName string) string {
	prefix := "&" + structVarName + "."
	return prefix + strings.Join(
		s.selectFieldNames(), ", "+prefix)
}

func isCreateTimestamp(f Field) bool {
	return f.CrudTag == CreateTimestamp
}

func isUpdateTimestamp(f Field) bool {
	return f.CrudTag == UpdateTimestamp
}

// CreateTimestampFields returns range of fields with auto created timestamp.
func (s Struct) CreateTimestampFields() []Field {
	return pickFields(s.AllFields(), isCreateTimestamp)
}

// UpdateTimestampFields returns range of fields with auto updated timestamp.
func (s Struct) UpdateTimestampFields() []Field {
	return pickFields(s.AllFields(), isUpdateTimestamp)
}
