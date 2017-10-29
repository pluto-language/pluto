package vm

import (
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
	"github.com/Zac-Garby/pluto/store"
)

// VirtualMachine is the base struct
// for the Pluto VM. It stores information
// such as the call stack and the last
// error thrown.
type VirtualMachine struct {
	frames      []*Frame
	frame       *Frame
	returnValue object.Object
	Error       *Error
}

// New returns a new virtual machine
func New() *VirtualMachine {
	return &VirtualMachine{
		frames:      make([]*Frame, 0),
		returnValue: nil,
		Error:       nil,
	}
}

// Run executes the supplied bytecode
func (vm *VirtualMachine) Run(code bytecode.Code, locals *store.Store, constants []object.Object, usePrelude bool) {
	frame := vm.makeFrame(code, store.New(), locals, constants)

	if usePrelude {
		frame.Use("std/prelude/*.pluto")
	}

	vm.pushFrame(frame)
	vm.runFrame(frame)
}

// RunDefault executes the bytecode with
// empty globals and locals
func (vm *VirtualMachine) RunDefault(code bytecode.Code, constants []object.Object) {
	vm.Run(code, store.New(), constants, false)
}

func (vm *VirtualMachine) makeFrame(code bytecode.Code, args, locals *store.Store, constants []object.Object) *Frame {
	frame := &Frame{
		code:      code,
		locals:    locals,
		offset:    0,
		stack:     newStack(),
		constants: constants,
		vm:        vm,
	}

	for k, v := range args.Names {
		locals.Names[k] = v
		locals.Define(v, args.GetName(v), true)
	}

	return frame
}

func (vm *VirtualMachine) pushFrame(frame *Frame) {
	vm.frames = append(vm.frames, frame)
}

func (vm *VirtualMachine) popFrame() *Frame {
	f := vm.frames[len(vm.frames)-1]
	vm.frames = vm.frames[:len(vm.frames)-1]
	return f
}

func (vm *VirtualMachine) runFrame(frame *Frame) {
	vm.pushFrame(frame)
	vm.popFrame().execute()
}

// ExtractValue returns the top value from the top frame
func (vm *VirtualMachine) ExtractValue() object.Object {
	if len(vm.frames) < 1 || len(vm.frames[0].stack.objects) < 1 {
		return nil
	}
	return vm.frames[0].stack.top()
}
