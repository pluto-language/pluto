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
	returnValue object.Object
	Error       error
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
func (vm *VirtualMachine) Run(code bytecode.Code, locals *Store, constants []object.Object) {
	frame := vm.makeFrame(code, NewStore(), locals, constants)

	vm.pushFrame(frame)
	vm.runFrame(frame)
}

// RunDefault executes the bytecode with
// empty globals and locals
func (vm *VirtualMachine) RunDefault(code bytecode.Code, constants []object.Object) {
	vm.Run(code, NewStore(), constants)
}

func (vm *VirtualMachine) makeFrame(code bytecode.Code, args, locals *Store, constants []object.Object) *Frame {
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
		locals.Data[v] = args.Data[v]
	}

	return frame
}

func (vm *VirtualMachine) pushFrame(frame *Frame) {
	vm.frames = append(vm.frames, frame)
}

func (vm *VirtualMachine) popFrame(frame *Frame) {
	vm.frames = vm.frames[:len(vm.frames)-1]
}

func (vm *VirtualMachine) runFrame(frame *Frame) {
	frame.execute()
}

// ExtractValue returns the top value from the top frame
func (vm *VirtualMachine) ExtractValue() object.Object {
	if len(vm.frames) < 1 || len(vm.frames[0].stack.objects) < 1 {
		return nil
	}
	return vm.frames[0].stack.top()
}
