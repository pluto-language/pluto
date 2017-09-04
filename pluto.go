package main

import (
	"fmt"
	"os"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/compiler"
	"github.com/Zac-Garby/pluto/parser"
	"github.com/Zac-Garby/pluto/vm"
)

func main() {
	compiler := compiler.New()

	p := parser.New(`

a = -5
a

`)
	program := p.Parse()

	if len(p.Errors) > 0 {
		p.PrintErrors()
		os.Exit(1)
	}

	compiler.CompileProgram(program)

	code, err := bytecode.Read(compiler.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	store := vm.NewStore()
	store.Names = compiler.Names

	machine := vm.New()
	machine.Run(code, store, compiler.Constants)

	if machine.Error != nil {
		fmt.Println(machine.Error)
		return
	}

	val := machine.ExtractValue()

	if val != nil {
		fmt.Println(">>", val)
	}
}
