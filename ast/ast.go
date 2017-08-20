package ast

import (
	"github.com/Zac-Garby/pluto/token"
)

type Node interface {
	Token() token.Token
}

type Statement interface {
	Node
	Stmt()
}

type Expression interface {
	Node
	Expr()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Tree() string {
	str := ""

	for _, stmt := range p.Statements {
		str += Tree(stmt, 0, "") + "\n"
	}

	return str
}
