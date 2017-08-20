package main

import (
	"bufio"
	"fmt"
	"os"

	_ "github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/lexer"
	"github.com/Zac-Garby/pluto/parser"
	_ "github.com/Zac-Garby/pluto/token"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		next := lexer.Lexer(text)
		parse := parser.New(next)
		program := parse.Parse()

		fmt.Println(program.Tree())
	}
}
