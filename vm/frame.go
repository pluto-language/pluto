package vm

import (
	"fmt"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

// Frame is a virtual machine frame. A frame is
// created for each function call, and stores
// things like the variables in that scope and
// the stack.
type Frame struct {
	previous *Frame          // the previous frame
	code     bytecode.Code   // the parsed bytecode
	offset   int             // the current instruction index
	vm       *VirtualMachine // the frame's virtual machine

	locals    store // the local namespace
	globals   store // the global namespace
	stack     stack // the object stack
	constants []object.Object
}

func (f *Frame) execute() {
	for {
		if f.offset >= len(f.code) {
			return
		}

		instruction := f.code[f.offset]

		f.doInstruction(instruction)

		f.offset++
	}
}

func (f *Frame) doInstruction(i bytecode.Instruction) {
	e, ok := effectors[i.Code]
	if !ok {
		f.vm.lastError = fmt.Errorf("evaluation: bytecode instruction %s not implemented", i.Name)
		return
	}

	e(f, i)
}
