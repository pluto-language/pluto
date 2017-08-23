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
