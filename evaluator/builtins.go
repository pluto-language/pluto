package evaluator

import (
	"fmt"
	"strings"

	"github.com/Zac-Garby/pluto/object"
)

type builtinFn func(map[string]object.Object, *object.Context) object.Object

type Builtin struct {
	Pattern []string
	Fn      builtinFn
}

func NewBuiltin(ptn string, types map[string]object.Type, fn builtinFn) Builtin {
	pattern := strings.Split(ptn, " ")

	typedFn := func(args map[string]object.Object, ctx *object.Context) object.Object {
		for key, t := range types {
			val := args[key]

			if val.Type() != t {
				return err(
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

var empty = make(map[string]object.Type)

var Builtins = []Builtin{
	NewBuiltin("print $obj", empty, printObj),
}

// print $obj
func printObj(args map[string]object.Object, ctx *object.Context) object.Object {
	fmt.Println(args["obj"])

	return NULL
}
