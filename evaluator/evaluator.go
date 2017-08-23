package evaluator

import (
	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/object"
)

var (
	NEXT  = new(object.Next)
	BREAK = new(object.Break)

	NULL  = new(object.Null)
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var (
	infixOverloads = map[string]string{
		"+":  "__plus $",
		"-":  "__minus $",
		"*":  "__times $",
		"/":  "__divide $",
		"**": "__exp $",
		"//": "__f_div $",
		`%`:  "__mod $",
		"==": "__eq $",
		"||": "__or $",
		"&&": "__and $",
		"|":  "__b_or $",
		"&":  "__b_and $",
		".":  "__get $",
	}

	prefixOverloads = map[string]string{
		"+": "__no_op",
		"-": "__negate",
		"!": "__invert",
	}
)

func EvaluateProgram(prog ast.Program, ctx *object.Context) object.Object {
	return evalProgram(&prog, ctx)
}

func eval(n ast.Node, ctx *object.Context) object.Object {
	/** Evaluation function naming **
	 * Every AST node evaluation function's name should be in the form:
	 *
	 *    evalNODE(ast.Node, *object.Context) object.Object
	 *
	 * ...where NODE is the actual name of the AST node struct type.
	 * For example: evalMatchExpression(node ast.Node, ctx *object.Context) object.Object
	 *
	 * Also, try to keep the switch branches below in alphabetical order.
	 */

	switch node := n.(type) {
	/* Not literals */
	case ast.AssignExpression:
		// return evalAssignExpression(node, ctx)
	case ast.BlockStatement:
		return evalBlockStatement(node, ctx)
	case ast.ClassStatement:
		// return evalClassStatement(node, ctx)
	case ast.DeclareExpression:
		// return evalDeclareExpression(node, ctx)
	case ast.DotExpression:
		// return evalDotExpression(node, ctx)
	case ast.ExpressionStatement:
		return eval(node.Expr, ctx)
	case ast.ForLoop:
		// return evalForLoop(node, ctx)
	case ast.IfExpression:
		// return evalIfExpression(node, ctx)
	case ast.InfixExpression:
		// return evalInfixExpression(node, ctx)
	case ast.MatchExpression:
		// return evalMatchExpression(node, ctx)
	case ast.MethodCall:
		// return evalMethodCall(node, ctx)
	case ast.ReturnStatement:
		// return evalReturnStatement(node, ctx)
	case ast.PrefixExpression:
		// return evalPrefixExpression(node, ctx)
	case ast.TryExpression:
		// return evalTryExpression(node, ctx)
	case ast.WhileLoop:
		// return evalWhileLoop(node, ctx)

	/* Literals */
	case ast.Array:
		return evalArray(node, ctx)
	case ast.BlockLiteral:
		return evalBlockLiteral(node, ctx)
	case ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case ast.Char:
		return &object.Char{Value: rune(node.Value)}
	case ast.Identifier:
		return evalIdentifier(node, ctx)
	case ast.Map:
		return evalMap(node, ctx)
	case ast.Null:
		return NULL
	case ast.Number:
		return &object.Number{Value: node.Value}
	case ast.String:
		return &object.String{Value: node.Value}
	case ast.Tuple:
		return evalTuple(node, ctx)
	}

	return err(ctx, "evaluation not yet implemented", "NotImplementedError")
}

func evalProgram(prog *ast.Program, ctx *object.Context) object.Object {
	if len(prog.Statements) == 0 {
		return NULL
	}

	var result object.Object

	for _, stmt := range prog.Statements {
		result = eval(stmt, ctx)

		if isErr(result) {
			return result
		}

		if ret, ok := result.(*object.ReturnValue); ok {
			return ret.Value
		}

		if _, ok := result.(*object.Next); ok {
			return NULL
		}

		switch obj := result.(type) {
		case *object.ReturnValue:
			return obj.Value
		case *object.Next, *object.Break:
			return NULL
		}
	}

	return result
}
