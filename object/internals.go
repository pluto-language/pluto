package object

import "github.com/Zac-Garby/pluto/ast"

/* Structs */
type (
	ReturnValue struct {
		Value Object
	}

	Next  struct{}
	Break struct{}

	Function struct {
		Pattern []ast.Expression
		Body    ast.Statement
		// Context context.Context
	}

	InitMethod struct {
		Fn Function
	}

	Method struct {
		Fn Function
	}
)

/* Type() methods */
func (_ *ReturnValue) Type() Type { return RETURN_VALUE }
func (_ *Next) Type() Type        { return NEXT }
func (_ *Break) Type() Type       { return BREAK }
func (_ *Function) Type() Type    { return FUNCTION }
func (_ *InitMethod) Type() Type  { return INIT }
func (_ *Method) Type() Type      { return METHOD }

/* Equals() methods */
func (r *ReturnValue) Equals(o Object) bool {
	if other, ok := o.(*ReturnValue); ok {
		return r.Value.Equals(other.Value)
	}

	return false
}

func (_ *Next) Equals(o Object) bool {
	_, ok := o.(*Next)
	return ok
}

func (_ *Break) Equals(o Object) bool {
	_, ok := o.(*Break)
	return ok
}

func (_ *Function) Equals(o Object) bool {
	_, ok := o.(*Function)
	return ok
}

func (_ *InitMethod) Equals(o Object) bool {
	_, ok := o.(*Function)
	return ok
}

func (_ *Method) Equals(o Object) bool {
	_, ok := o.(*Method)
	return ok
}

/* Stringer implementations */
func (r *ReturnValue) String() string {
	return r.Value.String()
}

func (_ *Next) String() string {
	return "<next>"
}

func (_ *Break) String() string {
	return "<break>"
}

func (_ *Function) String() string {
	return "<function>"
}

func (_ *InitMethod) String() string {
	return "<init method>"
}

func (_ *Method) String() string {
	return "<method>"
}
