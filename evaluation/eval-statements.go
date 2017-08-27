package evaluation

import (
	"github.com/Zac-Garby/pluto/ast"
)

func evalBlockStatement(block ast.BlockStatement, ctx *Context) Object {
	if len(block.Statements) == 0 {
		return O_NULL
	}

	var result Object

	for _, stmt := range block.Statements {
		result = eval(stmt, ctx)

		if isErr(result) || result != nil &&
			(result.Type() == RETURN_VALUE ||
				result.Type() == NEXT ||
				result.Type() == BREAK) {
			return result
		}
	}

	return result
}

func evalClassStatement(node ast.ClassStatement, ctx *Context) Object {
	o := &Class{Name: node.Name.Token().Literal}

	if node.Parent != nil {
		o.Parent = eval(node.Parent, ctx)
		if isErr(o.Parent) {
			return o.Parent
		}
	} else if o.Name != "Base" {
		o.Parent = ctx.Get("Base")
		if o.Parent == nil {
			panic("The prelude isn't loaded, so Base isn't defined!")
		}
	}

	for _, n := range node.Methods {
		var method Object

		switch m := n.(type) {
		case *ast.FunctionDefinition:
			fn := Function{
				Pattern: m.Pattern,
				Body:    m.Body,
				Context: ctx,
			}

			method = &Method{Fn: fn}
		case *ast.InitDefinition:
			fn := Function{
				Pattern: m.Pattern,
				Body:    m.Body,
				Context: ctx,
			}

			method = &InitMethod{Fn: fn}

			initPattern := append(
				[]ast.Expression{node.Name},
				fn.Pattern...,
			)

			onInit := func(self *Function, ctx, enclosed *Context) Object {
				enclosed.Assign("self", &Instance{Base: o, Data: make(map[string]Object)})

				result := eval(self.Body, enclosed)
				if isErr(result) {
					return result
				}

				return enclosed.Get("self")
			}

			initFn := &Function{
				Pattern: initPattern,
				Body:    m.Body,
				Context: ctx,
				OnCall:  onInit,
			}

			ctx.AddFunction(initFn)
		}

		o.Methods = append(o.Methods, method)
	}

	ctx.Assign(o.Name, o)

	return o
}

func evalReturnStatement(node ast.ReturnStatement, ctx *Context) Object {
	if node.Value == nil {
		return &ReturnValue{Value: O_NULL}
	}

	val := eval(node.Value, ctx)
	if isErr(val) {
		return val
	}

	return &ReturnValue{Value: val}
}
