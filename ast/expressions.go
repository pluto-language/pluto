package ast

import "github.com/Zac-Garby/pluto/token"

type (
	// Identifier is an identifier
	Identifier struct {
		Tok   token.Token
		Value string
	}

	// Number is a number literal
	Number struct {
		Tok   token.Token
		Value float64
	}

	// Boolean is a boolean literal
	Boolean struct {
		Tok   token.Token
		Value bool
	}

	// String is a string literal
	String struct {
		Tok   token.Token
		Value string
	}

	// Char is a character literal
	Char struct {
		Tok   token.Token
		Value byte
	}

	// Tuple is a tuple literal
	Tuple struct {
		Tok   token.Token
		Value []Expression
	}

	// Array is an array literal
	Array struct {
		Tok      token.Token
		Elements []Expression
	}

	// Map is a map literal
	Map struct {
		Tok   token.Token
		Pairs map[Expression]Expression
	}

	// BlockLiteral is a block literal
	BlockLiteral struct {
		Tok    token.Token
		Body   Statement
		Params []Expression
	}

	// Null is the null literal
	Null struct {
		Tok token.Token
	}

	// AssignExpression assigns an expression to a name
	AssignExpression struct {
		Tok         token.Token
		Name, Value Expression
	}

	// DeclareExpression declares a name as an expression
	DeclareExpression struct {
		Tok         token.Token
		Name, Value Expression
	}

	// PrefixExpression is a prefix operator expression
	PrefixExpression struct {
		Tok      token.Token
		Operator string
		Right    Expression
	}

	// InfixExpression is an infix operator expression
	InfixExpression struct {
		Tok         token.Token
		Operator    string
		Left, Right Expression
	}

	// DotExpression gets a value from a container
	DotExpression struct {
		Tok         token.Token
		Left, Right Expression
	}

	// IndexExpression gets a value from a collection
	IndexExpression struct {
		Tok               token.Token
		Collection, Index Expression
	}

	// Parameter is shorthand for putting an identifier into a pattern call
	Parameter struct {
		Tok  token.Token
		Name string
	}

	// Argument is an argument to a function
	Argument struct {
		Tok   token.Token
		Value Expression
	}

	// FunctionCall calls a function
	FunctionCall struct {
		Tok     token.Token
		Pattern []Expression
	}

	// QualifiedFunctionCall calls a function from a package
	QualifiedFunctionCall struct {
		Tok     token.Token
		Package Expression
		Pattern []Expression
	}

	// IfExpression executes Consequence or Alternative based on Condition
	IfExpression struct {
		Tok                      token.Token
		Condition                Expression
		Consequence, Alternative Statement
	}

	// Arm is a match-arm
	Arm struct {
		Exprs []Expression
		Body  Statement
	}

	// MatchExpression matches an expression against a list of arms
	MatchExpression struct {
		Tok  token.Token
		Exp  Expression
		Arms []Arm
	}

	// WhileLoop executes Body while Condition holds true
	WhileLoop struct {
		Tok       token.Token
		Condition Expression
		Body      Statement
	}

	// ForLoop executes Body for each element in a collection
	ForLoop struct {
		Tok             token.Token
		Var, Collection Expression
		Body            Statement
	}

	// TryExpression trys to execute Body, and matches Arms if it throws an error
	TryExpression struct {
		Tok     token.Token
		Body    Statement
		ErrName Expression
		Arms    []Arm
	}

	// MethodCall calls a method on an instance
	MethodCall struct {
		Tok      token.Token
		Instance Expression
		Pattern  []Expression
	}
)

// Expr tells the compiler this node is an expression
func (n Identifier) Expr() {}

// Token returns the node's token
func (n Identifier) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Number) Expr() {}

// Token returns the node's token
func (n Number) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Boolean) Expr() {}

// Token returns the node's token
func (n Boolean) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n String) Expr() {}

// Token returns the node's token
func (n String) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Char) Expr() {}

// Token returns the node's token
func (n Char) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Tuple) Expr() {}

// Token returns the node's token
func (n Tuple) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Array) Expr() {}

// Token returns the node's token
func (n Array) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Map) Expr() {}

// Token returns the node's token
func (n Map) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n BlockLiteral) Expr() {}

// Token returns the node's token
func (n BlockLiteral) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Null) Expr() {}

// Token returns the node's token
func (n Null) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n AssignExpression) Expr() {}

// Token returns the node's token
func (n AssignExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n DeclareExpression) Expr() {}

// Token returns the node's token
func (n DeclareExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n PrefixExpression) Expr() {}

// Token returns the node's token
func (n PrefixExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n InfixExpression) Expr() {}

// Token returns the node's token
func (n InfixExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n DotExpression) Expr() {}

// Token returns the node's token
func (n DotExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n IndexExpression) Expr() {}

// Token returns the node's token
func (n IndexExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Parameter) Expr() {}

// Token returns the node's token
func (n Parameter) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n Argument) Expr() {}

// Token returns the node's token
func (n Argument) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n FunctionCall) Expr() {}

// Token returns the node's token
func (n FunctionCall) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n QualifiedFunctionCall) Expr() {}

// Token returns the node's token
func (n QualifiedFunctionCall) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n IfExpression) Expr() {}

// Token returns the node's token
func (n IfExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n MatchExpression) Expr() {}

// Token returns the node's token
func (n MatchExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n WhileLoop) Expr() {}

// Token returns the node's token
func (n WhileLoop) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n ForLoop) Expr() {}

// Token returns the node's token
func (n ForLoop) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n TryExpression) Expr() {}

// Token returns the node's token
func (n TryExpression) Token() token.Token { return n.Tok }

// Expr tells the compiler this node is an expression
func (n MethodCall) Expr() {}

// Token returns the node's token
func (n MethodCall) Token() token.Token { return n.Tok }
