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

outer:
	for _, fn := range f.Functions {
		fnpat := fn.Pattern

		if len(fnpat) != len(pattern) {
			// Doesn't match
			continue outer
		}

		for i, item := range pattern {
			var (
				fItem = fnpat[i]
				isArg = item[0] == '$'
			)

			if isArg {
				if _, ok := fItem.(*ast.Parameter); !ok {
					// Doesn't match
					continue outer
				}
			} else {
				if id, ok := fItem.(*ast.Identifier); !ok {
					// Doesn't match
					continue outer
				} else if id.Value != item {
					// Doesn't match
					continue outer
				}
			}
		}

		return &fn
	}

	return nil
}
