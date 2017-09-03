package main

import (
	"fmt"
	"os"

	"github.com/Zac-Garby/pluto/ast"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
	"github.com/Zac-Garby/pluto/vm"
)

const version = "0.1.0"

type options struct {
	Parse       bool `short:"p" long:"parse" description:"Just parse the input - don't execute it."`
	Tree        bool `short:"t" long:"tree" description:"Pretty-print the AST."`
	Interactive bool `short:"i" long:"interactive" description:"Enter interactive mode after the file has been run"`
	NoPrelude   bool `short:"n" long:"no-prelude" description:"Don't load the prelude. Probably a bad idea."`
	NoColour    bool `short:"c" long:"no-colour" description:"Don't use coloured output."`
	Version     bool `short:"v" long:"version" description:"Print the version then quit"`

	Args struct {
		File string
	} `positional-args:"yes"`
}

var (
	opts options
	root string
)

func main() {
	in := []byte{10, 0, 0, 11, 0, 0, 25, 12, 0, 1, 11, 0, 1, 40, 0, 0}

	code, err := bytecode.Read(in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	machine := vm.New()

	store := vm.NewStore()

	store.Names = map[rune]string{
		0: "foo",
		1: "bar",
	}

	store.Patterns = map[rune]string{
		0: "print $obj",
	}

	store.Data = map[string]object.Object{
		"foo": &object.Number{Value: 14},
	}

	store.Functions.Functions = []object.Function{
		makePrintFunction(),
	}

	constants := []object.Object{
		&object.Number{Value: 3},
	}

	machine.Run(code, vm.NewStore(), store, constants)

	if machine.Error != nil {
		fmt.Println(machine.Error)
		return
	}

	val := machine.ExtractValue()

	fmt.Println(">>", val)
}

func makePrintFunction() object.Function {
	fn := object.Function{
		Pattern: []ast.Expression{
			&ast.Identifier{Value: "print"},
			&ast.Parameter{Name: "obj"},
		},
		Constants: []object.Object{
			object.NullObj,
		},
		Names: map[rune]string{
			0: "obj",
		},
	}

	// Prints the argument 'obj', then loads the constant 'null'
	bytes := []byte{11, 0, 0, 37, 10, 0, 0}

	code, err := bytecode.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fn.Body = code

	return fn
}
