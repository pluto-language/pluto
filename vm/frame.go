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

	locals    Store // the local namespace
	globals   Store // the global namespace
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

func (f *Frame) getName(arg rune) (string, bool) {
	if name, ok := f.locals.Names[arg]; ok {
		return name, true
	} else if name, ok := f.globals.Names[arg]; ok {
		return name, true
	}

	return "", false
}

func (f *Frame) searchName(name string) (object.Object, bool) {
	if val, ok := f.locals.Data[name]; ok {
		return val, true
	} else if val, ok := f.globals.Data[name]; ok {
		return val, true
	}

	return nil, false
}
