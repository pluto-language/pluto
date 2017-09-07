package compiler

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

// CompileExpression compiles an AST expression.
func (c *Compiler) CompileExpression(n ast.Expression) error {
	switch node := n.(type) {
	case *ast.InfixExpression:
		return c.compileInfix(node)
	case *ast.PrefixExpression:
		return c.compilePrefix(node)
	case *ast.Number:
		return c.compileNumber(node)
	case *ast.String:
		return c.compileString(node)
	case *ast.Boolean:
		return c.compileBoolean(node)
	case *ast.Char:
		return c.compileChar(node)
	case *ast.Null:
		return c.compileNull(node)
	case *ast.Identifier:
		return c.compileIdentifier(node)
	case *ast.Array:
		return c.compileArray(node)
	case *ast.Tuple:
		return c.compileTuple(node)
	case *ast.Map:
		return c.compileMap(node)
	case *ast.AssignExpression:
		return c.compileAssign(node)
	case *ast.IfExpression:
		return c.compileIf(node)
	case *ast.WhileLoop:
		return c.compileWhile(node)
	case *ast.FunctionCall:
		return c.compileFnCall(node)
	case *ast.Argument:
		return c.CompileExpression(node.Value)
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

	c.Bytes = append(c.Bytes, bytecode.LoadConst, high, low)

	return nil
}

func (c *Compiler) compileString(node *ast.String) error {
	obj := &object.String{Value: node.Value}
	c.Constants = append(c.Constants, obj)
	index := len(c.Constants) - 1

	if index >= 1<<16 {
		return fmt.Errorf("compiler: constant index %d greater than 1 << 16 (maximum uint16)", index)
	}

	low, high := runeToBytes(rune(index))

	c.Bytes = append(c.Bytes, bytecode.LoadConst, high, low)

	return nil
}

func (c *Compiler) compileBoolean(node *ast.Boolean) error {
	obj := &object.Boolean{Value: node.Value}
	c.Constants = append(c.Constants, obj)
	index := len(c.Constants) - 1

	if index >= 1<<16 {
		return fmt.Errorf("compiler: constant index %d greater than 1 << 16 (maximum uint16)", index)
	}

	low, high := runeToBytes(rune(index))

	c.Bytes = append(c.Bytes, bytecode.LoadConst, high, low)

	return nil
}

func (c *Compiler) compileChar(node *ast.Char) error {
	obj := &object.Char{Value: rune(node.Value)}
	c.Constants = append(c.Constants, obj)
	index := len(c.Constants) - 1

	if index >= 1<<16 {
		return fmt.Errorf("compiler: constant index %d greater than 1 << 16 (maximum uint16)", index)
	}

	low, high := runeToBytes(rune(index))

	c.Bytes = append(c.Bytes, bytecode.LoadConst, high, low)

	return nil
}

func (c *Compiler) compileNull(node *ast.Null) error {
	obj := object.NullObj
	c.Constants = append(c.Constants, obj)
	index := len(c.Constants) - 1

	if index >= 1<<16 {
		return fmt.Errorf("compiler: constant index %d greater than 1 << 16 (maximum uint16)", index)
	}

	low, high := runeToBytes(rune(index))

	c.Bytes = append(c.Bytes, bytecode.LoadConst, high, low)

	return nil
}

func (c *Compiler) compileIdentifier(node *ast.Identifier) error {
	var index int

	for i, name := range c.Names {
		if name == node.Value {
			index = i
			goto found
		}
	}

	// These two lines are executed if the name isn't found
	c.Names = append(c.Names, node.Value)
	index = len(c.Names) - 1

found:
	low, high := runeToBytes(rune(index))

	c.Bytes = append(c.Bytes, bytecode.LoadName, high, low)

	return nil
}

func (c *Compiler) compileAssign(node *ast.AssignExpression) error {
	if err := c.CompileExpression(node.Value); err != nil {
		return err
	}

	c.Names = append(c.Names, node.Name.(*ast.Identifier).Value)
	index := len(c.Names) - 1

	if index >= 1<<16 {
		return fmt.Errorf("compiler: name index %d greater than 1 << 16 (maximum uint16)", index)
	}

	low, high := runeToBytes(rune(index))

	c.Bytes = append(c.Bytes, bytecode.StoreName, high, low)

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

	op, ok := map[string]byte{
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
		"!=": bytecode.BinaryNotEqual,
		"<":  bytecode.BinaryLessThan,
		">":  bytecode.BinaryMoreThan,
		"<=": bytecode.BinaryLessEq,
		">=": bytecode.BinaryMoreEq,
	}[node.Operator]

	if !ok {
		return fmt.Errorf("compiler: operator %s not yet implemented", node.Operator)
	}

	c.Bytes = append(c.Bytes, op)

	return nil
}

func (c *Compiler) compilePrefix(node *ast.PrefixExpression) error {
	if err := c.CompileExpression(node.Right); err != nil {
		return err
	}

	op := map[string]byte{
		"+": bytecode.UnaryNoOp,
		"-": bytecode.UnaryNegate,
		"!": bytecode.UnaryInvert,
	}[node.Operator]

	c.Bytes = append(c.Bytes, op)

	return nil
}

func (c *Compiler) compileIf(node *ast.IfExpression) error {
	if err := c.CompileExpression(node.Condition); err != nil {
		return err
	}

	// JumpIfFalse (82) with 2 empty argument bytes
	c.Bytes = append(c.Bytes, bytecode.JumpIfFalse, 0, 0)
	condJump := len(c.Bytes) - 3

	if err := c.CompileStatement(node.Consequence); err != nil {
		return err
	}

	var skipJump int

	if node.Alternative != nil {
		// Jump past the alternative
		c.Bytes = append(c.Bytes, bytecode.Jump, 0, 0)
		skipJump = len(c.Bytes) - 3
	}

	// Set the jump target after the conditional
	condIndex := rune(len(c.Bytes))
	low, high := runeToBytes(condIndex)
	c.Bytes[condJump+1] = high
	c.Bytes[condJump+2] = low

	if node.Alternative != nil {
		if err := c.CompileStatement(node.Alternative); err != nil {
			return err
		}

		// Set the jump target after the conditional
		skipIndex := rune(len(c.Bytes))
		low, high = runeToBytes(skipIndex)
		c.Bytes[skipJump+1] = high
		c.Bytes[skipJump+2] = low
	}

	return nil
}

func (c *Compiler) compileArray(node *ast.Array) error {
	for _, elem := range node.Elements {
		if err := c.CompileExpression(elem); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Elements)))

	c.Bytes = append(c.Bytes, bytecode.MakeArray, high, low)

	return nil
}

func (c *Compiler) compileTuple(node *ast.Tuple) error {
	for _, elem := range node.Value {
		if err := c.CompileExpression(elem); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Value)))

	c.Bytes = append(c.Bytes, bytecode.MakeTuple, high, low)

	return nil
}

func (c *Compiler) compileMap(node *ast.Map) error {
	for key, val := range node.Pairs {
		if err := c.CompileExpression(key); err != nil {
			return err
		}

		if err := c.CompileExpression(val); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Pairs)))

	c.Bytes = append(c.Bytes, bytecode.MakeMap, high, low)

	return nil
}

func (c *Compiler) compileWhile(node *ast.WhileLoop) error {
	// Jump here to go to the next iteration
	start := len(c.Bytes) - 1

	if err := c.CompileExpression(node.Condition); err != nil {
		return err
	}

	// An empty jump to the end of the loop
	c.Bytes = append(c.Bytes, bytecode.JumpIfFalse, 0, 0)
	skipJump := len(c.Bytes) - 3

	// Compile the loop's body
	if err := c.CompileStatement(node.Body); err != nil {
		return err
	}

	// After the body, jump back to the beginning of the loop
	low, high := runeToBytes(rune(start))
	c.Bytes = append(c.Bytes, bytecode.Jump, high, low)

	// If the condition isn't met, jump to the end of the loop
	skipIndex := rune(len(c.Bytes))
	low, high = runeToBytes(skipIndex)
	c.Bytes[skipJump+1] = high
	c.Bytes[skipJump+2] = low

	return nil
}

func (c *Compiler) compileFnCall(node *ast.FunctionCall) error {
	var pstring bytes.Buffer

	for _, item := range node.Pattern {
		if id, ok := item.(*ast.Identifier); ok {
			pstring.WriteString(id.Value + " ")
		} else {
			pstring.WriteString("$ ")
		}
	}

	str := strings.TrimSpace(pstring.String())
	c.Patterns = append(c.Patterns, str)

	for _, item := range node.Pattern {
		if arg, ok := item.(*ast.Argument); ok {
			if err := c.CompileExpression(arg); err != nil {
				return err
			}
		}
	}

	low, high := runeToBytes(rune(len(c.Patterns) - 1))
	c.Bytes = append(c.Bytes, bytecode.Call, high, low)

	return nil
}
