package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	_ "github.com/Zac-Garby/pluto/bytecode"

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

	if _, err := flags.Parse(&opts); err != nil {
		return
	}

	color.NoColor = opts.NoColour

	if opts.Version {
		fmt.Printf("Pluto v%s\n", version)
		return
	}
}
