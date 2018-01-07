package store

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

// def defines fn in the function store. If it's
// already defined (i.e. a function with the same pattern
// already exists) the old function is overwritten.
func (f *FunctionStore) def(newf object.Function) {
outer:
	for i, fn := range f.Functions {
		var (
			fnpat = fn.Pattern
			nfpat = newf.Pattern
		)

		if len(fnpat) != len(nfpat) {
			// Doesn't match
			continue outer
		}

		for i, item := range nfpat {
			fItem := fnpat[i]

			_, isArg := item.(*ast.Argument)

			if isArg {
				if _, ok := fItem.(*ast.Parameter); !ok {
					// Doesn't match
					continue outer
				}
			} else {
				if id, ok := fItem.(*ast.Identifier); !ok {
					// Doesn't match
					continue outer
				} else if id.Value != item.Token().Literal {
					// Doesn't match
					continue outer
				}
			}
		}

		f.Functions[i] = newf

		return
	}

	f.Functions = append(f.Functions, newf)
}

// Define defines all functions in fs in the store. If
// any are already defined, the new one will overwrite
// the old one.
func (f *FunctionStore) Define(fs ...object.Function) {
	for _, fn := range fs {
		f.def(fn)
	}
}

// Clone duplicates a function store
func (f *FunctionStore) Clone() *FunctionStore {
	nfs := &FunctionStore{}

	for _, fn := range f.Functions {
		nfs.Functions = append(nfs.Functions, fn)
	}

	return nfs
}
