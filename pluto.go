package main

import (
	"fmt"
	"os"

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
	// bar = 100 + foo
	// bar
	in := []byte{10, 0, 0, 11, 0, 0, 25, 12, 0, 1, 11, 0, 1}

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

	store.Data = map[string]object.Object{
		"foo": &object.Number{Value: 14},
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

	fmt.Println(val)
}
