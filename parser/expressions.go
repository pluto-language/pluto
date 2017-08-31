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

	val, err := strconv.ParseFloat(p.cur.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %s as a number", p.cur.Literal)
		p.defaultErr(msg)
		return nil
	}

	lit.Value = val
	return lit
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
	expr.Right = p.parseExpression(prefix)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.next()

	if p.curIs(token.RPAREN) {
		return &ast.Tuple{
			Tok: p.cur,
		}
	}

	expr := p.parseExpression(lowest)
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
	}

	return &ast.Array{
		Tok:      p.cur,
		Elements: p.parseExpressionList(token.RSQUARE),
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
	expr.Condition = p.parseExpression(lowest)

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
	expr.Collection = p.parseExpression(lowest)

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

	expr.Pattern = append(expr.Pattern, p.parsePattern()...)

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
	expr.Condition = p.parseExpression(lowest)

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
	expr.Exp = p.parseExpression(lowest)

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
	expr.Value = p.parseExpression(lowest)

	return expr
}

func (p *Parser) parseShorthandAssignment(left ast.Expression) ast.Expression {
	expr := &ast.AssignExpression{
		Tok:  p.cur,
		Name: left,
	}

	op := p.cur.Literal

	p.next()
	right := p.parseExpression(lowest)

	expr.Value = &ast.InfixExpression{
		Left:     left,
		Operator: op[:len(op)-1],
		Right:    right,
	}

	return expr
}

func (p *Parser) parseDeclareExpression(left ast.Expression) ast.Expression {
	expr := &ast.DeclareExpression{
		Tok:  p.cur,
		Name: left,
	}

	p.next()
	expr.Value = p.parseExpression(lowest)

	return expr
}

func (p *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	expr := &ast.DotExpression{
		Tok:  p.cur,
		Left: left,
	}

	p.next()
	expr.Right = p.parseExpression(index)

	return expr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Tok:        p.cur,
		Collection: left,
	}

	p.next()
	expr.Index = p.parseExpression(lowest)

	if !p.expect(token.RSQUARE) {
		return nil
	}

	return expr
}

func (p *Parser) parseMethodCall(left ast.Expression) ast.Expression {
	return &ast.MethodCall{
		Tok:      p.cur,
		Instance: left,
		Pattern:  p.parsePattern(),
	}
}

func (p *Parser) parseQualifiedFunctionCall(left ast.Expression) ast.Expression {
	return &ast.QualifiedFunctionCall{
		Tok:     p.cur,
		Package: left,
		Pattern: p.parsePattern(),
	}
}

func (p *Parser) parsePattern() []ast.Expression {
	var pattern []ast.Expression

	for p.peekIs(p.argTokens...) || token.IsKeyword(p.peek.Type) {
		p.next()

		if token.IsKeyword(p.cur.Type) {
			pattern = append(pattern, &ast.Identifier{
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
			token.PARAM: func() ast.Expression {
				return arg(&ast.Identifier{
					Tok:   p.cur,
					Value: p.cur.Literal,
				})
			},
		}

		found := false
		for k, v := range p.prefixes {
			if _, hasHandler := handlers[k]; hasHandler {
				continue
			}

			if k == p.cur.Type {
				pattern = append(pattern, arg(v()))
				found = true
			}
		}

		if !found {
			handler := handlers[p.cur.Type]
			pattern = append(pattern, handler())
		}
	}

	if len(pattern) == 0 {
		p.defaultErr("expected at least one item in a pattern")
		return nil
	}

	return pattern
}
