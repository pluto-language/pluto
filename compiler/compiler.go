package compiler

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/object"
)

// Compiler compiles an AST into bytecode
type Compiler struct {
	// Generated code:
	Bytes []byte

	// Data:
	Constants       []object.Object
	Functions       []object.Function
	Names, Patterns []string
}

// New instantiates a new Compiler, and allocates
// memory for the members.
func New() Compiler {
	return Compiler{
		Bytes:     make([]byte, 0),
		Constants: make([]object.Object, 0),
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
