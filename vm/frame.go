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

		if f.vm.Error != nil {
			break
		}

		f.offset++
	}
}

func (f *Frame) doInstruction(i bytecode.Instruction) {
	e, ok := Effectors()[i.Code]
	if !ok {
		f.vm.Error = fmt.Errorf("evaluation: bytecode instruction %s not implemented", i.Name)
		return
	}

	e(f, i)
}

func (f *Frame) getName(arg rune) (string, bool) {
	index := int(arg)

	if index < len(f.locals.Names) {
		name := f.locals.Names[index]
		return name, true
	} else if f.previous != nil && index < len(f.previous.locals.Names) {
		name := f.previous.locals.Names[index]
		return name, true
	}

	return "", false
}

func (f *Frame) searchName(name string) (object.Object, bool) {
	if val, ok := f.locals.Data[name]; ok {
		return val, true
	} else if f.previous != nil {
		if val, ok := f.previous.locals.Data[name]; ok {
			return val, true
		}
	}

	return nil, false
}
