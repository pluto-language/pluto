package parser

import (
	"fmt"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/lexer"
	"github.com/Zac-Garby/pluto/token"
)

type prefixParser func() ast.Expression
type infixParser func(ast.Expression) ast.Expression

type Parser struct {
	Errors []Error

	lex       func() token.Token
	text      string
	cur, peek token.Token
	prefixes  map[token.Type]prefixParser
	infixes   map[token.Type]infixParser
	argTokens []token.Type
}

func New(text string) *Parser {
	p := &Parser{
		lex:    lexer.Lexer(text),
		text:   text,
		Errors: []Error{},
	}

	p.prefixes = map[token.Type]prefixParser{
		token.ID:      p.parseID,
		token.NUM:     p.parseNum,
		token.TRUE:    p.parseBool,
		token.FALSE:   p.parseBool,
		token.NULL:    p.parseNull,
		token.LSQUARE: p.parseArrayOrMap,
		token.STR:     p.parseString,
		token.CHAR:    p.parseChar,

		token.MINUS: p.parsePrefix,
		token.PLUS:  p.parsePrefix,
		token.BANG:  p.parsePrefix,

		token.LPAREN: p.parseGroupedExpression,
		token.IF:     p.parseIfExpression,
		token.BSLASH: p.parseFunctionCall,
		token.LBRACE: p.parseBlockLiteral,
		token.WHILE:  p.parseWhileLoop,
		token.FOR:    p.parseForLoop,
		token.MATCH:  p.parseMatchExpression,
		token.TRY:    p.parseTryExpression,
	}

	p.infixes = map[token.Type]infixParser{
		token.PLUS:     p.parseInfix,
		token.MINUS:    p.parseInfix,
		token.STAR:     p.parseInfix,
		token.SLASH:    p.parseInfix,
		token.EQ:       p.parseInfix,
		token.N_EQ:     p.parseInfix,
		token.LT:       p.parseInfix,
		token.GT:       p.parseInfix,
		token.OR:       p.parseInfix,
		token.AND:      p.parseInfix,
		token.B_OR:     p.parseInfix,
		token.B_AND:    p.parseInfix,
		token.EXP:      p.parseInfix,
		token.F_DIV:    p.parseInfix,
		token.MOD:      p.parseInfix,
		token.LTE:      p.parseInfix,
		token.GTE:      p.parseInfix,
		token.Q_MARK:   p.parseInfix,
		token.A_AND:    p.parseShorthandAssignment,
		token.A_B_AND:  p.parseShorthandAssignment,
		token.A_B_OR:   p.parseShorthandAssignment,
		token.A_EXP:    p.parseShorthandAssignment,
		token.A_F_DIV:  p.parseShorthandAssignment,
		token.A_MINUS:  p.parseShorthandAssignment,
		token.A_MOD:    p.parseShorthandAssignment,
		token.A_OR:     p.parseShorthandAssignment,
		token.A_PLUS:   p.parseShorthandAssignment,
		token.A_Q_MARK: p.parseShorthandAssignment,
		token.A_SLASH:  p.parseShorthandAssignment,
		token.A_STAR:   p.parseShorthandAssignment,
		token.ASSIGN:   p.parseAssignExpression,
		token.DECLARE:  p.parseDeclareExpression,
		token.DOT:      p.parseDotExpression,
		token.COLON:    p.parseMethodCall,
		token.LSQUARE:  p.parseIndexExpression,
	}

	p.argTokens = []token.Type{
		token.PARAM,
	}

	for k := range p.prefixes {
		if !isBlacklisted(k) {
			p.argTokens = append(p.argTokens, k)
		}
	}

	p.next()
	p.next()

	return p
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peek.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.cur.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.lex()

	if p.peek.Type == token.ILLEGAL {
		p.err(
			fmt.Sprintf("illegal token found: `%s`", p.peek.Literal),
			p.peek.Start,
			p.peek.End,
		)
	}
}

func (p *Parser) Parse() ast.Program {
	prog := ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.curIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}

		p.next()
	}

	return prog
}
