package main

import (
	"bytes"
	"fmt"
)

type Expr interface {
	PrettyPrint(indent int)
	exprNode()
}

type AndNode struct {
	left  Expr
	right Expr
}

func (n AndNode) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("AND(\n")
	n.left.PrettyPrint(indent + 4)
	fmt.Printf(",\n")
	n.right.PrettyPrint(indent + 4)
	fmt.Printf(")")
}

type OrNode struct {
	left  Expr
	right Expr
}

func (n OrNode) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("OR(\n")
	n.left.PrettyPrint(indent + 4)
	fmt.Printf(",\n")
	n.right.PrettyPrint(indent + 4)
	fmt.Printf(")")
}

type NotNode struct {
	child Expr
}

func (n NotNode) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("NOT(\n")
	n.child.PrettyPrint(indent + 4)
	fmt.Printf(")")
}

type TermNode struct {
	value string
}

func (n TermNode) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("TERM(%s)", n.value)
}

func (n TermNode) Render(buf *bytes.Buffer) {
	buf.WriteString("TERM(")
	buf.WriteString(n.value)
	buf.WriteString(")")
}

func printSpaces(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf(" ")
	}
}

// exprNode() dummy methods, so that only these types can be assigned to an Expr
func (n AndNode) exprNode()  {}
func (n OrNode) exprNode()   {}
func (n NotNode) exprNode()  {}
func (n TermNode) exprNode() {}
