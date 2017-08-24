package parser

import (
	"fmt"
	"strconv"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/token"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, prefixExists := p.prefixes[p.cur.Type]

	if !prefixExists {
		p.unexpectedTokenErr(p.cur.Type)
		return nil
	}

	left := prefix()

	for !p.peekIs(token.SEMI) && precedence < p.peekPrecedence() {
		infix, infixExists := p.infixes[p.peek.Type]

		if !infixExists {
			return left
		}

		p.next()
		left = infix(left)
	}

	return left
}

/**********************
 * Prefix expressions *
 **********************/

func (p *Parser) parseID() ast.Expression {
	if p.peekIs(p.argTokens...) {
		return p.parseFunctionCallStartingWith(&ast.Identifier{
			Tok:   p.cur,
			Value: p.cur.Literal,
		})
	}

	return &ast.Identifier{
		Tok:   p.cur,
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseNonFnID() ast.Expression {
	return &ast.Identifier{
		Tok:   p.cur,
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseNum() ast.Expression {
	lit := &ast.Number{
		Tok: p.cur,
	}

	if val, err := strconv.ParseFloat(p.cur.Literal, 64); err != nil {
		msg := fmt.Sprintf("could not parse %s as a number", p.cur.Literal)
		p.defaultErr(msg)
		return nil
	} else {
		lit.Value = val
		return lit
	}
}

func (p *Parser) parseBool() ast.Expression {
	return &ast.Boolean{
		Tok:   p.cur,
		Value: p.cur.Type == token.TRUE,
	}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.Null{
		Tok: p.cur,
	}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{
		Tok:   p.cur,
		Value: p.cur.Literal,
	}
}

func (p *Parser) parseChar() ast.Expression {
	return &ast.Char{
		Tok:   p.cur,
		Value: p.cur.Literal[0],
	}
}

func (p *Parser) parsePrefix() ast.Expression {
	expr := &ast.PrefixExpression{
		Tok:      p.cur,
		Operator: p.cur.Literal,
	}

	p.next()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.next()

	if p.curIs(token.RPAREN) {
		return &ast.Tuple{
			Tok: p.cur,
		}
	}

	expr := p.parseExpression(LOWEST)
	isTuple := false

	if p.peekIs(token.COMMA) {
		isTuple = true

		p.next()
		p.next()

		expr = &ast.Tuple{
			Tok: expr.Token(),
			Value: append(
				[]ast.Expression{expr},
				p.parseExpressionList(token.RPAREN)...,
			),
		}
	}

	if !isTuple && !p.expect(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseArrayOrMap() ast.Expression {
	p.next()

	if p.peekIs(token.COLON) || p.curIs(token.COLON) {
		pairs := p.parseExpressionPairs(token.RSQUARE)

		return &ast.Map{
			Tok:   p.cur,
			Pairs: pairs,
		}
	} else {
		return &ast.Array{
			Tok:      p.cur,
			Elements: p.parseExpressionList(token.RSQUARE),
		}
	}
}

func (p *Parser) parseBlockLiteral() ast.Expression {
	expr := &ast.BlockLiteral{
		Tok: p.cur,
	}

	if p.peekIs(token.B_OR) {
		p.next()
		expr.Params = p.parseParams(token.B_OR)

		if !p.expect(token.ARROW) {
			return nil
		}
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseWhileLoop() ast.Expression {
	expr := &ast.WhileLoop{
		Tok: p.cur,
	}

	if !p.expect(token.LPAREN) {
		return nil
	}

	p.next()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseForLoop() ast.Expression {
	expr := &ast.ForLoop{
		Tok: p.cur,
	}

	if !p.expect(token.LPAREN) {
		return nil
	}

	p.next()
	expr.Var = p.parseNonFnID()

	if !p.expect(token.COLON) {
		return nil
	}

	p.next()
	expr.Collection = p.parseExpression(LOWEST)

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseFunctionCall() ast.Expression {
	return p.parseFunctionCallStartingWith(nil)
}

func (p *Parser) parseFunctionCallStartingWith(start ast.Expression) ast.Expression {
	expr := &ast.FunctionCall{
		Tok: p.cur,
	}

	if start != nil {
		expr.Pattern = append(expr.Pattern, start)
	}

	keywordTokens := []token.Type{}
	for _, t := range token.Keywords {
		keywordTokens = append(keywordTokens, t)
	}

	for p.peekIs(p.argTokens...) || p.peekIs(keywordTokens...) {
		p.next()

		if p.curIs(keywordTokens...) {
			expr.Pattern = append(expr.Pattern, &ast.Identifier{
				Tok:   p.cur,
				Value: p.cur.Literal,
			})
			continue
		}

		arg := func(val ast.Expression) ast.Expression {
			return &ast.Argument{
				Tok:   p.cur,
				Value: val,
			}
		}

		type Handler func() ast.Expression

		handlers := map[token.Type]Handler{
			token.ID: func() ast.Expression {
				return &ast.Identifier{
					Tok: p.cur, Value: p.cur.Literal,
				}
			},
			token.PARAM: func() ast.Expression {
				return arg(&ast.Identifier{
					Tok: p.cur, Value: p.cur.Literal,
				})
			},
		}

		found := false
		for k, v := range p.prefixes {
			if _, hasHandler := handlers[k]; hasHandler {
				continue
			}

			if k == p.cur.Type {
				expr.Pattern = append(expr.Pattern, arg(v()))
				found = true
			}
		}

		if !found {
			handler := handlers[p.cur.Type]
			expr.Pattern = append(expr.Pattern, handler())
		}
	}

	if len(expr.Pattern) == 0 {
		p.defaultErr("expected at least one item in a pattern")
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{
		Tok: p.cur,
	}

	if !p.expect(token.LPAREN) {
		return nil
	}

	p.next()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peekIs(token.ELSE) {
		p.next()

		if !p.expect(token.LBRACE) {
			return nil
		}

		expr.Alternative = p.parseBlockStatement()
	} else if p.peekIs(token.ELIF) {
		p.next()

		expr.Alternative = &ast.BlockStatement{
			Tok: p.cur,
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Tok:  p.cur,
					Expr: p.parseIfExpression(),
				},
			},
		}
	}

	return expr
}

func (p *Parser) parseMatchExpression() ast.Expression {
	expr := &ast.MatchExpression{
		Tok: p.cur,
	}

	if !p.expect(token.LPAREN) {
		return nil
	}

	p.next()
	expr.Exp = p.parseExpression(LOWEST)

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Arms = p.parseMatchArms()
	if expr.Arms == nil {
		return nil
	}

	return expr
}

func (p *Parser) parseTryExpression() ast.Expression {
	expr := &ast.TryExpression{
		Tok: p.cur,
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	if !p.expect(token.CATCH) {
		return nil
	}

	if !p.expect(token.LPAREN) {
		return nil
	}

	p.next()
	expr.ErrName = p.parseNonFnID()

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Arms = p.parseMatchArms()

	return expr
}

/*********************
 * Infix expressions *
 *********************/

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Tok:      p.cur,
		Operator: p.cur.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.next()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	expr := &ast.AssignExpression{
		Tok:  p.cur,
		Name: left,
	}

	p.next()
	expr.Value = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseDeclareExpression(left ast.Expression) ast.Expression {
	expr := &ast.DeclareExpression{
		Tok:  p.cur,
		Name: left,
	}

	p.next()
	expr.Value = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	expr := &ast.DotExpression{
		Tok:  p.cur,
		Left: left,
	}

	p.next()
	expr.Right = p.parseExpression(DOT)

	return expr
}

func (p *Parser) parseMethodCall(left ast.Expression) ast.Expression {
	expr := &ast.MethodCall{
		Tok:      p.cur,
		Instance: left,
	}

	for p.peekIs(p.argTokens...) || token.IsKeyword(p.peek.Type) {
		p.next()

		if token.IsKeyword(p.cur.Type) {
			expr.Pattern = append(expr.Pattern, &ast.Identifier{
				Tok:   p.cur,
				Value: p.cur.Literal,
			})
			continue
		}

		arg := func(val ast.Expression) ast.Expression {
			return &ast.Argument{
				Tok:   p.cur,
				Value: val,
			}
		}

		type Handler func() ast.Expression

		handlers := map[token.Type]Handler{
			token.ID: func() ast.Expression {
				return &ast.Identifier{
					Tok:   p.cur,
					Value: p.cur.Literal,
				}
			},
			token.LPAREN: func() ast.Expression {
				return arg(p.parseGroupedExpression())
			},
			token.NUM: func() ast.Expression {
				return arg(p.parseNum())
			},
			token.NULL: func() ast.Expression {
				return arg(p.parseNull())
			},
			token.TRUE: func() ast.Expression {
				return arg(p.parseBool())
			},
			token.FALSE: func() ast.Expression {
				return arg(p.parseBool())
			},
			token.STR: func() ast.Expression {
				return arg(p.parseString())
			},
			token.PARAM: func() ast.Expression {
				return arg(&ast.Identifier{
					Tok:   p.cur,
					Value: p.cur.Literal,
				})
			},
			token.LSQUARE: func() ast.Expression {
				return arg(p.parseArrayOrMap())
			},
			token.LBRACE: func() ast.Expression {
				return arg(p.parseBlockLiteral())
			},
		}

		handler := handlers[p.cur.Type]
		expr.Pattern = append(expr.Pattern, handler())
	}

	if len(expr.Pattern) == 0 {
		p.defaultErr("expected at least one item in a pattern")
		return nil
	}

	return expr
}
