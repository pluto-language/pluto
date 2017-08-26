package object

import (
	"fmt"
	"strings"
)

type builtinFn func(map[string]Object, *Context) Object

type Builtin struct {
	Pattern []string
	Fn      builtinFn
}

func NewBuiltin(ptn string, types map[string]Type, fn builtinFn) Builtin {
	pattern := strings.Split(ptn, " ")

	typedFn := func(args map[string]Object, ctx *Context) Object {
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
	NewBuiltin("print $obj", empty, printObj),

	NewBuiltin("do $block", map[string]Type{
		"block": BLOCK,
	}, doBlock),

	NewBuiltin("do $block with $args", map[string]Type{
		"block": BLOCK,
		"args":  COLLECTION,
	}, doBlockWithArgs),

	NewBuiltin("do $block on $arg", map[string]Type{
		"block": BLOCK,
	}, doBlockOnArg),
}

// print $obj
func printObj(args map[string]Object, ctx *Context) Object {
	fmt.Println(args["obj"])

	return O_NULL
}

// do $block
func doBlock(args map[string]Object, ctx *Context) Object {
	block := args["block"].(*Block)

	return &AppliedBlock{
		Block: block,
		Args:  []Object{},
	}
}

// do $block with $args
func doBlockWithArgs(args map[string]Object, ctx *Context) Object {
	block := args["block"].(*Block)
	col := args["args"].(Collection)

	return &AppliedBlock{
		Block: block,
		Args:  col.Elements(),
	}
}

// do $block on $arg
func doBlockOnArg(args map[string]Object, ctx *Context) Object {
	block := args["block"].(*Block)
	arg := args["arg"]

	return &AppliedBlock{
		Block: block,
		Args:  []Object{arg},
	}
}
