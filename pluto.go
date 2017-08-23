package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Zac-Garby/pluto/evaluator"
	"github.com/Zac-Garby/pluto/object"

	"github.com/Zac-Garby/pluto/lexer"
	"github.com/Zac-Garby/pluto/parser"
	"github.com/jessevdk/go-flags"
)

const version = "0.1.0"

type options struct {
	Parse       bool `short:"p" long:"parse" description:"Just parse the input - don't execute it."`
	Tree        bool `short:"t" long:"tree" description:"Pretty-print the AST."`
	Interactive bool `short:"i" long:"interactive" description:"Enter interactive mode after the file has been run"`
	NoPrelude   bool `short:"n" long:"no-prelude" description:"Don't load the prelude. Probably a bad idea."`
	Version     bool `short:"v" long:"version" description:"Print the version then quit"`

	Args struct {
		File string
	} `positional-args:"yes"`
}

var opts options

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		return
	}

	if opts.Version {
		fmt.Printf("Pluto v%s\n", version)
		return
	}

	if len(opts.Args.File) == 0 {
		runREPL()
	} else {
		executeFile(opts.Args.File)
	}
}

func runREPL() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		execute(text, true)
	}
}

func executeFile(name string) {
	if code, err := ioutil.ReadFile(name); err != nil {
		panic(err)
	} else {
		execute(string(code), false)
	}
}

func execute(code string, showOutput bool) {
	next := lexer.Lexer(code)
	parse := parser.New(next)
	program := parse.Parse()

	if len(parse.Errors) > 0 {
		parse.PrintErrors()
		fmt.Println("\nExiting...")
		os.Exit(1)
	}

	if opts.Parse || opts.Tree {
		if opts.Tree {
			fmt.Println(program.Tree())
		}

		return
	}

	context := &object.Context{}
	result := evaluator.EvaluateProgram(program, context)

	if showOutput && result != evaluator.NULL {
		fmt.Println(result.String())
	}
}
