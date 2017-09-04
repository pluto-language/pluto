package compiler

import (
	"fmt"
	"reflect"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

// Compiler compiles an AST into bytecode
type Compiler struct {
	Bytes     []byte
	Constants []object.Object
}

// CompileExpression compiles an AST expression.
func (c *Compiler) CompileExpression(n ast.Expression) error {
	switch node := n.(type) {
	case *ast.InfixExpression:
		return c.compileInfix(node)
	case *ast.Number:
		return c.compileNumber(node)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(n))
	}
}

func (c *Compiler) compileNumber(node *ast.Number) error {
	obj := &object.Number{Value: node.Value}
	c.Constants = append(c.Constants, obj)
	index := len(c.Constants) - 1

	low, high := runeToBytes(rune(index))

	c.Bytes = append(c.Bytes, 10, high, low)

	return nil
}

func (c *Compiler) compileInfix(node *ast.InfixExpression) error {
	left, right := node.Left, node.Right

	if err := c.CompileExpression(left); err != nil {
		return err
	}

	if err := c.CompileExpression(right); err != nil {
		return err
	}

	var op byte

	switch node.Operator {
	case "+":
		op = bytecode.BinaryAdd
	case "-":
		op = bytecode.BinarySubtract
	}

	c.Bytes = append(c.Bytes, op)

	return nil
}
