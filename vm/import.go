package vm

import (
	"io/ioutil"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/compiler"
	"github.com/Zac-Garby/pluto/dir"
	"github.com/Zac-Garby/pluto/parser"
	"github.com/Zac-Garby/pluto/store"
)

// Use imports the sources found by the glob src into
// the frame
func (f *Frame) Use(src string) {
	sources, err := dir.LocateAnySources(src)
	if err != nil {
		f.vm.Error = Err(err.Error(), ErrUnknown)
		return
	}

	mergedTrees := ast.Program{}

	for _, source := range sources {
		src, err := ioutil.ReadFile(source)
		if err != nil {
			f.vm.Error = Err(err.Error(), ErrUnknown)

			return
		}

		var (
			str   = string(src)
			parse = parser.New(str, source)
			prog  = parse.Parse()
		)

		if len(parse.Errors) > 0 {
			parse.PrintErrors()
			f.vm.Error = Err("parse error", ErrSyntax)

			return
		}

		mergedTrees.Statements = append(mergedTrees.Statements, prog.Statements...)
	}

	cmp := compiler.New()

	if err = cmp.CompileProgram(mergedTrees); err != nil {
		f.vm.Error = Err(err.Error(), ErrUnknown)

		return
	}

	code, err := bytecode.Read(cmp.Bytes)
	if err != nil {
		f.vm.Error = Err(err.Error(), ErrUnknown)

		return
	}

	if err != nil {
		f.vm.Error = Err(err.Error(), ErrUnknown)

		return
	}

	store := &store.Store{
		Names:    cmp.Names,
		Patterns: cmp.Patterns,
		FunctionStore: &store.FunctionStore{
			Functions: cmp.Functions,
		},
	}

	machine := New()
	machine.Run(code, store, cmp.Constants, false)

	f.locals.ImportModule(store, src)
}
