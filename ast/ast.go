package ast

import (
	"github.com/Zac-Garby/pluto/token"
)

// Node is the base AST node struct
type Node interface {
	Token() token.Token
}

// Statement is a statement AST node
type Statement interface {
	Node
	Stmt()
}

// Expression is an expression AST node
type Expression interface {
	Node
	Expr()
}

// Program is a program, containing
// a list of statements
type Program struct {
	Statements []Statement
}

// Tree returns a tree representation of
// a program.
func (p *Program) Tree() string {
	str := ""

	for _, stmt := range p.Statements {
		str += Tree(stmt, 0, "") + "\n"
	}

	return str
}
