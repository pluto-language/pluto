package ast

import "github.com/Zac-Garby/pluto/token"

/* Identifier */
type Identifier struct {
	Tok   token.Token
	Value string
}

func (_ Identifier) Expr()              {}
func (n Identifier) Token() token.Token { return n.Tok }

/* Number */
type Number struct {
	Tok   token.Token
	Value float64
}

func (_ Number) Expr()              {}
func (n Number) Token() token.Token { return n.Tok }

/* Boolean */
type Boolean struct {
	Tok   token.Token
	Value bool
}

func (_ Boolean) Expr()              {}
func (n Boolean) Token() token.Token { return n.Tok }

/* String */
type String struct {
	Tok   token.Token
	Value string
}

func (_ String) Expr()              {}
func (n String) Token() token.Token { return n.Tok }

/* Char */
type Char struct {
	Tok   token.Token
	Value byte
}

func (_ Char) Expr()              {}
func (n Char) Token() token.Token { return n.Tok }

/* Tuple */
type Tuple struct {
	Tok   token.Token
	Value []Expression
}

func (_ Tuple) Expr()              {}
func (n Tuple) Token() token.Token { return n.Tok }

/* Array */
type Array struct {
	Tok      token.Token
	Elements []Expression
}

func (_ Array) Expr()              {}
func (n Array) Token() token.Token { return n.Tok }

/* Map */
type Map struct {
	Tok   token.Token
	Pairs map[Expression]Expression
}

func (_ Map) Expr()              {}
func (n Map) Token() token.Token { return n.Tok }

/* BlockLiteral */
type BlockLiteral struct {
	Tok    token.Token
	Body   Statement
	Params []Expression
}

func (_ BlockLiteral) Expr()              {}
func (n BlockLiteral) Token() token.Token { return n.Tok }

/* Null */
type Null struct {
	Tok token.Token
}

func (_ Null) Expr()              {}
func (n Null) Token() token.Token { return n.Tok }

/* AssignExpression */
type AssignExpression struct {
	Tok         token.Token
	Name, Value Expression
}

func (_ AssignExpression) Expr()              {}
func (n AssignExpression) Token() token.Token { return n.Tok }

/* DeclareExpression */
type DeclareExpression struct {
	Tok         token.Token
	Name, Value Expression
}

func (_ DeclareExpression) Expr()              {}
func (n DeclareExpression) Token() token.Token { return n.Tok }

/* PrefixExpression */
type PrefixExpression struct {
	Tok      token.Token
	Operator string
	Right    Expression
}

func (_ PrefixExpression) Expr()              {}
func (n PrefixExpression) Token() token.Token { return n.Tok }

/* InfixExpression */
type InfixExpression struct {
	Tok         token.Token
	Operator    string
	Left, Right Expression
}

func (_ InfixExpression) Expr()              {}
func (n InfixExpression) Token() token.Token { return n.Tok }

/* DotExpression */
type DotExpression struct {
	Tok         token.Token
	Left, Right Expression
}

func (_ DotExpression) Expr()              {}
func (n DotExpression) Token() token.Token { return n.Tok }

/* IndexExpression */
type IndexExpression struct {
	Tok               token.Token
	Collection, Index Expression
}

func (_ IndexExpression) Expr()              {}
func (i IndexExpression) Token() token.Token { return i.Tok }

/* Parameter */
type Parameter struct {
	Tok  token.Token
	Name string
}

func (_ Parameter) Expr()              {}
func (n Parameter) Token() token.Token { return n.Tok }

/* Argument */
type Argument struct {
	Tok   token.Token
	Value Expression
}

func (_ Argument) Expr()              {}
func (n Argument) Token() token.Token { return n.Tok }

/* FunctionCall */
type FunctionCall struct {
	Tok     token.Token
	Pattern []Expression
}

func (_ FunctionCall) Expr()              {}
func (n FunctionCall) Token() token.Token { return n.Tok }

/* QualifiedFunctionCall */
type QualifiedFunctionCall struct {
	Tok     token.Token
	Package Expression
	Pattern []Expression
}

func (_ QualifiedFunctionCall) Expr()              {}
func (n QualifiedFunctionCall) Token() token.Token { return n.Tok }

/* IfExpression */
type IfExpression struct {
	Tok                      token.Token
	Condition                Expression
	Consequence, Alternative Statement
}

func (_ IfExpression) Expr()              {}
func (n IfExpression) Token() token.Token { return n.Tok }

type Arm struct {
	Exprs []Expression
	Body  Statement
}

/* MatchExpression */
type MatchExpression struct {
	Tok  token.Token
	Exp  Expression
	Arms []Arm
}

func (_ MatchExpression) Expr()              {}
func (n MatchExpression) Token() token.Token { return n.Tok }

/* WhileLoop */
type WhileLoop struct {
	Tok       token.Token
	Condition Expression
	Body      Statement
}

func (_ WhileLoop) Expr()              {}
func (n WhileLoop) Token() token.Token { return n.Tok }

/* ForLoop */
type ForLoop struct {
	Tok             token.Token
	Var, Collection Expression
	Body            Statement
}

func (_ ForLoop) Expr()              {}
func (n ForLoop) Token() token.Token { return n.Tok }

/* TryExpression */
type TryExpression struct {
	Tok     token.Token
	Body    Statement
	ErrName Expression
	Arms    []Arm
}

func (_ TryExpression) Expr()              {}
func (n TryExpression) Token() token.Token { return n.Tok }

/* MethodCall */
type MethodCall struct {
	Tok      token.Token
	Instance Expression
	Pattern  []Expression
}

func (_ MethodCall) Expr()              {}
func (n MethodCall) Token() token.Token { return n.Tok }
