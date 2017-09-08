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

a = [1, 2, 3, 4, 5]
a[2]

`)
	program := p.Parse()

	if len(p.Errors) > 0 {
		p.PrintErrors()
		os.Exit(1)
	}

	err := compiler.CompileProgram(program)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	code, err := bytecode.Read(compiler.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(code)

	store := vm.NewStore()
	store.Names = compiler.Names
	store.Functions.Functions = compiler.Functions
	store.Patterns = compiler.Patterns

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
