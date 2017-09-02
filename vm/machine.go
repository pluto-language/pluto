package vm

import (
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

type store map[string]object.Object

// VirtualMachine is the base struct
// for the Pluto VM. It stores information
// such as the call stack and the last
// error thrown.
type VirtualMachine struct {
	frames      []*Frame
	frame       *Frame
	returnValue object.Object
	lastError   error
}

// New returns a new virtual machine
func New() *VirtualMachine {
	return &VirtualMachine{
		frames:      make([]*Frame, 0),
		returnValue: nil,
		lastError:   nil,
	}
}

// Run executes the supplied bytecode
func (vm *VirtualMachine) Run(code bytecode.Code, globals, locals store, constants []object.Object) {
	frame := vm.makeFrame(code, make(store), globals, locals, constants)

	vm.pushFrame(frame)
	vm.runFrame(frame)
}

// RunDefault executes the bytecode with
// empty globals and locals
func (vm *VirtualMachine) RunDefault(code bytecode.Code, constants []object.Object) {
	vm.Run(code, make(store), make(store), constants)
}

func (vm *VirtualMachine) makeFrame(code bytecode.Code, args, globals, locals store, constants []object.Object) *Frame {
	frame := &Frame{
		code:      code,
		locals:    locals,
		globals:   globals,
		offset:    0,
		stack:     newStack(),
		constants: constants,
	}

	for k, v := range args {
		locals[k] = v
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
	if len(vm.frames) == 0 {
		return nil
	}

	return vm.frames[0].stack.top()
}
