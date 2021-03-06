package compiler

import (
	"errors"
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
	case *ast.Parameter:
		return c.compileParameter(node)
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
	case *ast.FunctionCall:
		return c.compileFnCall(node)
	case *ast.QualifiedFunctionCall:
		return c.compileQualifiedFnCall(node)
	case *ast.Argument:
		return c.CompileExpression(node.Value)
	case *ast.IndexExpression:
		return c.compileIndex(node)
	case *ast.DotExpression:
		return c.compileDot(node)
	case *ast.EmissionExpression:
		return c.compileEmission(node)
	case *ast.BlockLiteral:
		return c.compileBlockLiteral(node)
	default:
		return fmt.Errorf("compiler: compilation not yet implemented for %s", reflect.TypeOf(n))
	}
}

func (c *Compiler) compileNumber(node *ast.Number) error {
	var (
		obj        = &object.Number{Value: node.Value}
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileString(node *ast.String) error {
	var (
		obj        = &object.String{Value: node.Value}
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileBoolean(node *ast.Boolean) error {
	var (
		obj        = &object.Boolean{Value: node.Value}
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileChar(node *ast.Char) error {
	var (
		obj        = &object.Char{Value: rune(node.Value)}
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileNull(node *ast.Null) error {
	var (
		obj        = object.NullObj
		index, err = c.addConst(obj)
	)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}

func (c *Compiler) compileIdentifier(node *ast.Identifier) error {
	return c.compileName(node.Value)
}

func (c *Compiler) compileParameter(node *ast.Parameter) error {
	return c.compileName(node.Name)
}

func (c *Compiler) addName(name string) (rune, error) {
	for i, n := range c.Names {
		if name == n {
			return rune(i), nil
		}
	}

	c.Names = append(c.Names, name)
	index := len(c.Names) - 1

	return rune(index), nil
}

func (c *Compiler) compileName(name string) error {
	index, err := c.addName(name)
	if err != nil {
		return err
	}

	c.loadName(rune(index))

	return nil
}

func (c *Compiler) compileAssign(node *ast.AssignExpression) error {
	if err := c.CompileExpression(node.Value); err != nil {
		return err
	}

	if id, ok := node.Name.(*ast.Identifier); ok {
		index, err := c.addName(id.Value)
		if err != nil {
			return err
		}

		if index >= maxRune {
			return fmt.Errorf("compiler: name index %d greater than 0xFFFF (maximum uint16)", index)
		}

		low, high := runeToBytes(rune(index))
		c.push(bytecode.StoreName, high, low)
	} else if indexpr, ok := node.Name.(*ast.IndexExpression); ok {
		if err := c.CompileExpression(indexpr.Collection); err != nil {
			return err
		}

		if err := c.CompileExpression(indexpr.Index); err != nil {
			return err
		}

		c.push(bytecode.StoreField)
	} else if dotexpr, ok := node.Name.(*ast.DotExpression); ok {
		if err := c.CompileExpression(dotexpr.Left); err != nil {
			return err
		}

		if id, ok := dotexpr.Right.(*ast.Identifier); ok {
			obj := &object.String{Value: id.Value}

			c.Constants = append(c.Constants, obj)
			index := len(c.Constants) - 1

			if index >= maxRune {
				return fmt.Errorf("compiler: constant index %d greater than 0xFFFF (maximum uint16)", index)
			}

			low, high := runeToBytes(rune(index))

			c.push(bytecode.LoadConst, high, low)
		} else {
			return errors.New("compiler: expected an identifier to the right of a dot")
		}

		c.push(bytecode.StoreField)
	} else {
		return errors.New("compiler: can only assign to identfiers and field accessors")
	}

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

	c.push(op)

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

	c.push(op)

	return nil
}

func (c *Compiler) compileIf(node *ast.IfExpression) error {
	if err := c.CompileExpression(node.Condition); err != nil {
		return err
	}

	// JumpIfFalse (82) with 2 empty argument bytes
	c.push(bytecode.JumpIfFalse, 0, 0)
	condJump := len(c.Bytes) - 3

	if err := c.CompileStatement(node.Consequence); err != nil {
		return err
	}

	var skipJump int

	if node.Alternative != nil {
		// Jump past the alternative
		c.push(bytecode.Jump, 0, 0)
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

	c.push(bytecode.MakeArray, high, low)

	return nil
}

func (c *Compiler) compileTuple(node *ast.Tuple) error {
	for _, elem := range node.Value {
		if err := c.CompileExpression(elem); err != nil {
			return err
		}
	}

	low, high := runeToBytes(rune(len(node.Value)))

	c.push(bytecode.MakeTuple, high, low)

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

	c.push(bytecode.MakeMap, high, low)

	return nil
}

func (c *Compiler) compileFnCall(node *ast.FunctionCall) error {
	var ptn []string

	for _, item := range node.Pattern {
		if id, ok := item.(*ast.Identifier); ok {
			ptn = append(ptn, id.Value)
		} else {
			ptn = append(ptn, "$")
		}
	}

	str := strings.Join(ptn, " ")
	c.Patterns = append(c.Patterns, str)

	for _, item := range node.Pattern {
		if arg, ok := item.(*ast.Argument); ok {
			if err := c.CompileExpression(arg); err != nil {
				return err
			}
		}
	}

	low, high := runeToBytes(rune(len(c.Patterns) - 1))
	c.push(bytecode.PushFn, high, low, bytecode.CallFn)

	return nil
}

func (c *Compiler) compileQualifiedFnCall(node *ast.QualifiedFunctionCall) error {
	var ptn []string

	for _, item := range node.Pattern {
		if id, ok := item.(*ast.Identifier); ok {
			ptn = append(ptn, id.Value)
		} else {
			ptn = append(ptn, "$")
		}
	}

	str := strings.Join(ptn, " ")
	c.Patterns = append(c.Patterns, str)

	for _, item := range node.Pattern {
		if arg, ok := item.(*ast.Argument); ok {
			if err := c.CompileExpression(arg); err != nil {
				return err
			}
		}
	}

	if err := c.CompileExpression(node.Base); err != nil {
		return err
	}

	low, high := runeToBytes(rune(len(c.Patterns) - 1))
	c.push(bytecode.PushQualFn, high, low, bytecode.CallFn)

	return nil
}

func (c *Compiler) compileIndex(node *ast.IndexExpression) error {
	if err := c.CompileExpression(node.Collection); err != nil {
		return err
	}

	if err := c.CompileExpression(node.Index); err != nil {
		return err
	}

	c.push(bytecode.LoadField)

	return nil
}

func (c *Compiler) compileDot(node *ast.DotExpression) error {
	if err := c.CompileExpression(node.Left); err != nil {
		return err
	}

	if id, ok := node.Right.(*ast.Identifier); ok {
		index, err := c.addConst(&object.String{Value: id.Value})
		if err != nil {
			return err
		}

		c.loadConst(index)
	} else {
		return errors.New("compiler: expected an identifier to the right of a dot")
	}

	c.push(bytecode.LoadField)

	return nil
}

func (c *Compiler) compileEmission(node *ast.EmissionExpression) error {
	for _, item := range node.Items {
		if item.IsInstruction {
			var (
				ib     byte
				hasArg bool
			)

			for b, data := range bytecode.Instructions {
				if data.Name == item.Instruction {
					ib = b
					hasArg = data.HasArg
					break
				}
			}

			c.push(ib)

			if hasArg {
				low, high := runeToBytes(item.Argument)
				c.push(high, low)
			}
		} else {
			if err := c.CompileExpression(item.Exp); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Compiler) compileBlockLiteral(node *ast.BlockLiteral) error {
	fcomp := New()

	if err := fcomp.CompileStatement(node.Body); err != nil {
		return err
	}

	instructions, err := bytecode.Read(fcomp.Bytes)
	if err != nil {
		return err
	}

	obj := &object.Block{
		Params:    node.Params,
		Body:      instructions,
		Constants: fcomp.Constants,
		Names:     fcomp.Names,
		Patterns:  fcomp.Patterns,
	}

	index, err := c.addConst(obj)

	if err != nil {
		return err
	}

	c.loadConst(index)

	return nil
}
