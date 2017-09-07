package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

// CompileStatement compiles an AST statement.
func (c *Compiler) CompileStatement(n ast.Statement) error {
	switch node := n.(type) {
	case *ast.ExpressionStatement:
		return c.CompileExpression(node.Expr)
	case *ast.BlockStatement:
		return c.compileBlockStatement(node)
	case *ast.FunctionDefinition:
		return c.compileFunctionDefinition(node)
	case *ast.ReturnStatement:
		return c.compileReturnStatement(node)
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

func (c *Compiler) compileFunctionDefinition(node *ast.FunctionDefinition) error {
	fcomp := New()

	if err := fcomp.CompileStatement(node.Body); err != nil {
		return err
	}

	instructions, err := bytecode.Read(fcomp.Bytes)
	if err != nil {
		return err
	}

	fn := object.Function{
		Pattern:   node.Pattern,
		Body:      instructions,
		Constants: fcomp.Constants,
		Names:     fcomp.Names,
	}

	c.Functions = append(c.Functions, fn)

	return nil
}

func (c *Compiler) compileReturnStatement(node *ast.ReturnStatement) error {
	if node.Value != nil {
		if err := c.CompileExpression(node.Value); err != nil {
			return err
		}
	}

	c.Bytes = append(c.Bytes, bytecode.Return)

	return nil
}
