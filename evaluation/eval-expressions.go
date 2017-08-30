package evaluation

import (
	"math"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
)

func evalIdentifier(node ast.Identifier, ctx *Context) Object {
	val := ctx.Get(node.Value)

	if val != nil {
		return val
	}

	return err(ctx, "`%s` is not found", "NotFoundError", node.Value)
}

func evalTuple(node ast.Tuple, ctx *Context) Object {
	elements := evalExpressions(node.Value, ctx)

	if len(elements) == 1 && isErr(elements[0]) {
		return elements[0]
	}

	return &Tuple{Value: elements}
}

func evalArray(node ast.Array, ctx *Context) Object {
	elements := evalExpressions(node.Elements, ctx)

	if len(elements) == 1 && isErr(elements[0]) {
		return elements[0]
	}

	return &Array{Value: elements}
}

func evalMap(node ast.Map, ctx *Context) Object {
	m := &Map{
		Values: make(map[string]Object),
		Keys:   make(map[string]Object),
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

func evalBlockLiteral(node ast.BlockLiteral, ctx *Context) Object {
	return &Block{Params: node.Params, Body: node.Body}
}

func evalPrefixExpression(node ast.PrefixExpression, ctx *Context) Object {
	right := eval(node.Right, ctx)
	if isErr(right) {
		return right
	}

	return evalPrefix(node.Operator, right, ctx)
}

func evalPrefix(op string, right Object, ctx *Context) Object {
	if instance, ok := right.(*Instance); ok {
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
		return err(ctx, "unknown operator: %s%s", "TypeError", op, right.Type)
	}
}

func evalInstancePrefix(op string, right *Instance, ctx *Context) Object {
	fnName, ok := prefixOverloads[op]
	if !ok {
		return err(ctx, "cannot overload operator %s", "NotFoundError", op)
	}

	if method := right.Base.(*Class).GetMethod(fnName); method != nil {
		args := map[string]Object{
			"self": right,
		}

		enclosed := ctx.EncloseWith(args)

		return eval(method.Fn.Body, enclosed)
	}

	return err(ctx, "unknown operator: %s%s. try overloading %s", "NotFoundError", op, right.Base.String(), fnName)
}

func evalMinusPrefix(right Object, ctx *Context) Object {
	if right.Type() != NUMBER {
		return err(ctx, "unknown operator: -%s", "TypeError", right.Type())
	}

	return &Number{Value: -right.(*Number).Value}
}

func evalInfixExpression(node ast.InfixExpression, ctx *Context) Object {
	left := eval(node.Left, ctx)
	if isErr(left) {
		return left
	}

	right := eval(node.Right, ctx)
	if isErr(right) {
		return right
	}

	op := node.Operator

	if lCol, ok := left.(Collection); ok {
		if rCol, ok := right.(Collection); ok {
			return evalCollectionInfix(op, lCol, rCol, ctx)
		}
	}

	if instance, ok := left.(*Instance); ok {
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
		if left == O_NULL {
			return right
		}

		return left
	}

	if left.Type() == NUMBER && right.Type() == NUMBER {
		return evalNumberInfix(op, left.(*Number), right.(*Number), ctx)
	}

	if (left.Type() == CHAR || left.Type() == STRING) &&
		(right.Type() == CHAR || right.Type() == STRING) {
		return evalCharStringInfix(op, left, right, ctx)
	}

	if left.Type() == CHAR && right.Type() == NUMBER {
		ch := string(left.(*Char).Value)
		amount := int(math.Floor(right.(*Number).Value))

		return &String{Value: strings.Repeat(ch, amount)}
	}

	if lCol, ok := left.(Collection); ok {
		if right.Type() == NUMBER && op == "*" {
			var result []Object
			elems := lCol.Elements()
			amount := int(math.Floor(right.(*Number).Value))

			for i := 0; i < amount; i++ {
				result = append(result, elems...)
			}

			return MakeCollection(left.Type(), result, ctx)
		}
	}

	return err(ctx, "unknown operator: %s %s %s", "TypeError", left.Type(), op, right.Type())
}

func evalCollectionInfix(op string, left, right Collection, ctx *Context) Object {
	l := left.Elements()
	r := right.Elements()

	switch op {
	case "==":
		return boolObj(left.Equals(right))
	case "!=":
		return boolObj(!left.Equals(right))
	case "+":
		return MakeCollection(left.Type(), append(l, r...), ctx)
	case "-":
		var elems []Object

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

		return MakeCollection(left.Type(), elems, ctx)
	case "&", "&&":
		var elems []Object

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

		return MakeCollection(left.Type(), elems, ctx)
	case "|", "||":
		var elems []Object

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

		return MakeCollection(left.Type(), elems, ctx)
	default:
		return err(ctx, "unknown operator: %s %s %s", "TypeError", left.Type(), op, right.Type())
	}
}

func evalNumberInfix(op string, left, right *Number, ctx *Context) Object {
	l := left.Value
	r := right.Value

	switch op {
	case "+":
		return &Number{Value: l + r}
	case "-":
		return &Number{Value: l - r}
	case "*":
		return &Number{Value: l * r}
	case "/":
		return &Number{Value: l / r}
	case "&":
		return &Number{Value: float64(int(l) & int(r))}
	case "|":
		return &Number{Value: float64(int(l) | int(r))}
	case "**":
		return &Number{Value: math.Pow(l, r)}
	case "//":
		return &Number{Value: math.Floor(l / r)}
	case "%":
		return &Number{Value: math.Mod(l, r)}
	case "<":
		return boolObj(l < r)
	case ">":
		return boolObj(l > r)
	case "<=":
		return boolObj(l <= r)
	case ">=":
		return boolObj(l >= r)
	default:
		return err(ctx, "unknown operator: %s %s %s", "TypeError", left.Type(), op, right.Type())
	}
}

func evalCharStringInfix(op string, left, right Object, ctx *Context) Object {
	var l, r string

	if lch, ok := left.(*Char); ok {
		l = string(lch.Value)
	} else if lstr, ok := left.(*String); ok {
		l = lstr.Value
	}

	if rch, ok := right.(*Char); ok {
		r = string(rch.Value)
	} else if rstr, ok := right.(*String); ok {
		r = rstr.Value
	}

	switch op {
	case "+":
		return &String{Value: l + r}
	case "-":
		var val string

		for _, ch := range strings.Split(l, "") {
			if ch != r {
				val += ch
			}
		}

		return &String{Value: val}
	default:
		return err(ctx, "unknown operator: %s %s %s", "TypeError", left.Type(), op, right.Type())
	}
}

func evalInstanceInfix(op string, left *Instance, right Object, ctx *Context) Object {
	fnName, ok := infixOverloads[op]
	if !ok {
		return err(ctx, "cannot overload operator %s", "NotFoundError", op)
	}

	if method := left.Base.(*Class).GetMethod(fnName); method != nil {
		methodPattern := method.Fn.Pattern

		args := map[string]Object{
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

func evalAssignExpression(node ast.AssignExpression, ctx *Context) Object {
	right := eval(node.Value, ctx)
	if isErr(right) {
		return right
	}

	if dot, ok := node.Name.(*ast.DotExpression); ok {
		o := eval(dot.Left, ctx)

		ctr, ok := o.(Container)
		if !ok {
			return err(ctx, "can only access fields of containers", "TypeError")
		}

		if id, ok := dot.Right.(*ast.Identifier); ok {
			ctr.Set(&String{Value: id.Value}, right)
			return right
		}

		return err(ctx, "an identifier is expected to follow a dot operator", "SyntaxError")
	}

	if left, ok := node.Name.(*ast.Identifier); ok {
		ctx.Assign(left.Value, right)
		return right
	}

	if left, ok := node.Name.(*ast.IndexExpression); ok {
		rightobj := eval(left.Index, ctx)
		if isErr(rightobj) {
			return rightobj
		}

		leftobj := eval(left.Collection, ctx)
		if isErr(leftobj) {
			return leftobj
		}

		if col, ok := leftobj.(Collection); ok {
			if index, ok := rightobj.(*Number); ok {
				col.SetIndex(int(index.Value), right)
				return O_NULL
			}

			return err(ctx, "cannot index a collection with %s - expected a <number>", "TypeError", rightobj.Type())
		}

		if cont, ok := leftobj.(Container); ok {
			cont.Set(rightobj, right)
			return O_NULL
		}

		return err(ctx, "can only index collections and containers", "TypeError")
	}

	return err(ctx, "can only assign to identifiers and index expressions!", "SyntaxError")
}

func evalDeclareExpression(node ast.DeclareExpression, ctx *Context) Object {
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

func evalDotExpression(node ast.DotExpression, ctx *Context) Object {
	left := eval(node.Left, ctx)

	if field, ok := node.Right.(*ast.Identifier); ok {
		if cnt, ok := left.(Container); ok {
			return cnt.Get(&String{Value: field.Value})
		}

		return err(ctx, "cannot access fields of %s", "TypeError", left.Type())
	}

	return err(ctx, "an identifier is expected after a dot '.'", "SyntaxError")
}

func evalWhileLoop(node ast.WhileLoop, ctx *Context) Object {
	var steps []Object

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

		if result.Type() == BREAK {
			break
		}

		steps = append(steps, result)
	}

	return &Array{Value: steps}
}

func evalForLoop(node ast.ForLoop, ctx *Context) Object {
	var (
		v    = node.Var
		body = node.Body
		col  = eval(node.Collection, ctx)
	)

	if isErr(col) {
		return col
	}

	var items []Object
	if collection, ok := col.(Collection); ok {
		items = collection.Elements()
	} else {
		return err(ctx, "cannot perform a for-loop over the non-collection type: %s", "TypeError", col.Type())
	}

	var steps []Object

	for _, item := range items {
		enclosed := ctx.EncloseWith(map[string]Object{
			v.Token().Literal: item,
		})

		result := eval(body, enclosed)
		if isErr(result) {
			return result
		}

		if result.Type() == BREAK {
			break
		}

		steps = append(steps, result)
	}

	return &Array{Value: steps}
}

func evalFunctionCall(node ast.FunctionCall, ctx *Context) Object {
	fn := ctx.GetFunction(node.Pattern)

	if fn == nil {
		ps := patternString(node.Pattern)

		return err(ctx, "no function matching the pattern: %s", "NotFoundError", ps)
	}

	args := make(map[string]Object)

	var result Object

	if function, ok := fn.(*Function); ok {
		result = applyFunction(function, node.Pattern, ctx)
	}

	if function, ok := fn.(Builtin); ok {
		for i, item := range node.Pattern {
			fItem := function.Pattern[i]

			if arg, ok := item.(*ast.Argument); ok {
				if fItem[0] == '$' {
					evaled := eval(arg.Value, ctx)
					if isErr(evaled) {
						return evaled
					}

					args[fItem[1:]] = evaled
				}
			}
		}

		result = function.Fn(args, ctx)
	}

	if result == nil {
		return O_NULL
	}

	return unwrapReturnValue(result)
}

func evalQualifiedFunctionCall(node ast.QualifiedFunctionCall, ctx *Context) Object {
	pkgID, ok := node.Package.(*ast.Identifier)
	if !ok {
		return err(ctx, "the package name (before ::) must be an identifier", "SyntaxError")
	}

	pkgName := pkgID.Value

	pkg, ok := ctx.Packages[pkgName]
	if !ok {
		return err(ctx, "a package named '%s' cannot be found. are you sure you've imported it?", "NotFoundError", pkgName)
	}

	fn := pkg.GetFunction(node.Pattern)

	if fn == nil {
		ps := patternString(node.Pattern)

		return err(ctx, "no function in package '%s' matching the pattern: %s", "NotFoundError", pkgName, ps)
	}

	if function, ok := fn.(*Function); ok {
		return applyFunction(function, node.Pattern, ctx)
	}

	return O_NULL
}

func patternString(pattern []ast.Expression) string {
	var patternStrings []string

	for _, item := range pattern {
		if id, ok := item.(*ast.Identifier); ok {
			patternStrings = append(patternStrings, id.Value)
		} else {
			patternStrings = append(patternStrings, "$")
		}
	}

	return strings.Join(patternStrings, " ")
}

func applyFunction(fn *Function, pattern []ast.Expression, ctx *Context) Object {
	args := make(map[string]Object)

	for i, item := range pattern {
		fItem := fn.Pattern[i]

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

	enclosed := fn.Context.EncloseWith(args)

	if fn.OnCall != nil {
		result := fn.OnCall(fn, ctx, enclosed)

		if result != nil || isErr(result) {
			return result
		}
	}

	return unwrapReturnValue(eval(fn.Body, enclosed))
}

func evalFunctionDefinition(node ast.FunctionDefinition, ctx *Context) Object {
	function := &Function{
		Pattern: node.Pattern,
		Body:    node.Body,
		Context: ctx,
	}

	ctx.AddFunction(function)

	return O_NULL
}

func evalIfExpression(node ast.IfExpression, ctx *Context) Object {
	condition := eval(node.Condition, ctx)
	if isErr(condition) {
		return condition
	}

	if isTruthy(condition) {
		return eval(node.Consequence, ctx)
	} else if node.Alternative != nil {
		return eval(node.Alternative, ctx)
	} else {
		return O_NULL
	}
}

func evalMatchExpression(node ast.MatchExpression, ctx *Context) Object {
	val := eval(node.Exp, ctx)
	if isErr(val) {
		return val
	}

	var matched ast.Statement

	for _, arm := range node.Arms {
		var (
			exprs = arm.Exprs
			body  = arm.Body

			m = false
		)

		if exprs == nil {
			m = true
		} else {
			for _, expr := range exprs {
				e := eval(expr, ctx)
				if isErr(e) {
					return e
				}

				if e.Equals(val) {
					m = true
				}
			}
		}

		if m {
			matched = body
			break
		}
	}

	if matched != nil {
		result := eval(matched, ctx.Enclose())
		return unwrapReturnValue(result)
	}

	return O_NULL
}

func evalMethodCall(node ast.MethodCall, ctx *Context) Object {
	inst := eval(node.Instance, ctx)
	if isErr(inst) {
		return inst
	}

	instance, ok := inst.(*Instance)
	if !ok {
		return err(ctx, "can only call a method on type <instance>. got %s", "TypeError", inst.Type())
	}

	base := instance.Base.(*Class)

	var (
		pattern  = node.Pattern
		function *Method
	)

	for _, fn := range base.GetMethods() {
		if len(pattern) != len(fn.Fn.Pattern) {
			continue
		}

		matched := true

		for i, item := range pattern {
			fItem := fn.Fn.Pattern[i]

			if itemID, ok := item.(*ast.Identifier); ok {
				if fItemID, ok := fItem.(*ast.Identifier); ok {
					if itemID.Value != fItemID.Value {
						matched = false
					}
				}
			} else if _, ok := item.(*ast.Argument); !ok {
				matched = false
			} else if _, ok := fItem.(*ast.Parameter); !ok {
				matched = false
			}
		}

		if matched {
			function = &fn
			break
		}
	}

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

		return err(
			ctx,
			"could not find a method of %s matching %s",
			"NotFoundError",
			base.Name,
			patternString,
		)
	}

	args := make(map[string]Object)

	for i, item := range node.Pattern {
		fItem := function.Fn.Pattern[i]

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

	args["self"] = instance
	enclosed := ctx.EncloseWith(args)

	return eval(function.Fn.Body, enclosed)
}

func evalTryExpression(node ast.TryExpression, ctx *Context) Object {
	v := eval(node.Body, ctx)

	val, ok := v.(*Instance)
	if !ok || val.Base.(*Class).Name != "Error" {
		return err(ctx, "the error value in a try-expression must be an Error instance", "TypeError")
	}

	tag, ok := val.Data["tag"]
	if !ok {
		return err(ctx, "error value doesn't have a 'tag' field", "NotFoundError")
	}

	msg, ok := val.Data["msg"]
	if !ok {
		return err(ctx, "error value doesn't have a 'msg' field", "NotFoundError")
	}

	var matched ast.Statement

	for _, arm := range node.Arms {
		var (
			exprs = arm.Exprs
			body  = arm.Body

			m = false
		)

		if exprs == nil {
			m = true
		} else {
			for _, expr := range exprs {
				e := eval(expr, ctx)
				if isErr(e) {
					return e
				}

				if e.Type() != STRING {
					return err(
						ctx,
						"all catch-arm predicate values must be strings. found %s",
						"TypeError",
						e.Type(),
					)
				}

				if e.Equals(tag) {
					m = true
				}
			}
		}

		if m {
			matched = body
			break
		}
	}

	if matched != nil {
		errObj := &Map{
			Values: make(map[string]Object),
			Keys:   make(map[string]Object),
		}

		errObj.Set(&String{Value: "tag"}, tag)
		errObj.Set(&String{Value: "msg"}, msg)

		enclosed := ctx.EncloseWith(map[string]Object{
			node.ErrName.(*ast.Identifier).Value: errObj,
		})

		r := eval(matched, enclosed)
		return unwrapReturnValue(r)
	}

	return val
}

func evalIndexExpression(node ast.IndexExpression, ctx *Context) Object {
	var (
		left  = node.Collection
		right = node.Index
	)

	leftobj := eval(left, ctx)
	if isErr(leftobj) {
		return leftobj
	}

	rightobj := eval(right, ctx)
	if isErr(rightobj) {
		return rightobj
	}

	if col, ok := leftobj.(Collection); ok {
		if index, ok := rightobj.(*Number); ok {
			return col.GetIndex(int(index.Value))
		}

		return err(ctx, "cannot index a collection with %s - expected a <number>", "TypeError", rightobj.Type())
	}

	if cont, ok := leftobj.(Container); ok {
		return cont.Get(rightobj)
	}

	return err(ctx, "can only index collections and containers", "TypeError")
}
