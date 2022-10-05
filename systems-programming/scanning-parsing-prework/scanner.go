package main

import (
	"strings"
)

type TokenType int

const (
	AND TokenType = iota
	OR
	NOT
	STRING
	EOF
)

var keywordToType = map[string]TokenType{
	"AND": AND,
	"OR":  OR,
	"NOT": NOT,
}

type Token struct {
	tokenType TokenType
	lexeme    string
}

type Scanner struct {
	src string
	idx int
}

func newScanner(src string) *Scanner {
	return &Scanner{
		strings.TrimSpace(src),
		0}
}

var tokenEOF = Token{EOF, ""}
var tokenNOT = Token{NOT, ""}

func (s *Scanner) isAtEnd() bool {
	return s.idx >= len(s.src)
}

// Given a starting index, stop at the first whitespace, and return the corresponding token
func (s *Scanner) scan() *Token {
	if s.isAtEnd() {
		return &tokenEOF
	}

	idx := s.idx
	for {
		if idx == len(s.src) {
			break
		}
		if s.src[idx] == ' ' {
			break
		}
		idx++
	}

	if s.src[s.idx] == '-' {
		s.idx++
		return &Token{NOT, "-"}
	}

	lexeme := s.src[s.idx:idx]
	s.idx = idx + 1
	return &Token{s.getTokenType(lexeme), lexeme}
}

func (s *Scanner) getTokenType(lexeme string) TokenType {
	reservedWordTokenType, exists := keywordToType[lexeme]
	if exists {
		return reservedWordTokenType
	} else {
		return STRING
	}
}
