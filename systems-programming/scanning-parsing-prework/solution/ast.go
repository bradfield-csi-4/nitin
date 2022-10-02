package main

import (
	"bytes"
	"fmt"
)

func printSpaces(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf(" ")
	}
}

type Query interface {
	PrettyPrint(indent int)
	Render(buf *bytes.Buffer)
}

type OrNode struct {
	children []Query
}

func (n *OrNode) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("OR(\n")
	for i, child := range n.children {
		if i > 0 {
			fmt.Printf(",\n")
		}
		child.PrettyPrint(indent + 4)
	}
	fmt.Printf(")")
}

func (n *OrNode) Render(buf *bytes.Buffer) {
	buf.WriteString("OR(")
	for i, child := range n.children {
		if i > 0 {
			buf.WriteString(", ")
		}
		child.Render(buf)
	}
	buf.WriteString(")")
}

type AndNode struct {
	children []Query
}

func (n *AndNode) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("AND(\n")
	for i, child := range n.children {
		if i > 0 {
			fmt.Printf(",\n")
		}
		child.PrettyPrint(indent + 4)
	}
	fmt.Printf(")")
}

func (n *AndNode) Render(buf *bytes.Buffer) {
	buf.WriteString("AND(")
	for i, child := range n.children {
		if i > 0 {
			buf.WriteString(", ")
		}
		child.Render(buf)
	}
	buf.WriteString(")")
}

type NotNode struct {
	child Query
}

func (n *NotNode) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("NOT(\n")
	n.child.PrettyPrint(indent + 4)
	fmt.Printf(")")
}

func (n *NotNode) Render(buf *bytes.Buffer) {
	buf.WriteString("NOT(")
	n.child.Render(buf)
	buf.WriteString(")")
}

type Term struct {
	value string
}

func (t *Term) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("TERM(%s)", t.value)
}

func (t *Term) Render(buf *bytes.Buffer) {
	buf.WriteString("TERM(")
	buf.WriteString(t.value)
	buf.WriteString(")")
}

type Phrase struct {
	value string
}

func (p *Phrase) PrettyPrint(indent int) {
	printSpaces(indent)
	fmt.Printf("PHRASE(%s)", p.value)
}

func (p *Phrase) Render(buf *bytes.Buffer) {
	buf.WriteString("PHRASE(")
	buf.WriteString(p.value)
	buf.WriteString(")")
}

type Grouping struct {
	child Query
}

func (g *Grouping) PrettyPrint(indent int) {
	g.child.PrettyPrint(indent)
}

func (g *Grouping) Render(buf *bytes.Buffer) {
	g.child.Render(buf)
}
