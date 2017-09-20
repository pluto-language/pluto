package compiler

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

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
	case *ast.WhileLoop:
		return c.compileWhile(node)
	case *ast.NextStatement:
		return c.compileNext(node)
	case *ast.BreakStatement:
		return c.compileBreak(node)
	case *ast.UseStatement:
		return c.compileUse(node)
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
		Patterns:  fcomp.Patterns,
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

	c.push(bytecode.Return)

	return nil
}

func (c *Compiler) compileWhile(node *ast.WhileLoop) error {
	c.push(bytecode.LoopStart)

	// Jump here to go to the next iteration
	start := len(c.Bytes) - 1

	if err := c.CompileExpression(node.Condition); err != nil {
		return err
	}

	// An empty jump to the end of the loop
	c.push(bytecode.JumpIfFalse, 0, 0)
	skipJump := len(c.Bytes) - 3

	// Compile the loop's body
	if err := c.CompileStatement(node.Body); err != nil {
		return err
	}

	// After the body, jump back to the beginning of the loop
	low, high := runeToBytes(rune(start))
	c.push(bytecode.Jump, high, low)

	// If the condition isn't met, jump to the end of the loop
	skipIndex := rune(len(c.Bytes))
	low, high = runeToBytes(skipIndex)
	c.Bytes[skipJump+1] = high
	c.Bytes[skipJump+2] = low

	c.push(bytecode.LoopEnd)

	return nil
}

func (c *Compiler) compileNext(node *ast.NextStatement) error {
	c.push(bytecode.Next)

	return nil
}

func (c *Compiler) compileBreak(node *ast.BreakStatement) error {
	c.push(bytecode.Break)

	return nil
}

func (c *Compiler) compileUse(node *ast.UseStatement) error {
	pkg := node.Package

	if strings.HasPrefix(node.Package, "./") {
		dir, _ := filepath.Split(node.Tok.Start.File)
		pkg = filepath.Join(dir, pkg)
	}

	obj := &object.String{Value: pkg}
	c.Constants = append(c.Constants, obj)
	index := len(c.Constants) - 1

	if index >= maxRune {
		return fmt.Errorf("compiler: constant index %d greater than 0xFFFF (maximum uint16)", index)
	}

	low, high := runeToBytes(rune(index))

	c.push(bytecode.Use, high, low)

	return nil
}
