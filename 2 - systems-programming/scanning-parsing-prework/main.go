package main

import (
	"fmt"
)

func main() {
	src := "hello world OR alice -bob"

	var token *Token
	var tokens []*Token

	s := newScanner(src)
	for {
		token = s.scan()
		if *token == tokenEOF {
			break
		}
		tokens = append(tokens, token)
		fmt.Println(*token)
	}

	fmt.Println()

	p := &Parser{tokens, 0}
	ast := p.parseQuery()
	ast.PrettyPrint(4)
}
