package main

import (
	"fmt"
	"log"
)

func main() {
	data := []byte("hello AND world OR alice AND NOT bob")

	fmt.Println("Scanning:")
	s := newScanner(data)
	for {
		token, err := s.Scan()
		if err != nil {
			log.Fatal(err)
		}
		printSpaces(4)
		fmt.Println(token)
		if token == tokenEOF {
			break
		}
	}

	fmt.Println("\nParsing:")
	s.reset()
	p := &Parser{s}
	query, err := p.parseQuery()
	if err != nil {
		log.Fatal(err)
	}
	query.PrettyPrint(4)
	fmt.Println()
}
