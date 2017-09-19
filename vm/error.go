package vm

import (
	"fmt"
)

// ErrType is a runtime error type
type ErrType string

const (
	// ErrUnknown is thrown when there's an error, but the vm isn't
	// sure of what nature it is
	ErrUnknown = "Unknown"

	// ErrNoInstruction is thrown if an instruction in the bytecode
	// isn't yet implemented
	ErrNoInstruction = "NoInstruction"
)

// Error is a runtime error thrown in the virtual machine
type Error struct {
	Type    ErrType
	Message string
}

// Err creates a new runtime error with the given message and type
func Err(msg string, t ErrType) *Error {
	return &Error{
		Type:    t,
		Message: msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
