package object

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
)

/* Structs */
type (
	// ReturnValue is a value which has been returned via a return statement
	ReturnValue struct {
		Value Object
	}

	// Next is a single-value type which is created from a next statement
	Next struct{}

	// Break is a single-value type which is created from a break statement
	Break struct{}

	// Function is a normal Pluto function, referenced by its pattern
	Function struct {
		Pattern []ast.Expression
		Body    bytecode.Code
		OnCall  func(self *Function) Object
	}

	// InitMethod is an initializer method on a class
	InitMethod struct {
		Fn Function
	}

	// Method is a regular method on a class
	Method struct {
		Fn Function
	}
)

/* Type() methods */

// Type returns the type of the object
func (r *ReturnValue) Type() Type { return ReturnValueType }

// Type returns the type of the object
func (n *Next) Type() Type { return NextType }

// Type returns the type of the object
func (b *Break) Type() Type { return BreakType }

// Type returns the type of the object
func (f *Function) Type() Type { return FunctionType }

// Type returns the type of the object
func (i *InitMethod) Type() Type { return InitType }

// Type returns the type of the object
func (m *Method) Type() Type { return MethodType }

/* Equals() methods */

// Equals checks if two objects are equal to each other
func (r *ReturnValue) Equals(o Object) bool {
	if other, ok := o.(*ReturnValue); ok {
		return r.Value.Equals(other.Value)
	}

	return false
}

// Equals checks if two objects are equal to each other
func (n *Next) Equals(o Object) bool {
	_, ok := o.(*Next)
	return ok
}

// Equals checks if two objects are equal to each other
func (b *Break) Equals(o Object) bool {
	_, ok := o.(*Break)
	return ok
}

// Equals checks if two objects are equal to each other
func (f *Function) Equals(o Object) bool {
	_, ok := o.(*Function)
	return ok
}

// Equals checks if two objects are equal to each other
func (i *InitMethod) Equals(o Object) bool {
	_, ok := o.(*Function)
	return ok
}

// Equals checks if two objects are equal to each other
func (m *Method) Equals(o Object) bool {
	_, ok := o.(*Method)
	return ok
}

/* Stringer implementations */
func (r *ReturnValue) String() string {
	return r.Value.String()
}

func (n *Next) String() string {
	return "<next>"
}

func (b *Break) String() string {
	return "<break>"
}

func (f *Function) String() string {
	return "<function>"
}

func (i *InitMethod) String() string {
	return "<init method>"
}

func (m *Method) String() string {
	return "<method>"
}
