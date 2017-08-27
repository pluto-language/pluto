package evaluation

import (
	"strings"
)

func MakeCollection(t Type, elems []Object, ctx *Context) Object {
	switch t {
	case ARRAY:
		return &Array{Value: elems}
	case TUPLE:
		return &Tuple{Value: elems}
	case STRING:
		var strs []string

		for _, elem := range elems {
			strs = append(strs, elem.String())
		}

		return &String{Value: strings.Join(strs, "")}
	default:
		return Err(ctx, "could not form a collection of type %s", "TypeError", t)
	}
}
