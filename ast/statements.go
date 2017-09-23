package ast

import "github.com/Zac-Garby/pluto/token"

type (
	// ExpressionStatement is an expression which acts a statement
	ExpressionStatement struct {
		Tok  token.Token
		Expr Expression
	}

	// BlockStatement is a list of statements
	BlockStatement struct {
		Tok        token.Token
		Statements []Statement
	}

	// FunctionDefinition defines a function
	FunctionDefinition struct {
		Tok     token.Token
		Pattern []Expression
		Body    Statement
	}

	// ReturnStatement returns an expression from a BlockStatement
	ReturnStatement struct {
		Tok   token.Token
		Value Expression
	}

	// NextStatement goes to the next iteration of a loop
	NextStatement struct {
		Tok token.Token
	}

	// BreakStatement breaks a loop
	BreakStatement struct {
		Tok token.Token
	}

	// UseStatement imports a package into the current scope
	UseStatement struct {
		Tok     token.Token
		Package string
	}

	// WhileLoop executes Body while Condition holds true
	WhileLoop struct {
		Tok       token.Token
		Condition Expression
		Body      Statement
	}

	// ForLoop executes Body for each element in a collection
	ForLoop struct {
		Tok token.Token

		// for (Init; Condition; Increment) { Body }
		Init      Expression
		Condition Expression
		Increment Expression
		Body      Statement
	}
)

// Stmt tells the compiler this node is a statement
func (n ExpressionStatement) Stmt() {}

//Token returns this node's token
func (n ExpressionStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n BlockStatement) Stmt() {}

//Token returns this node's token
func (n BlockStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n FunctionDefinition) Stmt() {}

//Token returns this node's token
func (n FunctionDefinition) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n ReturnStatement) Stmt() {}

//Token returns this node's token
func (n ReturnStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n NextStatement) Stmt() {}

//Token returns this node's token
func (n NextStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n BreakStatement) Stmt() {}

//Token returns this node's token
func (n BreakStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n UseStatement) Stmt() {}

//Token returns this node's token
func (n UseStatement) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n WhileLoop) Stmt() {}

//Token returns this node's token
func (n WhileLoop) Token() token.Token { return n.Tok }

// Stmt tells the compiler this node is a statement
func (n ForLoop) Stmt() {}

//Token returns this node's token
func (n ForLoop) Token() token.Token { return n.Tok }
