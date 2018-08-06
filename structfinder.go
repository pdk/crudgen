package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
)

// FindStructs parses a .go file and finds structs defined therein. Returns a
// map of struct names to Structs (which contains info discovered about the
// structs).
func FindStructs(sourceFileName string) map[string]*Struct {

	fs := token.NewFileSet()

	f, err := parser.ParseFile(fs, sourceFileName, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("%s", err)
	}

	finder := StructFinder{
		Structs: make(map[string]*Struct),
	}

	ast.Walk(&finder, f)

	return finder.Structs
}

// StructFinder is a thing that can Visit() nodes.
type StructFinder struct {
	Structs map[string]*Struct
}

// Visit checks if the current node is a struct, and processes it
// if so.
func (v *StructFinder) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	d, ok := n.(*ast.TypeSpec)
	if !ok {
		return v
	}

	// it's a type declaration

	d2, ok := d.Type.(*ast.StructType)
	if !ok {
		return v
	}

	// it's a struct declaration

	s := NewStruct(d.Name.Name)

	for _, f := range (*d2.Fields).List {
		for _, n := range f.Names {

			f2, ok := n.Obj.Decl.(*ast.Field)
			if ok && f2.Tag != nil {
				tag := reflect.StructTag(stripSingleQuotes(f2.Tag.Value))
				s.AppendField(n.Name, tag.Get("db"), tag.Get("crud"))

				continue
			}

			s.AppendField(n.Name, "", "")
		}
	}

	v.Structs[s.Name] = s

	return v
}

func stripSingleQuotes(s string) string {
	if len(s) < 2 {
		return s
	}

	return s[1 : len(s)-1]
}
