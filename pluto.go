package main

import (
	"fmt"
	"os"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/compiler"
	"github.com/Zac-Garby/pluto/vm"
)

func main() {
	compiler := compiler.New()

	compiler.CompileExpression(&ast.InfixExpression{
		Left:     &ast.Number{Value: 6},
		Right:    &ast.Number{Value: 5},
		Operator: "==",
	})

	code, err := bytecode.Read(compiler.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	machine := vm.New()
	machine.Run(code, vm.NewStore(), vm.NewStore(), compiler.Constants)

	if machine.Error != nil {
		fmt.Println(machine.Error)
		return
	}

	val := machine.ExtractValue()

	fmt.Println(">>", val)
}
