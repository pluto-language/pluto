package evaluation

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
)

type args map[string]Object
type builtinFn func(args, *Context) Object

// Builtin represents a builtin function, i.e. a
// function which is implemented in Go, as opposed
// to Pluto
type Builtin struct {
	Pattern []string
	Fn      builtinFn
}

// NewBuiltin returns a builtin with the given pattern and function,
// and performs type checking on the arguments
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

// GetBuiltins returns a slice of all defined builtins,
// and defines them if they aren't already
func GetBuiltins() []Builtin {
	if len(builtins) == 0 {
		builtins = []Builtin{
			NewBuiltin("print $obj", printObj, empty),

			NewBuiltin("print $format with $args", printObjWithArgs, map[string]Type{
				"format": StringType,
				"args":   CollectionType,
			}),

			NewBuiltin("do $block", doBlock, map[string]Type{
				"block": BlockType,
			}),

			NewBuiltin("do $block with $args", doBlockWithArgs, map[string]Type{
				"block": BlockType,
				"args":  CollectionType,
			}),

			NewBuiltin("do $block on $arg", doBlockOnArg, map[string]Type{
				"block": BlockType,
			}),

			NewBuiltin("map $block over $collection", mapBlockOverCollection, map[string]Type{
				"block":      BlockType,
				"collection": CollectionType,
			}),

			NewBuiltin("format $format with $args", formatWithArgs, map[string]Type{
				"format": StringType,
				"args":   CollectionType,
			}),

			NewBuiltin("$start to $end", startToEnd, map[string]Type{
				"start": NumberType,
				"end":   NumberType,
			}),

			NewBuiltin(
				"slice $collection from $start to $end",
				sliceCollectionFromStartToEnd,
				map[string]Type{
					"collection": CollectionType,
					"start":      NumberType,
					"end":        NumberType,
				},
			),

			NewBuiltin(
				"slice $collection from $start",
				sliceCollectionFromStart,
				map[string]Type{
					"collection": CollectionType,
					"start":      NumberType,
				},
			),

			NewBuiltin(
				"slice $collection to $end",
				sliceCollectionToEnd,
				map[string]Type{
					"collection": CollectionType,
					"end":        NumberType,
				},
			),

			NewBuiltin(
				"filter $collection by $predicate",
				filterCollectionByPredicate,
				map[string]Type{
					"collection": CollectionType,
					"predicate":  BlockType,
				},
			),

			NewBuiltin("round $number", roundNumber, map[string]Type{
				"number": NumberType,
			}),

			NewBuiltin("floor $number", floorNumber, map[string]Type{
				"number": NumberType,
			}),

			NewBuiltin("ceil $number", ceilNumber, map[string]Type{
				"number": NumberType,
			}),

			NewBuiltin("keys of $map", keysOfMap, map[string]Type{
				"map": MapType,
			}),

			NewBuiltin("values of $map", valuesOfMap, map[string]Type{
				"map": MapType,
			}),

			NewBuiltin("pairs of $map", pairsOfMap, map[string]Type{
				"map": MapType,
			}),

			NewBuiltin("prompt $prefix", promptPrefix, map[string]Type{
				"prefix": StringType,
			}),

			NewBuiltin("type of $instance", typeOfInstance, map[string]Type{
				"instance": InstanceType,
			}),

			NewBuiltin("new $class", newClass, map[string]Type{
				"class": ClassType,
			}),
		}
	}

	return builtins
}

// print $obj
func printObj(args args, ctx *Context) Object {
	fmt.Println(args["obj"])

	return NullObj
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

	return NullObj
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

// prompt $prefix
func promptPrefix(args args, ctx *Context) Object {
	var (
		prefix = args["prefix"].(*String)
		reader = bufio.NewReader(os.Stdin)
	)

	fmt.Print(prefix.Value)

	if text, err := reader.ReadString('\n'); err != nil {
		panic(err)
	} else {
		return &String{Value: strings.TrimSpace(text)}
	}
}

func evalBlock(block *Block, args []Object, ctx *Context) Object {
	if len(block.Params) != len(args) {
		var params []string

		for _, param := range block.Params {
			params = append(params, param.(*ast.Identifier).Value)
		}

		return Err(
			ctx,
			"wrong number of arguments applied to a block. expected %d, got %d (params: %s)", "TypeError",
			len(block.Params),
			len(args),
			strings.Join(params, ", "),
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

		if IsErr(mapped) {
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
		return Err(ctx, "$start must be less than $end", "OutOfBoundsError")
	}

	if sVal < 0 || eVal < 0 {
		return Err(ctx, "neither $start nor $end can be less than 0", "OutOfBoundsError")
	}

	if eVal >= len(elems) {
		return Err(ctx, "$end must be contained by $collection", "OutOfBoundsError")
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
		return Err(ctx, "$start is out of bounds", "OutOfBoundsError")
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
		return Err(ctx, "$end is out of bounds", "OutOfBoundsError")
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

		if IsErr(result) {
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

// type of $instance
func typeOfInstance(args args, ctx *Context) Object {
	return args["instance"].(*Instance).Base
}

// new $class
func newClass(args args, ctx *Context) Object {
	class := args["class"].(*Class)

	return &Instance{
		Base: class,
		Data: make(map[string]Object),
	}
}
