package parser

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/token"
)

func (p *Parser) parseStatement() ast.Statement {
	var stmt ast.Statement

	if p.curIs(token.Semi) {
		return nil
	} else if p.curIs(token.Return) {
		stmt = p.parseReturnStatement()
	} else if p.curIs(token.Def) {
		stmt = p.parseDefStatement()
	} else if p.curIs(token.Next) {
		stmt = p.parseNextStatement()
	} else if p.curIs(token.Break) {
		stmt = p.parseBreakStatement()
	} else if p.curIs(token.Import) {
		stmt = p.parseImportStatement()
	} else if p.curIs(token.Use) {
		stmt = p.parseUseStatement()
	} else if p.curIs(token.While) {
		stmt = p.parseWhileLoop()
	} else if p.curIs(token.For) {
		stmt = p.parseForLoop()
	} else {
		stmt = p.parseExpressionStatement()
	}

	if !p.expect(token.Semi) {
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

	for !p.curIs(token.RightBrace) && !p.curIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.next()
	}

	return block
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Tok:  p.cur,
		Expr: p.parseExpression(lowest),
	}

	if stmt.Expr == nil {
		return nil
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	if p.peek.Type == token.Semi {
		return &ast.ReturnStatement{
			Tok: p.cur,
		}
	}

	stmt := &ast.ReturnStatement{
		Tok: p.cur,
	}

	p.next()
	stmt.Value = p.parseExpression(lowest)

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
	stmt.Pattern = p.parsePatternCall(token.LeftBrace)

	if len(stmt.Pattern) == 0 {
		p.defaultErr("expected at least one item in a pattern")
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseImportStatement() ast.Statement {
	stmt := &ast.ImportStatement{
		Tok: p.cur,
	}

	if !p.expect(token.String) {
		return nil
	}

	stmt.Package = p.cur.Literal

	return stmt
}

func (p *Parser) parseUseStatement() ast.Statement {
	stmt := &ast.UseStatement{
		Tok: p.cur,
	}

	if !p.expect(token.String) {
		return nil
	}

	stmt.Package = p.cur.Literal

	return stmt
}

func (p *Parser) parseWhileLoop() ast.Statement {
	stmt := &ast.WhileLoop{
		Tok: p.cur,
	}

	if !p.expect(token.LeftParen) {
		return nil
	}

	p.next()
	stmt.Condition = p.parseExpression(lowest)

	if !p.expect(token.RightParen) {
		return nil
	}

	if !p.expect(token.LeftBrace) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForLoop() ast.Statement {
	stmt := &ast.ForLoop{
		Tok: p.cur,
	}

	if !p.expect(token.LeftParen) {
		return nil
	}

	p.next()
	stmt.Var = p.parseID()

	if !p.expect(token.In) {
		return nil
	}

	p.next()
	stmt.Collection = p.parseExpression(lowest)

	if !p.expect(token.RightParen) {
		return nil
	}

	if !p.expect(token.LeftBrace) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}
