package main

import (
	"flag"
	"log"
	"os"

	"github.com/pdk/crudgen/crudlib"
)

func init() {
	flag.StringVar(&packageName, "package", "", "name of package for the created .go file")
	flag.StringVar(&outFileName, "out", "", "name of file to create/write")
	flag.StringVar(&structName, "struct", "", "name of struct to generate CRUD for")
	flag.StringVar(&tableName, "table", "", "name of the database table")
	flag.Var(&bindStyle, "bindstyle", "bind style, one of 'question', 'dollar' or 'named'")
	flag.StringVar(&selectName, "select", "Select", "name of the function for executing a query")

	flag.Parse()
}

var (
	packageName string
	structName  string
	tableName   string
	bindStyle   crudlib.BindStyle
	outFileName string
	selectName  string
)

func main() {

	if tableName == "" || packageName == "" {
		log.Fatalf("-table & -package are required arguments.\n")
		os.Exit(1)
	}

	structs := make(map[string]*Struct)
	for _, fn := range flag.Args() {
		merge(structs, FindStructs(fn))
	}

	s := selectStruct(structName, structs)

	outFile := resolveOutFile(outFileName)

	PrintTemplate(outFile, packageName, tableName, *s, bindStyle, selectName)
}

// selectStruct figures out which of the structs found in the .go file to process.
func selectStruct(structName string, structs map[string]*Struct) *Struct {

	if len(structs) == 0 {
		log.Fatalf("no structs found.\n")
	}

	if structName != "" {
		s, ok := structs[structName]
		if !ok {
			log.Fatalf("struct %s not found.\n", structName)
		}

		return s
	}

	if len(structs) != 1 {
		log.Fatalf("multiple structs found. use -struct to specify which to process.\n")
	}

	// get the struct name of the one item in the map
	for k := range structs {
		structName = k
	}

	return structs[structName]
}

// resolveOutFile figures out where we're gonna send the output.
func resolveOutFile(outFileName string) *os.File {

	if outFileName == "" {
		return os.Stdout
	}

	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatalf("%s", err)
	}

	return outFile
}

func merge(dst, src map[string]*Struct) {
	for k, v := range src {
		dst[k] = v
	}
}
