package parser

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/token"
)

func (p *Parser) parseStatement() ast.Statement {
	var stmt ast.Statement

	if p.curIs(token.SEMI) {
		return nil
	} else if p.curIs(token.RETURN) {
		stmt = p.parseReturnStatement()
	} else if p.curIs(token.DEF) {
		stmt = p.parseDefStatement()
	} else if p.curIs(token.NEXT) {
		stmt = p.parseNextStatement()
	} else if p.curIs(token.BREAK) {
		stmt = p.parseBreakStatement()
	} else if p.curIs(token.CLASS) {
		stmt = p.parseClassDeclaration()
	} else {
		stmt = p.parseExpressionStatement()
	}

	if !p.expect(token.SEMI) {
		return nil
	}

	return stmt
}

func (p *Parser) parseBlockStatement() ast.Statement {
	block := &ast.BlockStatement{
		Tok:        p.cur,
		Statements: []ast.Statement{},
	}

	p.next()

	for !p.curIs(token.RBRACE) && !p.curIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.next()
	}

	return block
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	return &ast.ExpressionStatement{
		Tok:  p.cur,
		Expr: p.parseExpression(LOWEST),
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	if p.peek.Type == token.SEMI {
		return &ast.ReturnStatement{
			Tok: p.cur,
		}
	}

	stmt := &ast.ReturnStatement{
		Tok: p.cur,
	}

	p.next()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseNextStatement() ast.Statement {
	return &ast.NextStatement{
		Tok: p.cur,
	}
}

func (p *Parser) parseBreakStatement() ast.Statement {
	return &ast.BreakStatement{
		Tok: p.cur,
	}
}

func (p *Parser) parseDefStatement() ast.Statement {
	stmt := &ast.FunctionDefinition{
		Tok: p.cur,
	}

	p.next()
	stmt.Pattern = p.parsePatternCall(token.LBRACE)

	if len(stmt.Pattern) == 0 {
		p.defaultErr("expected at least one item in a pattern")
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseInitStatement() ast.Statement {
	stmt := &ast.InitDefinition{
		Tok: p.cur,
	}

	p.next()
	stmt.Pattern = p.parsePatternCall(token.LBRACE)

	if len(stmt.Pattern) == 0 {
		p.defaultErr("expected at least one item in a pattern")
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseClassDeclaration() ast.Statement {
	stmt := &ast.ClassStatement{
		Tok: p.cur,
	}

	if !p.expect(token.ID) {
		return nil
	}

	stmt.Name = p.parseNonFnID()

	if p.peekIs(token.EXTENDS) {
		p.next()
		p.next()

		stmt.Parent = p.parseNonFnID()
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	p.next()

	if !p.curIs(token.RBRACE) {
		return stmt
	}

	for p.curIs(token.INIT, token.DEF) {
		if p.curIs(token.INIT) {
			stmt.Methods = append(stmt.Methods, p.parseInitStatement())
		} else {
			stmt.Methods = append(stmt.Methods, p.parseDefStatement())
		}

		if !p.expect(token.SEMI) {
			return nil
		}

		if p.peekIs(token.INIT, token.DEF) {
			p.next()
		}
	}

	if !p.expect(token.RBRACE) {
		return nil
	}

	return stmt
}
