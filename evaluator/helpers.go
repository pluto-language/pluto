package evaluator

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/object"
)

func evalExpressions(exprs []ast.Expression, ctx *object.Context) []object.Object {
	var result []object.Object

	for _, expr := range exprs {
		o := Evaluate(expr, ctx)

		if isErr(o) {
			return []object.Object{o}
		}

		result = append(result, o)
	}

	return result
}

func evalBlockStatement(block ast.BlockStatement, ctx *object.Context) object.Object {
	if len(block.Statements) == 0 {
		return NULL
	}

	var result object.Object

	for _, stmt := range block.Statements {
		result = Evaluate(stmt, ctx)

		if isErr(result) || result != nil &&
			(result.Type() == object.RETURN_VALUE ||
				result.Type() == object.NEXT ||
				result.Type() == object.BREAK) {
			return result
		}
	}

	return result
}

func unwrapReturnValue(o object.Object) object.Object {
	if ret, ok := o.(*object.ReturnValue); ok {
		return ret.Value
	}

	return o
}

func isTruthy(o object.Object) bool {
	if o == NULL || o == FALSE {
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
