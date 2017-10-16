package vm

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/compiler"
	"github.com/Zac-Garby/pluto/dir"
	"github.com/Zac-Garby/pluto/parser"
)

// Use imports the sources found by the glob src into
// the frame
func (f *Frame) Use(src string) {
	sources, err := dir.LocateAnySources(src)
	if err != nil {
		f.vm.Error = Err(err.Error(), ErrUnknown)
		return
	}

	for _, source := range sources {
		src, err := ioutil.ReadFile(source)
		if err != nil {
			f.vm.Error = Err(err.Error(), ErrUnknown)
		}

		var (
			str   = string(src)
			cmp   = compiler.New()
			parse = parser.New(str, source)
			prog  = parse.Parse()
		)

		if len(parse.Errors) > 0 {
			parse.PrintErrors()
			f.vm.Error = Err("parse error", ErrSyntax)
			return
		}

		err = cmp.CompileProgram(prog)
		if err != nil {
			f.vm.Error = Err(err.Error(), ErrUnknown)
			return
		}

		code, err := bytecode.Read(cmp.Bytes)
		if err != nil {
			f.vm.Error = Err(err.Error(), ErrUnknown)
			return
		}

		store := &Store{
			Names:    cmp.Names,
			Patterns: cmp.Patterns,
			Data:     []*item{},
			FunctionStore: &FunctionStore{
				Functions: cmp.Functions,
			},
		}

		machine := New()
		machine.Run(code, store, cmp.Constants, false)

		sourceName := strings.Split(filepath.Base(source), ".")[0]

		f.locals.ImportModule(store, sourceName)
	}
}
