package vm

import (
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

// VirtualMachine is the base struct
// for the Pluto VM. It stores information
// such as the call stack and the last
// error thrown.
type VirtualMachine struct {
	frames      []*Frame
	frame       *Frame
	returnValue *object.Object
	lastError   *object.Object
}

// Run executes the supplied bytecode
func (vm *VirtualMachine) Run(code bytecode.Code, globals, locals map[string]*object.Object) {
	frame := &Frame{
		code:            code,
		locals:          locals,
		globals:         globals,
		lastInstruction: 0,
		stack:           []*object.Object{},
	}
}
