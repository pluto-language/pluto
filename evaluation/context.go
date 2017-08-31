package evaluation

import (
	"github.com/Zac-Garby/pluto/ast"
)

// Context is an evaluation scope, containing
// defined variables and functions, as well as
// the outer scope, and imported packages
type Context struct {
	Store     map[string]Object
	Functions []*Function
	Packages  map[string]*Package

	Outer *Context
}

// NewContext returns a new, empty context
func NewContext() *Context {
	return &Context{
		Store:    make(map[string]Object),
		Packages: make(map[string]*Package),
	}
}

// Enclose creates a new, empty context,
// with Outer predefined
func (c *Context) Enclose() *Context {
	return &Context{
		Store:    make(map[string]Object),
		Outer:    c,
		Packages: make(map[string]*Package),
	}
}

// EncloseWith is the same as Enclose, but predefines
// the variables in 'args'
func (c *Context) EncloseWith(args map[string]Object) *Context {
	return &Context{
		Store:    args,
		Outer:    c,
		Packages: make(map[string]*Package),
	}
}

// Get searches for a variable
func (c *Context) Get(key string) Object {
	if obj, ok := c.Store[key]; ok {
		return obj
	} else if c.Outer != nil {
		return c.Outer.Get(key)
	} else {
		return nil
	}
}

// Assign assigns a value to a variable. Can bubble up
// to parent scopes
func (c *Context) Assign(key string, obj Object) {
	if c.Outer != nil {
		if v := c.Outer.Get(key); v != nil {
			c.Outer.Assign(key, obj)
			return
		}
	}

	c.Store[key] = obj
}

// Declare declares a variable as a value. Cannot
// bubble up to parent scopes
func (c *Context) Declare(key string, obj Object) {
	c.Store[key] = obj
}

// AddFunction defines a function
func (c *Context) AddFunction(fn Object) {
	if _, ok := fn.(*Function); !ok {
		panic("Not a function!")
	}

	c.Functions = append(c.Functions, fn.(*Function))
}

// GetFunction gets a function matching the given
// pattern
func (c *Context) GetFunction(pattern []ast.Expression) interface{} {
	for _, fn := range c.Functions {
		if len(pattern) != len(fn.Pattern) {
			continue
		}

		matched := true

		for i, item := range pattern {
			fItem := fn.Pattern[i]

			if itemID, ok := item.(*ast.Identifier); ok {
				if fItemID, ok := fItem.(*ast.Identifier); ok {
					if itemID.Value != fItemID.Value {
						matched = false
					}
				}
			} else if _, ok := item.(*ast.Argument); !ok {
				matched = false
			} else if _, ok := fItem.(*ast.Parameter); !ok {
				matched = false
			}
		}

		if matched {
			return fn
		}
	}

	if c.Outer != nil {
		return c.Outer.GetFunction(pattern)
	}

	for _, fn := range GetBuiltins() {
		if len(pattern) != len(fn.Pattern) {
			continue
		}

		matched := true

		for i, item := range pattern {
			fItem := fn.Pattern[i]

			if itemID, ok := item.(*ast.Identifier); ok {
				if fItem[0] != '$' {
					if itemID.Value != fItem {
						matched = false
					}
				}
			} else if _, ok := item.(*ast.Argument); !ok {
				matched = false
			} else if !(fItem[0] == '$') {
				matched = false
			}
		}

		if matched {
			return fn
		}
	}

	return nil
}
