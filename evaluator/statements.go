package evaluator

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/object"
)

func evalBlockStatement(block ast.BlockStatement, ctx *object.Context) object.Object {
	if len(block.Statements) == 0 {
		return NULL
	}

	var result object.Object

	for _, stmt := range block.Statements {
		result = eval(stmt, ctx)

		if isErr(result) || result != nil &&
			(result.Type() == object.RETURN_VALUE ||
				result.Type() == object.NEXT ||
				result.Type() == object.BREAK) {
			return result
		}
	}

	return result
}

func evalClassStatement(node ast.ClassStatement, ctx *object.Context) object.Object {
	o := &object.Class{Name: node.Name.Token().Literal}

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
		var method object.Object

		switch m := n.(type) {
		case *ast.FunctionDefinition:
			fn := object.Function{
				Pattern: m.Pattern,
				Body:    m.Body,
				Context: ctx,
			}

			method = &object.Method{Fn: fn}
		case *ast.InitDefinition:
			fn := object.Function{
				Pattern: m.Pattern,
				Body:    m.Body,
				Context: ctx,
			}

			method = &object.InitMethod{Fn: fn}

			initPattern := append(
				[]ast.Expression{node.Name},
				fn.Pattern...,
			)

			onInit := func(self *object.Function, ctx, enclosed *object.Context) object.Object {
				enclosed.Assign("self", &object.Instance{Base: o})

				result := eval(self.Body, enclosed)
				if isErr(result) {
					return result
				}

				return enclosed.Get("self")
			}

			initFn := &object.Function{
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

func evalReturnStatement(node ast.ReturnStatement, ctx *object.Context) object.Object {
	if node.Value == nil {
		return &object.ReturnValue{Value: NULL}
	}

	val := eval(node.Value, ctx)
	if isErr(val) {
		return val
	}

	return &object.ReturnValue{Value: val}
}
