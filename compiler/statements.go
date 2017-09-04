package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/pluto/ast"
)

// CompileStatement compiles an AST statement.
func (c *Compiler) CompileStatement(n ast.Statement) error {
	switch node := n.(type) {
	case *ast.ExpressionStatement:
		return c.CompileExpression(node.Expr)
	case *ast.BlockStatement:
		return c.compileBlockStatement(node)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(n))
	}
}

func (c *Compiler) compileBlockStatement(node *ast.BlockStatement) error {
	for _, stmt := range node.Statements {
		if err := c.CompileStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}
