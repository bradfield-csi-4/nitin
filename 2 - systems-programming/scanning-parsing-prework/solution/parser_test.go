package main

import (
	"bytes"
	"testing"
)

func TestParser(t *testing.T) {
	for _, testCase := range []struct {
		input    []byte
		expected string
	}{
		{
			[]byte("alice AND bob"),
			"AND(TERM(alice), TERM(bob))",
		},
		{
			[]byte("alice AND -bob"),
			"AND(TERM(alice), NOT(TERM(bob)))",
		},
		{
			[]byte("\"hello, dave\""),
			"PHRASE(\"hello, dave\")",
		},
		{
			[]byte("carol OR \"hello, \\\" dave\""),
			"OR(TERM(carol), PHRASE(\"hello, \\\" dave\"))",
		},
		{
			[]byte("(alice OR bob) and (carol AND (dave or eve))"),
			"AND(OR(TERM(alice), TERM(bob)), AND(TERM(carol), OR(TERM(dave), TERM(eve))))",
		},
	} {
		s := newScanner(testCase.input)
		p := &Parser{s}
		q, err := p.parseQuery()
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		q.Render(&buf)
		actual := buf.String()
		if testCase.expected != actual {
			t.Fatalf("unexpected parser output:\nexpected:\t%v\nactual:\t%v", testCase.expected, actual)
		}
	}
}
