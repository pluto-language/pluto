package evaluation

import (
	"fmt"
	"math"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
)

type args map[string]Object
type builtinFn func(args, *Context) Object

type Builtin struct {
	Pattern []string
	Fn      builtinFn
}

func NewBuiltin(ptn string, fn builtinFn, types map[string]Type) Builtin {
	pattern := strings.Split(ptn, " ")

	typedFn := func(args args, ctx *Context) Object {
		for key, t := range types {
			val := args[key]

			if !is(val, t) {
				return Err(
					ctx,
					"the $%s parameter of %s must be of type %s, not %s",
					"TypeError",
					key, ptn,
					t, val.Type(),
				)
			}
		}

		return fn(args, ctx)
	}

	return Builtin{
		Pattern: pattern,
		Fn:      typedFn,
	}
}

var empty = make(map[string]Type)

var builtins = []Builtin{}

func GetBuiltins() []Builtin {
	if len(builtins) == 0 {
		builtins = []Builtin{
			NewBuiltin("print $obj", printObj, empty),

			NewBuiltin("print $format with $args", printObjWithArgs, map[string]Type{
				"format": STRING,
				"args":   COLLECTION,
			}),

			NewBuiltin("do $block", doBlock, map[string]Type{
				"block": BLOCK,
			}),

			NewBuiltin("do $block with $args", doBlockWithArgs, map[string]Type{
				"block": BLOCK,
				"args":  COLLECTION,
			}),

			NewBuiltin("do $block on $arg", doBlockOnArg, map[string]Type{
				"block": BLOCK,
			}),

			NewBuiltin("map $block over $collection", mapBlockOverCollection, map[string]Type{
				"block":      BLOCK,
				"collection": COLLECTION,
			}),

			NewBuiltin("format $format with $args", formatWithArgs, map[string]Type{
				"format": STRING,
				"args":   COLLECTION,
			}),

			NewBuiltin("$start to $end", startToEnd, map[string]Type{
				"start": NUMBER,
				"end":   NUMBER,
			}),

			NewBuiltin(
				"slice $collection from $start to $end",
				sliceCollectionFromStartToEnd,
				map[string]Type{
					"collection": COLLECTION,
					"start":      NUMBER,
					"end":        NUMBER,
				},
			),

			NewBuiltin(
				"slice $collection from $start",
				sliceCollectionFromStart,
				map[string]Type{
					"collection": COLLECTION,
					"start":      NUMBER,
				},
			),

			NewBuiltin(
				"slice $collection to $end",
				sliceCollectionToEnd,
				map[string]Type{
					"collection": COLLECTION,
					"end":        NUMBER,
				},
			),

			NewBuiltin(
				"filter $collection by $predicate",
				filterCollectionByPredicate,
				map[string]Type{
					"collection": COLLECTION,
					"predicate":  BLOCK,
				},
			),

			NewBuiltin("round $number", roundNumber, map[string]Type{
				"number": NUMBER,
			}),

			NewBuiltin("floor $number", floorNumber, map[string]Type{
				"number": NUMBER,
			}),

			NewBuiltin("ceil $number", ceilNumber, map[string]Type{
				"number": NUMBER,
			}),

			NewBuiltin("keys of $map", keysOfMap, map[string]Type{
				"map": MAP,
			}),

			NewBuiltin("values of $map", valuesOfMap, map[string]Type{
				"map": MAP,
			}),

			NewBuiltin("pairs of $map", pairsOfMap, map[string]Type{
				"map": MAP,
			}),
		}
	}

	return builtins
}

// print $obj
func printObj(args args, ctx *Context) Object {
	fmt.Println(args["obj"])

	return O_NULL
}

// Not a builtin
// Formats a string according to args
func format(format string, args []Object) string {
	for _, f := range args {
		format = strings.Replace(format, "{}", f.String(), 1)
	}

	return format
}

// print $format with $args
func printObjWithArgs(args args, ctx *Context) Object {
	var (
		str   = args["format"].(*String)
		elems = args["args"].(Collection)
	)

	result := format(str.Value, elems.Elements())

	fmt.Println(result)

	return O_NULL
}

// format $format with $args
func formatWithArgs(args args, ctx *Context) Object {
	var (
		str   = args["format"].(*String)
		elems = args["args"].(Collection)
	)

	result := format(str.Value, elems.Elements())

	return &String{Value: result}
}

func evalBlock(block *Block, args []Object, ctx *Context) Object {
	if len(block.Params) != len(args) {
		return err(
			ctx,
			"wrong number of arguments applied to a block. expected %d, got %d", "TypeError",
			len(block.Params),
			len(args),
		)
	}

	apArgs := make(map[string]Object)

	for i, param := range block.Params {
		apArgs[param.(*ast.Identifier).Value] = args[i]
	}

	return eval(block.Body, ctx.EncloseWith(apArgs))
}

// do $block
func doBlock(args args, ctx *Context) Object {
	block := args["block"].(*Block)

	return evalBlock(block, []Object{}, ctx)
}

// do $block with $args
func doBlockWithArgs(args args, ctx *Context) Object {
	var (
		block = args["block"].(*Block)
		col   = args["args"].(Collection)
	)

	return evalBlock(block, col.Elements(), ctx)
}

// do $block on $arg
func doBlockOnArg(args args, ctx *Context) Object {
	var (
		block = args["block"].(*Block)
		arg   = args["arg"]
	)

	return evalBlock(block, []Object{arg}, ctx)
}

// map $block over $collection
func mapBlockOverCollection(args args, ctx *Context) Object {
	var (
		block = args["block"].(*Block)
		col   = args["collection"].(Collection)
	)

	var result []Object

	for i, item := range col.Elements() {
		mapped := evalBlock(block, []Object{
			&Number{Value: float64(i)},
			item,
		}, ctx)

		if isErr(mapped) {
			return mapped
		}

		result = append(result, mapped)
	}

	return MakeCollection(col.Type(), result, ctx)
}

// $start to $end
func startToEnd(args args, ctx *Context) Object {
	var (
		start = args["start"].(*Number)
		end   = args["end"].(*Number)

		sVal = int(start.Value)
		eVal = int(end.Value)
	)

	if eVal < sVal {
		result := &Array{Value: []Object{}}

		for i := sVal; i >= eVal; i-- {
			result.Value = append(result.Value, &Number{Value: float64(i)})
		}

		return result
	} else if eVal > sVal {
		result := &Array{Value: []Object{}}

		for i := sVal; i < eVal+1; i++ {
			result.Value = append(result.Value, &Number{Value: float64(i)})
		}

		return result
	}

	return &Array{Value: []Object{start}}
}

// slice $collection from $start to $end
func sliceCollectionFromStartToEnd(args args, ctx *Context) Object {
	var (
		col   = args["collection"].(Collection)
		start = args["start"].(*Number)
		end   = args["end"].(*Number)

		elems = col.Elements()
		sVal  = int(start.Value)
		eVal  = int(end.Value)
	)

	if sVal >= eVal {
		return err(ctx, "$start must be less than $end", "OutOfBoundsError")
	}

	if sVal < 0 || eVal < 0 {
		return err(ctx, "neither $start nor $end can be less than 0", "OutOfBoundsError")
	}

	if eVal >= len(elems) {
		return err(ctx, "$end must be contained by $collection", "OutOfBoundsError")
	}

	return &Array{Value: elems[sVal:eVal]}
}

// slice $collection from $start
func sliceCollectionFromStart(args args, ctx *Context) Object {
	var (
		col   = args["collection"].(Collection)
		start = args["start"].(*Number)

		elems = col.Elements()
		index = int(start.Value)
	)

	if index < 0 || index >= len(elems) {
		return err(ctx, "$start is out of bounds", "OutOfBoundsError")
	}

	return &Array{Value: elems[index:]}
}

// slice $collection to $end
func sliceCollectionToEnd(args args, ctx *Context) Object {
	var (
		col = args["collection"].(Collection)
		end = args["end"].(*Number)

		elems = col.Elements()
		index = int(end.Value)
	)

	if index < 0 || index >= len(elems) {
		return err(ctx, "$end is out of bounds", "OutOfBoundsError")
	}

	return &Array{Value: elems[:index]}
}

// filter $collection by $predicate
func filterCollectionByPredicate(args args, ctx *Context) Object {
	var (
		col  = args["collection"].(Collection)
		pred = args["predicate"].(*Block)

		filtered = []Object{}
	)

	for i, item := range col.Elements() {
		result := evalBlock(pred, []Object{
			&Number{Value: float64(i)},
			item,
		}, ctx)

		if isErr(result) {
			return result
		}

		if isTruthy(result) {
			filtered = append(filtered, item)
		}
	}

	return MakeCollection(col.Type(), filtered, ctx)
}

// round $number
func roundNumber(args args, ctx *Context) Object {
	num := args["number"].(*Number).Value

	return &Number{Value: math.Floor(num + 0.5)}
}

// floor $number
func floorNumber(args args, ctx *Context) Object {
	num := args["number"].(*Number).Value

	return &Number{Value: math.Floor(num)}
}

// ceil $number
func ceilNumber(args args, ctx *Context) Object {
	num := args["number"].(*Number).Value

	return &Number{Value: math.Ceil(num)}
}

// keys of $map
func keysOfMap(args args, ctx *Context) Object {
	var (
		m    = args["map"].(*Map)
		keys = []Object{}
	)

	for _, k := range m.Keys {
		keys = append(keys, k)
	}

	return &Array{Value: keys}
}

// values of $map
func valuesOfMap(args args, ctx *Context) Object {
	var (
		m    = args["map"].(*Map)
		vals = []Object{}
	)

	for _, v := range m.Values {
		vals = append(vals, v)
	}

	return &Array{Value: vals}
}

// pairs of $map
func pairsOfMap(args args, ctx *Context) Object {
	var (
		m     = args["map"].(*Map)
		keys  = m.Keys
		vals  = m.Values
		pairs = []Object{}
	)

	for hash, key := range keys {
		val := vals[hash]

		pairs = append(pairs, &Tuple{Value: []Object{key, val}})
	}

	return &Array{Value: pairs}
}
