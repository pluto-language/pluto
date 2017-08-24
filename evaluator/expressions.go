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

func evalPrefixExpression(node ast.PrefixExpression, ctx *object.Context) object.Object {
	right := eval(node.Right, ctx)
	if isErr(right) {
		return right
	}

	return evalPrefix(node.Operator, right, ctx)
}

func evalPrefix(op string, right object.Object, ctx *object.Context) object.Object {
	if instance, ok := right.(*object.Instance); ok {
		return evalInstancePrefix(op, instance, ctx)
	}

	switch op {
	case "-":
		return evalMinusPrefix(right, ctx)
	case "+":
		return right
	case "!":
		return boolObj(!isTruthy(right))
	default:
		return err(ctx, "unknown operator: %s%s", "NotFoundError", op, right.Type)
	}
}

func evalInstancePrefix(op string, right *object.Instance, ctx *object.Context) object.Object {
	fnName, ok := prefixOverloads[op]
	if !ok {
		return err(ctx, "cannot overload operator %s", "NotFoundError", op)
	}

	if method := right.Base.(*object.Class).GetMethod(fnName); method != nil {
		args := map[string]object.Object{
			"self": right,
		}

		enclosed := ctx.EncloseWith(args)

		return eval(method.Fn.Body, enclosed)
	}

	return err(ctx, "unknown operator: %s%s. try overloading %s", "NotFoundError", op, right.Base.String(), fnName)
}

func evalMinusPrefix(right object.Object, ctx *object.Context) object.Object {
	if right.Type() != object.NUMBER {
		return err(ctx, "unknown operator: -%s", "NotFoundError", right.Type())
	}

	return &object.Number{Value: -right.(*object.Number).Value}
}
