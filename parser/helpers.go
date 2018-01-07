package parser

import (
	"fmt"
	"math"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/token"
)

func (p *Parser) curIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.cur.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) peekIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.peek.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) expect(ts ...token.Type) bool {
	if p.peekIs(ts...) {
		p.next()
		return true
	}

	p.peekErr(ts...)
	return false
}

func (p *Parser) expectCur(ts ...token.Type) bool {
	if p.curIs(ts...) {
		p.next()
		return true
	}

	p.curErr(ts...)
	return false
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	exprs := []ast.Expression{}

	if p.curIs(end) {
		return exprs
	}

	exprs = append(exprs, p.parseExpression(lowest))

	for p.peekIs(token.Comma) {
		p.next()

		if p.peekIs(end) {
			p.next()
			return exprs
		}

		p.next()
		exprs = append(exprs, p.parseExpression(lowest))
	}

	if !p.expect(end) {
		return nil
	}

	return exprs
}

func (p *Parser) parseEmissionList() []ast.EmittedItem {
	items := []ast.EmittedItem{}

	if p.curIs(token.GreaterThan) {
		return items
	}

	items = append(items, p.parseEmittedItem())

	for p.peekIs(token.Comma) {
		p.next()

		if p.peekIs(token.GreaterThan) {
			p.next()
			return items
		}

		p.next()
		items = append(items, p.parseEmittedItem())
	}

	if !p.expect(token.GreaterThan) {
		return nil
	}

	return items
}

func (p *Parser) parseEmittedItem() ast.EmittedItem {
	if p.curIs(token.ID) {
		for _, data := range bytecode.Instructions {
			if data.Name == p.cur.Literal {
				item := ast.EmittedItem{
					IsInstruction: true,
					Instruction:   data.Name,
				}

				if p.peekIs(token.Number) {
					p.next()

					var (
						num = p.parseNum()
						arg = num.(*ast.Number).Value
					)

					if math.Floor(arg) != arg {
						p.Err(fmt.Sprintf("non-integer instruction argument %g", arg), p.cur.Start, p.cur.End)
						return item
					} else if arg < 0 {
						p.Err(fmt.Sprintf("instruction argument %g is less than 0", arg), p.cur.Start, p.cur.End)
						return item
					} else if arg > 1<<16 {
						p.Err(fmt.Sprintf("instruction argument %g is more than 0xFFFF (maximum uint16)", arg), p.cur.Start, p.cur.End)
						return item
					}

					item.Argument = rune(arg)
				}

				return item
			}
		}
	}

	return ast.EmittedItem{
		Exp: p.parseExpression(compare), // precedence = compare because > is an operator
	}
}

func (p *Parser) parseExpressionPairs(end token.Type) map[ast.Expression]ast.Expression {
	pairs := map[ast.Expression]ast.Expression{}

	if p.curIs(token.Colon) {
		p.next()
		return pairs
	}

	key, val := p.parsePair()
	pairs[key] = val

	for p.peekIs(token.Comma) {
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
	key := p.parseExpression(index)

	if !p.expect(token.Colon) {
		return nil, nil
	}

	p.next()

	value := p.parseExpression(lowest)

	return key, value
}

func (p *Parser) parseParams(end token.Type) []ast.Expression {
	params := []ast.Expression{}

	if p.peekIs(end) {
		p.next()
		return params
	}

	p.next()
	params = append(params, p.parseID())

	for p.peekIs(token.Comma) {
		p.next()
		p.next()
		params = append(params, p.parseID())
	}

	if !p.expect(end) {
		return nil
	}

	return params
}

func (p *Parser) parsePatternCall(end token.Type) []ast.Expression {
	var (
		pattern      []ast.Expression
		allowedItems = []token.Type{token.ID, token.Param}
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
