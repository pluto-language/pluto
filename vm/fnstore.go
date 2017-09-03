package vm

import (
	"strings"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/object"
)

// FunctionStore stores the functions in a frame,
// and handles the searching of them.
type FunctionStore struct {
	Functions []object.Function
}

// SearchString searches a function store for a function
// matching the given pattern. The pattern is in the
// format: "print $ and $".
func (f *FunctionStore) SearchString(search string) *object.Function {
	pattern := strings.Split(search, " ")

	for _, fn := range f.Functions {
		fnpat := fn.Pattern

		for i, item := range pattern {
			var (
				fItem = fnpat[i]
				isArg = item[0] == '$'
			)

			if isArg {
				if _, ok := fItem.(*ast.Argument); !ok {
					// Doesn't match
					goto nomatch
				}
			} else {
				if id, ok := fItem.(*ast.Identifier); !ok {
					// Doesn't match
					goto nomatch
				} else if id.Value != item {
					// Doesn't match
					goto nomatch
				}
			}
		}

		return &fn

	nomatch:
		continue
	}

	return nil
}
