package parser

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/token"
)

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	exprs := []ast.Expression{}

	if p.curIs(end) {
		return exprs
	}

	exprs = append(exprs, p.parseExpression(LOWEST))

	for p.peekIs(token.COMMA) {
		p.next()

		if p.peekIs(end) {
			p.next()
			return exprs
		}

		p.next()
		exprs = append(exprs, p.parseExpression(LOWEST))
	}

	if !p.expect(end) {
		return nil
	}

	return exprs
}

func (p *Parser) parseExpressionPairs(end token.Type) map[ast.Expression]ast.Expression {
	pairs := map[ast.Expression]ast.Expression{}

	if p.curIs(token.COLON) {
		p.next()
		return pairs
	}

	key, val := p.parsePair()
	pairs[key] = val

	for p.peekIs(token.COMMA) {
		p.next()

		if p.peekIs(end) {
			p.next()
			return pairs
		}

		p.next()
		key, val = p.parsePair()
		pairs[key] = val
	}

	if !p.expect(end) {
		return nil
	}

	return pairs
}

func (p *Parser) parsePair() (ast.Expression, ast.Expression) {
	key := p.parseExpression(DOT)

	if !p.expect(token.COLON) {
		return nil, nil
	}

	p.next()

	value := p.parseExpression(LOWEST)

	return key, value
}

func (p *Parser) parseParams(end token.Type) []ast.Expression {
	params := []ast.Expression{}

	if p.peekIs(end) {
		p.next()
		return params
	}

	p.next()
	params = append(params, p.parseNonFnID())

	for p.peekIs(token.COMMA) {
		p.next()
		p.next()
		params = append(params, p.parseNonFnID())
	}

	if !p.expect(end) {
		return nil
	}

	return params
}

func (p *Parser) parseMatchArm() ast.Arm {
	var (
		left  = []ast.Expression{}
		right ast.Statement
	)

	if p.curIs(token.STAR) {
		left = nil
		p.next()
	} else {
		left = p.parseExpressionList(token.F_ARROW)
	}

	p.next()

	if p.curIs(token.LBRACE) {
		right = p.parseBlockStatement()
		p.next()
	} else {
		right = p.parseStatement()
	}

	return ast.Arm{
		Exprs: left,
		Body:  right,
	}
}

func (p *Parser) parseMatchArms() []ast.Arm {
	var arms []ast.Arm

	for !p.curIs(token.RBRACE) {
		p.next()
		arm := p.parseMatchArm()

		arms = append(arms, arm)

		if p.peekIs(token.RBRACE) {
			p.next()
		}
	}

	return arms
}

func (p *Parser) parsePatternCall(end token.Type) []ast.Expression {
	var (
		pattern      []ast.Expression
		allowedItems = []token.Type{token.ID, token.PARAM}
	)

	for _, t := range token.Keywords {
		allowedItems = append(allowedItems, t)
	}

	for !p.curIs(end) {
		tok := p.cur

		if !p.expectCur(allowedItems...) {
			return nil
		}

		if tok.Type == token.ID || token.IsKeyword(tok.Type) {
			pattern = append(pattern, &ast.Identifier{
				Tok:   tok,
				Value: tok.Literal,
			})
		} else {
			pattern = append(pattern, &ast.Parameter{
				Tok:  tok,
				Name: tok.Literal,
			})
		}
	}

	return pattern
}
