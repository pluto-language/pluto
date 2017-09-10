package object

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
)

/* Structs */
type (
	// Function is a normal Pluto function, referenced by its pattern
	Function struct {
		Pattern   []ast.Expression
		Body      bytecode.Code
		Constants []Object
		Names     []string
		Patterns  []string
		OnCall    func(self *Function) Object
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
func (f *Function) Type() Type { return FunctionType }

// Type returns the type of the object
func (i *InitMethod) Type() Type { return InitType }

// Type returns the type of the object
func (m *Method) Type() Type { return MethodType }

/* Equals() methods */

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
func (f *Function) String() string {
	return "<function>"
}

func (i *InitMethod) String() string {
	return "<init method>"
}

func (m *Method) String() string {
	return "<method>"
}
