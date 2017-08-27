package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Zac-Garby/pluto/evaluation"

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
	ctx := &evaluation.Context{
		Store: make(map[string]evaluation.Object),
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		execute(text, true, ctx)
	}
}

func executeFile(name string) {
	if code, err := ioutil.ReadFile(name); err != nil {
		panic(err)
	} else {
		ctx := &evaluation.Context{
			Store: make(map[string]evaluation.Object),
		}

		execute(string(code), false, ctx)
	}
}

func execute(code string, showOutput bool, ctx *evaluation.Context) {
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

	result := evaluation.EvaluateProgram(program, ctx)

	if showOutput && !result.Equals(evaluation.O_NULL) {
		fmt.Println(result.String())
	}
}
