package main

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pdk/crudgen/crudlib"
)

// Values used to identify particular columns for insert/update/delete
// operations.
const (
	AutoIncr = "autoincr"
	Key      = "key"
)

// Field is what we found in a struct declaration.
type Field struct {
	Name    string
	DBTag   string
	CrudTag string
}

// Struct is a struct declaration we found in the source file.
type Struct struct {
	Name   string
	Fields []Field
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

// ColumnName returns either the DBTag value or snake-case of the Name.
func (f Field) ColumnName() string {
	if f.DBTag != "" {
		return f.DBTag
	}

	return strcase.ToSnake(f.Name)
}

func (s Struct) insertColumnNames() []string {
	v := []string{}
	for _, f := range s.Fields {
		if f.CrudTag != AutoIncr {
			v = append(v, f.ColumnName())
		}
	}

	return v
}

func (s Struct) insertFieldNames() []string {
	v := []string{}
	for _, f := range s.Fields {
		if f.CrudTag != AutoIncr {
			v = append(v, f.Name)
		}
	}

	return v
}

func (s Struct) valueColumnNames() []string {
	v := []string{}
	for _, f := range s.Fields {
		if f.CrudTag != AutoIncr && f.CrudTag != Key {
			v = append(v, f.ColumnName())
		}
	}

	return v
}

func (s Struct) valueFieldNames() []string {
	v := []string{}
	for _, f := range s.Fields {
		if f.CrudTag != AutoIncr && f.CrudTag != Key {
			v = append(v, f.Name)
		}
	}

	return v
}

func (s Struct) keyColumnNames() []string {
	v := []string{}
	for _, f := range s.Fields {
		if f.CrudTag == AutoIncr || f.CrudTag == Key {
			v = append(v, f.ColumnName())
		}
	}

	return v
}

func (s Struct) keyFieldNames() []string {
	v := []string{}
	for _, f := range s.Fields {
		if f.CrudTag == AutoIncr || f.CrudTag == Key {
			v = append(v, f.Name)
		}
	}

	return v
}

// AutoIncrColumnName returns the name of the column marked with crud tag
// AutoIncr.
func (s Struct) AutoIncrColumnName() string {
	for _, f := range s.Fields {
		if f.CrudTag == AutoIncr {
			return f.ColumnName()
		}
	}

	return ""
}

// AutoIncrFieldName returns the name of the column marked with crud tag
// AutoIncr.
func (s Struct) AutoIncrFieldName() string {
	for _, f := range s.Fields {
		if f.CrudTag == AutoIncr {
			return f.Name
		}
	}

	return ""
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
