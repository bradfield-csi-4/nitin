package main

import (
	"bytes"
	"log"
	"os"
	"sort"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const src string = `package foo

import (
	"fmt"
	"time"
)

func baz() {
	fmt.Println("Hello, world!")
}

type A int

const b = "testing"

func bar() {
	fmt.Println(time.Now())
}`

// Moves all top-level functions to the end, sorted in alphabetical order.
// The "source file" is given as a string (rather than e.g. a filename).
func SortFunctions(src string) (string, error) {
	f, err := decorator.Parse(src)
	if err != nil {
		return "", err
	}

	var functions []dst.Decl
	var result []dst.Decl

	for i := 0; i < len(f.Decls); i++ {
		_, isFunc := f.Decls[i].(*dst.FuncDecl)
		if isFunc {
			functions = append(functions, f.Decls[i])
		} else {
			result = append(result, f.Decls[i])
		}
	}

	sort.Slice(functions, func(i, j int) bool {
		return functions[i].(*dst.FuncDecl).Name.Name < functions[j].(*dst.FuncDecl).Name.Name
	})

	result = append(result, functions...)
	f.Decls = result

	var buff bytes.Buffer
	err = decorator.Fprint(&buff, f)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func main() {
	f, err := decorator.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	// Print AST
	err = dst.Fprint(os.Stdout, f, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Convert AST back to source
	err = decorator.Print(f)
	if err != nil {
		log.Fatal(err)
	}
}
