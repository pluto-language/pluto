package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Zac-Garby/pluto/evaluation"
	"github.com/Zac-Garby/pluto/parser"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
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
	if r, exists := os.LookupEnv("PLUTO"); exists {
		root = r
	} else {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}

		root = filepath.Join(usr.HomeDir, "pluto")
	}

	fmt.Println(root)

	if _, err := flags.Parse(&opts); err != nil {
		return
	}

	color.NoColor = opts.NoColour

	if opts.Version {
		fmt.Printf("Pluto v%s\n", version)
		return
	}

	if len(opts.Args.File) == 0 {
		runREPL(&evaluation.Context{
			Store: make(map[string]evaluation.Object),
		})
	} else {
		executeFile(opts.Args.File)
	}
}

func runREPL(ctx *evaluation.Context) {
	if !opts.NoPrelude {
		importPrelude(ctx)
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

		if !opts.NoPrelude {
			importPrelude(ctx)
		}

		execute(string(code), false, ctx)

		if opts.Interactive {
			runREPL(ctx)
		}
	}
}

func importPrelude(ctx *evaluation.Context) {
	srcPath := filepath.Join(root, "libraries", "prelude.pluto")

	if prelude, err := ioutil.ReadFile(srcPath); err != nil {
		panic(err)
	} else {
		oldTreeFlag := opts.Tree

		opts.Tree = false
		execute(string(prelude), false, ctx)
		opts.Tree = oldTreeFlag
	}
}

func execute(code string, showOutput bool, ctx *evaluation.Context) {
	parse := parser.New(code)
	program := parse.Parse()

	if len(parse.Errors) > 0 {
		parse.PrintErrors()

		return
	}

	if opts.Parse || opts.Tree {
		if opts.Tree {
			fmt.Println(program.Tree())
		}

		return
	}

	result := evaluation.EvaluateProgram(program, ctx)

	if (showOutput && !result.Equals(evaluation.O_NULL)) || evaluation.IsErr(result) {
		fmt.Println(result.String())
	}
}
