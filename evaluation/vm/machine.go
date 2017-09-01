package vm

import (
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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
