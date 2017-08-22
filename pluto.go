package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/Zac-Garby/pluto/context"
	"github.com/Zac-Garby/pluto/lexer"
	_ "github.com/Zac-Garby/pluto/object"
	"github.com/Zac-Garby/pluto/parser"
	"github.com/jessevdk/go-flags"
)

const VERSION = "0.1.0"

type Options struct {
	Parse       bool `short:"p" long:"parse" description:"Just parse the input - don't execute it."`
	Tree        bool `short:"t" long:"tree" description:"Pretty-print the AST."`
	Interactive bool `short:"i" long:"interactive" description:"Enter interactive mode after the file has been run"`
	NoPrelude   bool `short:"n" long:"no-prelude" description:"Don't load the prelude. Probably a bad idea."`
	Version     bool `short:"v" long:"version" description:"Print the version then quit"`

	Args struct {
		File string
	} `positional-args:"yes"`
}

var opts Options

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		return
	}

	if opts.Version {
		fmt.Printf("Pluto v%s\n", VERSION)
		return
	}

	if len(opts.Args.File) == 0 {
		REPL()
	} else {
		executeFile(opts.Args.File)
	}
}

func REPL() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		execute(text)
	}
}

func executeFile(name string) {
	if code, err := ioutil.ReadFile(name); err != nil {
		panic(err)
	} else {
		execute(string(code))
	}
}

func execute(code string) {
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
}
