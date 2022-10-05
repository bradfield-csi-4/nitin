package main

type Parser struct {
	tokens  []*Token
	current int
}

func (p *Parser) parseQuery() Expr {
	return p.parseOrQuery()
}

func (p *Parser) parseOrQuery() Expr {
	expr := p.parseAndQuery()

	for p.match(OR) {
		p.current++
		expr = OrNode{expr, p.parseAndQuery()}
	}

	return expr
}

func (p *Parser) parseAndQuery() Expr {
	expr := p.parseNotQuery()

	for p.match(AND) {
		p.current++
		expr = AndNode{expr, p.parseNotQuery()}
	}

	return expr
}

func (p *Parser) parseNotQuery() Expr {
	for p.match(NOT) {
		p.current++
		return NotNode{p.parseNotQuery()}
	}
	return p.parseTermNode()
}

func (p *Parser) parseTermNode() Expr {
	currentToken := p.tokens[p.current]
	p.current++
	if p.isAtEnd() {
		return TermNode{currentToken.lexeme}
	}

	nextToken := p.tokens[p.current]
	if match(nextToken, STRING) {
		return AndNode{TermNode{currentToken.lexeme}, p.parseTermNode()}
	} else if match(nextToken, NOT) {
		return AndNode{TermNode{currentToken.lexeme}, p.parseNotQuery()}
	}

	return TermNode{currentToken.lexeme}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) match(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return match(p.tokens[p.current], tokenType)
}

func match(token *Token, tokenType TokenType) bool {
	return token.tokenType == tokenType
}
