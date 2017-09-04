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

// New instantiates a new Compiler, and allocates
// memory for the members.
func New() Compiler {
	return Compiler{
		Bytes:     make([]byte, 0),
		Constants: make([]object.Object, 16),
	}
}

// CompileProgram compiles a complete parsed program.
func (c *Compiler) CompileProgram(p ast.Program) error {
	for _, stmt := range p.Statements {
		if err := c.CompileStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

// CompileStatement compiles an AST statement.
func (c *Compiler) CompileStatement(n ast.Statement) error {
	switch node := n.(type) {
	case *ast.ExpressionStatement:
		return c.CompileExpression(node.Expr)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(n))
	}
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

	if index >= 1<<16 {
		return fmt.Errorf("compiler: constant index %d greater than 1 << 16 (maximum uint16)", index)
	}

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

	op := map[string]byte{
		"+":  bytecode.BinaryAdd,
		"-":  bytecode.BinarySubtract,
		"*":  bytecode.BinaryMultiply,
		"/":  bytecode.BinaryDivide,
		"**": bytecode.BinaryExponent,
		"//": bytecode.BinaryFloorDiv,
		"%":  bytecode.BinaryFloorDiv,
		"||": bytecode.BinaryOr,
		"&&": bytecode.BinaryAnd,
		"|":  bytecode.BinaryBitOr,
		"&":  bytecode.BinaryBitAnd,
		"==": bytecode.BinaryEquals,
	}[node.Operator]

	c.Bytes = append(c.Bytes, op)

	return nil
}
