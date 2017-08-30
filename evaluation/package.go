package evaluation

import "github.com/Zac-Garby/pluto/ast"

// Package represents an imported package
type Package struct {
	Context *Context // the context containing the package's functions
	Used    bool     // whether 'use <this package>' has been called
	Sources []string // unglobbed source files

	/* metadata read from 'collections.yaml' */
	Meta struct {
		Title, Description, Version string
		Tags, Authors, Sources      []string
	}
}

// GetFunction directly calls the context's GetFunction()
func (p *Package) GetFunction(ptn []ast.Expression) interface{} {
	return p.Context.GetFunction(ptn)
}
