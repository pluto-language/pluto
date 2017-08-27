package evaluator

import (
	"reflect"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/object"
)

var (
	NEXT  = new(object.Next)
	BREAK = new(object.Break)

	NULL  = new(object.Null)
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}

	err   = object.Err
	isErr = object.IsErr
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

	var result object.Object

	switch node := n.(type) {
	/* Not literals */
	case *ast.AssignExpression:
		result = evalAssignExpression(*node, ctx)
	case *ast.BlockStatement:
		result = evalBlockStatement(*node, ctx)
	case *ast.BreakStatement:
		result = BREAK
	case *ast.ClassStatement:
		result = evalClassStatement(*node, ctx)
	case *ast.DeclareExpression:
		result = evalDeclareExpression(*node, ctx)
	case *ast.DotExpression:
		result = evalDotExpression(*node, ctx)
	case *ast.ExpressionStatement:
		result = eval(node.Expr, ctx)
	case *ast.ForLoop:
		result = evalForLoop(*node, ctx)
	case *ast.FunctionDefinition:
		result = evalFunctionDefinition(*node, ctx)
	case *ast.FunctionCall:
		result = evalFunctionCall(*node, ctx)
	case *ast.IfExpression:
		result = evalIfExpression(*node, ctx)
	case *ast.InfixExpression:
		result = evalInfixExpression(*node, ctx)
	case *ast.MatchExpression:
		result = evalMatchExpression(*node, ctx)
	case *ast.MethodCall:
		result = evalMethodCall(*node, ctx)
	case *ast.NextStatement:
		result = NEXT
	case *ast.ReturnStatement:
		result = evalReturnStatement(*node, ctx)
	case *ast.PrefixExpression:
		result = evalPrefixExpression(*node, ctx)
	case *ast.TryExpression:
		result = evalTryExpression(*node, ctx)
	case *ast.WhileLoop:
		result = evalWhileLoop(*node, ctx)

	/* Literals */
	case *ast.Array:
		result = evalArray(*node, ctx)
	case *ast.BlockLiteral:
		result = evalBlockLiteral(*node, ctx)
	case *ast.Boolean:
		result = &object.Boolean{Value: node.Value}
	case *ast.Char:
		result = &object.Char{Value: rune(node.Value)}
	case *ast.Identifier:
		result = evalIdentifier(*node, ctx)
	case *ast.Map:
		result = evalMap(*node, ctx)
	case *ast.Null:
		result = NULL
	case *ast.Number:
		result = &object.Number{Value: node.Value}
	case *ast.String:
		result = &object.String{Value: node.Value}
	case *ast.Tuple:
		result = evalTuple(*node, ctx)
	default:
		return err(ctx, "evaluation for %s not yet implemented", "NotImplementedError", reflect.TypeOf(n))
	}

	return result
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
