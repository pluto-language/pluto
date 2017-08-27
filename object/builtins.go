package object

import (
	"fmt"
	"strings"
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

var (
	O_NULL  = new(Null)
	O_TRUE  = &Boolean{Value: true}
	O_FALSE = &Boolean{Value: false}

	empty = make(map[string]Type)
)

var Builtins = []Builtin{
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
}

// print $obj
func printObj(args args, ctx *Context) Object {
	fmt.Println(args["obj"])

	return O_NULL
}

// do $block
func doBlock(args args, ctx *Context) Object {
	block := args["block"].(*Block)

	return block
}

// do $block with $args
func doBlockWithArgs(args args, ctx *Context) Object {
	var (
		block = args["block"].(*Block)
		// col   = args["args"].(Collection)
	)

	return block
}

// do $block on $arg
func doBlockOnArg(args args, ctx *Context) Object {
	var (
		block = args["block"].(*Block)
		// arg   = args["arg"]
	)

	return block
}
