package evaluator

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/object"
)

func evalIdentifier(node ast.Identifier, ctx *object.Context) object.Object {
	val := ctx.Get(node.Value)

	if val != nil {
		return val
	}

	return err(ctx, "`%s` is not found", "NotFoundError", node.Value)
}

func evalTuple(node ast.Tuple, ctx *object.Context) object.Object {
	elements := evalExpressions(node.Value, ctx)

	if len(elements) == 1 && isErr(elements[0]) {
		return elements[0]
	}

	return &object.Tuple{Value: elements}
}

func evalArray(node ast.Array, ctx *object.Context) object.Object {
	elements := evalExpressions(node.Elements, ctx)

	if len(elements) == 1 && isErr(elements[0]) {
		return elements[0]
	}

	return &object.Array{Value: elements}
}

func evalMap(node ast.Map, ctx *object.Context) object.Object {
	var dict map[object.Object]object.Object

	for k, v := range node.Pairs {
		key := eval(k, ctx)
		if isErr(key) {
			return key
		}

		value := eval(v, ctx)
		if isErr(value) {
			return value
		}

		dict[key] = value
	}

	return &object.Map{Pairs: dict}
}

func evalBlockLiteral(node ast.BlockLiteral, ctx *object.Context) object.Object {
	return &object.Block{Params: node.Params, Body: node.Body}
}
