package vm

import (
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/evaluation/object"
)

// Frame is a virtual machine frame. A frame is
// created for each function call, and stores
// things like the variables in that scope and
// the stack.
type Frame struct {
	previous        *Frame
	code            bytecode.Code
	lastInstruction int

	local, global    map[string]*object.Object
	constants, stack []object.Object
}
