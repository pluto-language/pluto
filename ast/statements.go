package ast

import "github.com/Zac-Garby/pluto/token"

/* ExpressionStatement */
type ExpressionStatement struct {
	Tok  token.Token
	Expr Expression
}

func (_ ExpressionStatement) Stmt()              {}
func (n ExpressionStatement) Token() token.Token { return n.Tok }

/* BlockStatement */
type BlockStatement struct {
	Tok        token.Token
	Statements []Statement
}

func (_ BlockStatement) Stmt()              {}
func (n BlockStatement) Token() token.Token { return n.Tok }

/* FunctionDefinition */
type FunctionDefinition struct {
	Tok     token.Token
	Pattern []Expression
	Body    Statement
}

func (_ FunctionDefinition) Stmt()              {}
func (n FunctionDefinition) Token() token.Token { return n.Tok }

/* InitDefinition */
type InitDefinition struct {
	Tok     token.Token
	Pattern []Expression
	Body    Statement
}

func (_ InitDefinition) Stmt()              {}
func (n InitDefinition) Token() token.Token { return n.Tok }

/* ReturnStatement */
type ReturnStatement struct {
	Tok   token.Token
	Value Expression
}

func (_ ReturnStatement) Stmt()              {}
func (n ReturnStatement) Token() token.Token { return n.Tok }

/* NextStatement */
type NextStatement struct {
	Tok token.Token
}

func (_ NextStatement) Stmt()              {}
func (n NextStatement) Token() token.Token { return n.Tok }

/* BreakStatement */
type BreakStatement struct {
	Tok token.Token
}

func (_ BreakStatement) Stmt()              {}
func (n BreakStatement) Token() token.Token { return n.Tok }

/* ClassStatement */
type ClassStatement struct {
	Tok          token.Token
	Name, Parent Expression
	Methods      []Statement
}

func (_ ClassStatement) Stmt()              {}
func (n ClassStatement) Token() token.Token { return n.Tok }

/* ImportStatement */
type ImportStatement struct {
	Tok     token.Token
	Package string
}

func (_ ImportStatement) Stmt()              {}
func (n ImportStatement) Token() token.Token { return n.Tok }

/* UseStatement */
type UseStatement struct {
	Tok     token.Token
	Package string
}

func (_ UseStatement) Stmt()              {}
func (n UseStatement) Token() token.Token { return n.Tok }
