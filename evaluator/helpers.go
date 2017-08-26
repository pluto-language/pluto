package evaluator

import (
	"strings"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/object"
)

func evalExpressions(exprs []ast.Expression, ctx *object.Context) []object.Object {
	var result []object.Object

	for _, expr := range exprs {
		o := eval(expr, ctx)

		if isErr(o) {
			return []object.Object{o}
		}

		result = append(result, o)
	}

	return result
}

func unwrapReturnValue(o object.Object) object.Object {
	if ret, ok := o.(*object.ReturnValue); ok {
		return ret.Value
	}

	return o
}

func makeCollection(t object.Type, elems []object.Object, ctx *object.Context) object.Object {
	switch t {
	case object.ARRAY:
		return &object.Array{Value: elems}
	case object.TUPLE:
		return &object.Tuple{Value: elems}
	case object.STRING:
		var strs []string

		for _, elem := range elems {
			strs = append(strs, elem.String())
		}

		return &object.String{Value: strings.Join(strs, "")}
	default:
		return err(ctx, "could not form a collection of type %s", "TypeError", t)
	}
}

func executeAppliedBlock(ab *object.AppliedBlock, ctx *object.Context) object.Object {
	if len(ab.Args) != len(ab.Block.Params) {
		return err(
			ctx,
			"wrong amount of arguments to a block. expected %d, got %d",
			"TypeError",
			len(ab.Block.Params),
			len(ab.Args),
		)
	}

	argDict := make(map[string]object.Object)

	for i, param := range ab.Block.Params {
		pval := param.(*ast.Identifier).Value
		argDict[pval] = ab.Args[i]
	}

	enclosed := ctx.EncloseWith(argDict)
	return eval(ab.Block.Body, enclosed)
}

func isTruthy(o object.Object) bool {
	if o.Equals(NULL) || o.Equals(FALSE) {
		return false
	}

	if num, ok := o.(*object.Number); ok {
		return num.Value != 0
	}

	if col, ok := o.(object.Collection); ok {
		return len(col.Elements()) != 0
	}

	return true
}

func boolObj(t bool) object.Object {
	if t {
		return TRUE
	}

	return FALSE
}
