package evaluation

import "github.com/Zac-Garby/pluto/ast"

type Context struct {
	Store     map[string]Object
	Functions []*Function

	Outer *Context
}

func (c *Context) Enclose() *Context {
	return &Context{
		Outer: c,
	}
}

func (c *Context) EncloseWith(args map[string]Object) *Context {
	return &Context{
		Store: args,
		Outer: c,
	}
}

func (c *Context) Get(key string) Object {
	if obj, ok := c.Store[key]; ok {
		return obj
	} else if c.Outer != nil {
		return c.Outer.Get(key)
	} else {
		return nil
	}
}

func (c *Context) Assign(key string, obj Object) {
	if c.Outer != nil {
		if v := c.Outer.Get(key); v != nil {
			c.Outer.Assign(key, obj)
			return
		}
	}

	c.Store[key] = obj
}

func (c *Context) Declare(key string, obj Object) {
	c.Store[key] = obj
}

func (c *Context) AddFunction(fn Object) {
	if _, ok := fn.(*Function); !ok {
		panic("Not a function!")
	}

	c.Functions = append(c.Functions, fn.(*Function))
}

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
