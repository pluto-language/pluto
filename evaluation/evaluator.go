package evaluation

import (
	"reflect"

	"github.com/Zac-Garby/pluto/ast"
)

var (
	// NextObj is a predefined instance of Next
	NextObj = new(Next)

	// BreakObj is a predefined instance of Break
	BreakObj = new(Break)

	// NullObj is a predefined instance of Null
	NullObj = new(Null)

	// TrueObj is a predefined instance of Boolean, whose value is true
	TrueObj = &Boolean{Value: true}

	// FalseObj is a predefined instance of Boolean, whose value is false
	FalseObj = &Boolean{Value: false}
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

// EvaluateProgram evaluates a program in the given context
func EvaluateProgram(prog ast.Program, ctx *Context) Object {
	return evalProgram(&prog, ctx)
}

func eval(n ast.Node, ctx *Context) Object {
	/** Evaluation function naming **
	 * Every AST node evaluation function's name should be in the form:
	 *
	 *    evalNODE(ast.Node, *Context) Object
	 *
	 * ...where NODE is the actual name of the AST node struct type.
	 * For example: evalMatchExpression(node ast.Node, ctx *Context) Object
	 *
	 * Also, try to keep the switch branches below in alphabetical order.
	 */

	var result Object

	switch node := n.(type) {
	/* Not literals */
	case *ast.AssignExpression:
		result = evalAssignExpression(*node, ctx)
	case *ast.BlockStatement:
		result = evalBlockStatement(*node, ctx)
	case *ast.BreakStatement:
		result = BreakObj
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
	case *ast.ImportStatement:
		result = evalImportStatement(*node, ctx)
	case *ast.IndexExpression:
		result = evalIndexExpression(*node, ctx)
	case *ast.InfixExpression:
		result = evalInfixExpression(*node, ctx)
	case *ast.MatchExpression:
		result = evalMatchExpression(*node, ctx)
	case *ast.MethodCall:
		result = evalMethodCall(*node, ctx)
	case *ast.NextStatement:
		result = NextObj
	case *ast.QualifiedFunctionCall:
		result = evalQualifiedFunctionCall(*node, ctx)
	case *ast.ReturnStatement:
		result = evalReturnStatement(*node, ctx)
	case *ast.PrefixExpression:
		result = evalPrefixExpression(*node, ctx)
	case *ast.TryExpression:
		result = evalTryExpression(*node, ctx)
	case *ast.UseStatement:
		result = evalUseStatement(*node, ctx)
	case *ast.WhileLoop:
		result = evalWhileLoop(*node, ctx)

	/* Literals */
	case *ast.Array:
		result = evalArray(*node, ctx)
	case *ast.BlockLiteral:
		result = evalBlockLiteral(*node, ctx)
	case *ast.Boolean:
		result = &Boolean{Value: node.Value}
	case *ast.Char:
		result = &Char{Value: rune(node.Value)}
	case *ast.Identifier:
		result = evalIdentifier(*node, ctx)
	case *ast.Map:
		result = evalMap(*node, ctx)
	case *ast.Null:
		result = NullObj
	case *ast.Number:
		result = &Number{Value: node.Value}
	case *ast.String:
		result = &String{Value: node.Value}
	case *ast.Tuple:
		result = evalTuple(*node, ctx)
	default:
		return Err(ctx, "evaluation for %s not yet implemented", "NotImplementedError", reflect.TypeOf(n))
	}

	return result
}

func evalProgram(prog *ast.Program, ctx *Context) Object {
	if len(prog.Statements) == 0 {
		return NullObj
	}

	var result Object

	for _, stmt := range prog.Statements {
		result = eval(stmt, ctx)

		if IsErr(result) {
			return result
		}

		switch obj := result.(type) {
		case *ReturnValue:
			return obj.Value
		case *Next, *Break:
			return NullObj
		}
	}

	return result
}
