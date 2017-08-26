package evaluator

import (
	"math"
	"strings"

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
	m := &object.Map{
		Values: make(map[string]object.Object),
		Keys:   make(map[string]object.Object),
	}

	for k, v := range node.Pairs {
		key := eval(k, ctx)
		if isErr(key) {
			return key
		}

		value := eval(v, ctx)
		if isErr(value) {
			return value
		}

		m.Set(key, value)
	}

	return m
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

func evalInfixExpression(node ast.InfixExpression, ctx *object.Context) object.Object {
	left := eval(node.Left, ctx)
	if isErr(left) {
		return left
	}

	right := eval(node.Right, ctx)
	if isErr(right) {
		return right
	}

	op := node.Operator

	if lCol, ok := left.(object.Collection); ok {
		if rCol, ok := right.(object.Collection); ok {
			return evalCollectionInfix(op, lCol, rCol, ctx)
		}
	}

	if instance, ok := left.(*object.Instance); ok {
		return evalInstanceInfix(op, instance, right, ctx)
	}

	switch op {
	case "&&":
		return boolObj(isTruthy(left) && isTruthy(right))
	case "||":
		return boolObj(isTruthy(left) || isTruthy(right))
	case "==":
		return boolObj(left.Equals(right))
	case "!=":
		return boolObj(!left.Equals(right))
	case "?":
		if left == NULL {
			return right
		}

		return left
	}

	if left.Type() == object.NUMBER && right.Type() == object.NUMBER {
		return evalNumberInfix(op, left.(*object.Number), right.(*object.Number), ctx)
	}

	if (left.Type() == object.CHAR || left.Type() == object.STRING) &&
		(right.Type() == object.CHAR || right.Type() == object.STRING) {
		return evalCharStringInfix(op, left, right, ctx)
	}

	if left.Type() == object.CHAR && right.Type() == object.NUMBER {
		ch := string(left.(*object.Char).Value)
		amount := int(math.Floor(right.(*object.Number).Value))

		return &object.String{Value: strings.Repeat(ch, amount)}
	}

	if lCol, ok := left.(object.Collection); ok {
		if right.Type() == object.NUMBER && op == "*" {
			var result []object.Object
			elems := lCol.Elements()
			amount := int(math.Floor(right.(*object.Number).Value))

			for i := 0; i < amount; i++ {
				result = append(result, elems...)
			}

			return makeCollection(left.Type(), result, ctx)
		}
	}

	return err(ctx, "unknown operator: %s %s %s", "NotFoundError", left.Type(), op, right.Type())
}

func evalCollectionInfix(op string, left, right object.Collection, ctx *object.Context) object.Object {
	l := left.Elements()
	r := right.Elements()

	switch op {
	case "==":
		return boolObj(left.Equals(right))
	case "!=":
		return boolObj(!left.Equals(right))
	case "+":
		return makeCollection(left.Type(), append(l, r...), ctx)
	case "-":
		var elems []object.Object

		for _, el := range l {
			blacklisted := false

			for _, rel := range r {
				if el.Equals(rel) {
					blacklisted = true
				}
			}

			if !blacklisted {
				elems = append(elems, el)
			}
		}

		return makeCollection(left.Type(), elems, ctx)
	case "&", "&&":
		var elems []object.Object

		for _, el := range l {
			both := false

			for _, rel := range r {
				if el.Equals(rel) {
					both = true
					break
				}
			}

			if both {
				elems = append(elems, el)
			}
		}

		return makeCollection(left.Type(), elems, ctx)
	case "|", "||":
		var elems []object.Object

		for _, el := range append(l, r...) {
			unique := true

			for _, rel := range elems {
				if el.Equals(rel) {
					unique = false
					break
				}
			}

			if unique {
				elems = append(elems, el)
			}
		}

		return makeCollection(left.Type(), elems, ctx)
	default:
		return err(ctx, "unknown operator: %s %s %s", "NotFoundError", left.Type(), op, right.Type())
	}
}

func evalNumberInfix(op string, left, right *object.Number, ctx *object.Context) object.Object {
	l := left.Value
	r := right.Value

	switch op {
	case "+":
		return &object.Number{Value: l + r}
	case "-":
		return &object.Number{Value: l - r}
	case "*":
		return &object.Number{Value: l * r}
	case "/":
		return &object.Number{Value: l / r}
	case "&":
		return &object.Number{Value: float64(int(l) & int(r))}
	case "|":
		return &object.Number{Value: float64(int(l) | int(r))}
	case "**":
		return &object.Number{Value: math.Pow(l, r)}
	case "//":
		return &object.Number{Value: math.Floor(l / r)}
	case "%":
		return &object.Number{Value: math.Mod(l, r)}
	case "<":
		return boolObj(l < r)
	case ">":
		return boolObj(l > r)
	case "<=":
		return boolObj(l <= r)
	case ">=":
		return boolObj(l >= r)
	default:
		return err(ctx, "unknown operator: %s %s %s", "NotFoundError", left.Type(), op, right.Type())
	}
}

func evalCharStringInfix(op string, left, right object.Object, ctx *object.Context) object.Object {
	var l, r string

	if lch, ok := left.(*object.Char); ok {
		l = string(lch.Value)
	} else if lstr, ok := left.(*object.String); ok {
		l = lstr.Value
	}

	if rch, ok := right.(*object.Char); ok {
		r = string(rch.Value)
	} else if rstr, ok := right.(*object.String); ok {
		r = rstr.Value
	}

	switch op {
	case "+":
		return &object.String{Value: l + r}
	case "-":
		var val string

		for _, ch := range strings.Split(l, "") {
			if ch != r {
				val += ch
			}
		}

		return &object.String{Value: val}
	default:
		return err(ctx, "unknown operator: %s %s %s", "NotFoundError", left.Type(), op, right.Type())
	}
}

func evalInstanceInfix(op string, left *object.Instance, right object.Object, ctx *object.Context) object.Object {
	fnName, ok := infixOverloads[op]
	if !ok {
		return err(ctx, "cannot overload operator %s", "NotFoundError", op)
	}

	if method := left.Base.(*object.Class).GetMethod(fnName); method != nil {
		methodPattern := method.Fn.Pattern

		args := map[string]object.Object{
			"self": left,
		}

		for _, item := range methodPattern {
			if param, ok := item.(*ast.Parameter); ok {
				args[param.Name] = right
				break
			}
		}

		enclosed := ctx.EncloseWith(args)

		return eval(method.Fn.Body, enclosed)
	}

	return err(
		ctx, "unknown operator: %s %s %s. try overloading %s",
		"NotFoundError",
		left.Base.String(),
		op,
		right.Type(),
		fnName,
	)
}

func evalAssignExpression(node ast.AssignExpression, ctx *object.Context) object.Object {
	right := eval(node.Value, ctx)
	if isErr(right) {
		return right
	}

	if dot, ok := node.Name.(*ast.DotExpression); ok {
		o := eval(dot.Left, ctx)
		if isErr(o) {
			return o
		}

		ctr, ok := o.(object.Container)
		if !ok {
			return err(ctx, "can only access fields of containers", "TypeError")
		}

		if id, ok := dot.Right.(*ast.Identifier); ok {
			ctr.Set(&object.String{Value: id.Value}, right)
			return right
		}

		return err(ctx, "an identifier is expected to follow a dot operator", "SyntaxError")
	}

	if left, ok := node.Name.(*ast.Identifier); ok {
		ctx.Assign(left.Value, right)
		return right
	}

	return err(ctx, "can only assign to identifiers!", "SyntaxError")
}

func evalDeclareExpression(node ast.DeclareExpression, ctx *object.Context) object.Object {
	right := eval(node.Value, ctx)
	if isErr(right) {
		return right
	}

	if left, ok := node.Name.(*ast.Identifier); ok {
		ctx.Declare(left.Value, right)
		return right
	}

	return err(ctx, "cannot declare a non-identifier!", "SyntaxError")
}

func evalDotExpression(node ast.DotExpression, ctx *object.Context) object.Object {
	left := eval(node.Left, ctx)
	if isErr(left) {
		return left
	}

	if field, ok := node.Right.(*ast.Identifier); ok {
		if cnt, ok := left.(object.Container); ok {
			return cnt.Get(&object.String{Value: field.Value})
		}

		return err(ctx, "cannot access fields of %s", "TypeError", left.Type())
	}

	return err(ctx, "an identifier is expected after a dot '.'", "SyntaxError")
}

func evalWhileLoop(node ast.WhileLoop, ctx *object.Context) object.Object {
	var steps []object.Object

	for {
		condition := eval(node.Condition, ctx)
		if isErr(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result := eval(node.Body, ctx.Enclose())
		if isErr(result) {
			return result
		}

		if result.Type() == object.BREAK {
			break
		}

		steps = append(steps, result)
	}

	return &object.Array{Value: steps}
}

func evalForLoop(node ast.ForLoop, ctx *object.Context) object.Object {
	var (
		v    = node.Var
		body = node.Body
		col  = eval(node.Collection, ctx)
	)

	if isErr(col) {
		return col
	}

	var items []object.Object
	if collection, ok := col.(object.Collection); ok {
		items = collection.Elements()
	} else {
		return err(ctx, "cannot perform a for-loop over the non-collection type: %s", "TypeError", col.Type())
	}

	var steps []object.Object

	for _, item := range items {
		enclosed := ctx.EncloseWith(map[string]object.Object{
			v.Token().Literal: item,
		})

		result := eval(body, enclosed)
		if isErr(result) {
			return result
		}

		if result.Type() == object.BREAK {
			break
		}

		steps = append(steps, result)
	}

	return &object.Array{Value: steps}
}

func evalFunctionCall(node ast.FunctionCall, ctx *object.Context) object.Object {
	function := ctx.GetFunction(node.Pattern)

	if function == nil {
		var patternStrings []string

		for _, item := range node.Pattern {
			if id, ok := item.(*ast.Identifier); ok {
				patternStrings = append(patternStrings, id.Value)
			} else {
				patternStrings = append(patternStrings, "$")
			}
		}

		patternString := strings.Join(patternStrings, " ")

		return err(ctx, "no function matching the pattern: %s", "NotFoundError", patternString)
	}

	args := make(map[string]object.Object)

	var result object.Object

	if function.Type() == object.FUNCTION {
		for i, item := range node.Pattern {
			fItem := function.Pattern[i]

			if arg, ok := item.(*ast.Argument); ok {
				if param, ok := fItem.(*ast.Parameter); ok {
					evaled := eval(arg.Value, ctx)
					if isErr(evaled) {
						return evaled
					}

					args[param.Name] = evaled
				}
			}
		}

		enclosed := ctx.EncloseWith(args)
		onCallResult := function.OnCall(*function, ctx, enclosed)

		if onCallResult != nil || isErr(onCallResult) {
			return onCallResult
		}

		result = eval(function.Body, enclosed)
		if isErr(result) {
			return result
		}
	}

	if result == nil {
		return NULL
	}

	return unwrapReturnValue(result)
}

func evalFunctionDefinition(node ast.FunctionDefinition, ctx *object.Context) object.Object {
	function := &object.Function{
		Pattern: node.Pattern,
		Body:    node.Body,
		Context: ctx,
	}

	ctx.AddFunction(function)

	return NULL
}
