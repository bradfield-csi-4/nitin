package main

import (
	"fmt"
)

func main() {
	src := "hello AND world OR alice AND NOT bob"

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
	ast := p.parse()
	ast.PrettyPrint(4)
}
