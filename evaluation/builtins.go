package evaluation

import (
	"fmt"
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
		}
	}

	return builtins
}

// print $obj
func printObj(args args, ctx *Context) Object {
	fmt.Println(args["obj"])

	return O_NULL
}

// format $format with $args
func formatWithArgs(args args, ctx *Context) Object {
	var (
		format  = args["format"].(*String)
		formats = args["args"].(Collection)
	)

	// if format = "Hello, {}!" and args = ["world"]
	// the result will be "Hello, world!"

	result := format.Value

	for _, f := range formats.Elements() {
		result = strings.Replace(result, "{}", f.String(), 1)
	}

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
