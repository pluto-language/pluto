package evaluation

import (
	"github.com/Zac-Garby/pluto/ast"
)

func evalExpressions(exprs []ast.Expression, ctx *Context) []Object {
	var result []Object

	for _, expr := range exprs {
		o := eval(expr, ctx)

		if isErr(o) {
			return []Object{o}
		}

		result = append(result, o)
	}

	return result
}

func unwrapReturnValue(o Object) Object {
	if ret, ok := o.(*ReturnValue); ok {
		return ret.Value
	}

	return o
}

func isTruthy(o Object) bool {
	if o.Equals(O_NULL) || o.Equals(O_FALSE) {
		return false
	}

	if num, ok := o.(*Number); ok {
		return num.Value != 0
	}

	if col, ok := o.(Collection); ok {
		return len(col.Elements()) != 0
	}

	return true
}

func boolObj(t bool) Object {
	if t {
		return O_TRUE
	}

	return O_FALSE
}
