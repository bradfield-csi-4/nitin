package main

import (
	"fmt"
)

type Parser struct {
	s *Scanner
}

func (p *Parser) parseQuery() (Query, error) {
	return p.parseOrQuery()
}

func (p *Parser) parseOrQuery() (Query, error) {
	var children []Query
	andQuery, err := p.parseAndQuery()
	if err != nil {
		return nil, err
	}
	children = append(children, andQuery)
	for {
		token, err := p.s.Peek()
		if err != nil {
			return nil, err
		}
		if token.tokenType != OR {
			break
		}
		err = p.s.Consume(OR)
		if err != nil {
			return nil, err
		}
		andQuery, err := p.parseAndQuery()
		if err != nil {
			return nil, err
		}
		children = append(children, andQuery)
	}
	if len(children) > 1 {
		return &OrNode{
			children: children,
		}, nil
	} else {
		return children[0], nil
	}
}

func (p *Parser) parseAndQuery() (Query, error) {
	var children []Query
	notQuery, err := p.parseNotQuery()
	if err != nil {
		return nil, err
	}
	children = append(children, notQuery)
	for {
		token, err := p.s.Peek()
		if err != nil {
			return nil, err
		}
		if token.tokenType == EOF || token.tokenType == R_PAREN || token.tokenType == OR {
			break
		} else if token.tokenType == AND {
			err = p.s.Consume(AND)
			if err != nil {
				return nil, err
			}
		}
		notQuery, err := p.parseNotQuery()
		if err != nil {
			return nil, err
		}
		children = append(children, notQuery)
	}
	if len(children) > 1 {
		return &AndNode{
			children: children,
		}, nil
	} else {
		return children[0], nil
	}
}

func (p *Parser) parseNotQuery() (Query, error) {
	notCount := 0
	for {
		token, err := p.s.Peek()
		if err != nil {
			return nil, err
		}
		if token.tokenType != NOT {
			break
		}
		err = p.s.Consume(NOT)
		if err != nil {
			return nil, err
		}
		notCount++
	}
	child, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	if notCount%2 == 0 {
		return child, nil
	} else {
		return &NotNode{
			child: child,
		}, nil
	}
}

func (p *Parser) parseAtom() (Query, error) {
	token, err := p.s.Scan()
	if err != nil {
		return nil, err
	}
	switch token.tokenType {
	case TERM:
		return &Term{
			value: token.literal,
		}, nil
	case PHRASE:
		return &Phrase{
			value: token.literal,
		}, nil
	case L_PAREN:
		child, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		err = p.s.Consume(R_PAREN)
		if err != nil {
			return nil, err
		}
		return &Grouping{
			child: child,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected token type %s", typeToString[token.tokenType])
	}
}
