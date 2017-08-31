package evaluation

import (
	"strings"
)

// MakeCollection creates a collection of type t with the elements provided
func MakeCollection(t Type, elems []Object, ctx *Context) Object {
	switch t {
	case ArrayType:
		return &Array{Value: elems}
	case TupleType:
		return &Tuple{Value: elems}
	case StringType:
		var strs []string

		for _, elem := range elems {
			strs = append(strs, elem.String())
		}

		return &String{Value: strings.Join(strs, "")}
	default:
		return Err(ctx, "could not form a collection of type %s", "TypeError", t)
	}
}
