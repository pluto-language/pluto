package main

import (
	"fmt"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/compiler"
	"github.com/Zac-Garby/pluto/graph"
	"github.com/Zac-Garby/pluto/object"
	"github.com/Zac-Garby/pluto/parser"
	"github.com/Zac-Garby/pluto/vm"
)

func main() {
	store := vm.NewStore()
	execute(`
a = 0

while (a < 10) {
	a = a + 1

	if (a > 5) {
		break
	}
}
`, store)

	/* for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")

		if obj, err := execute(text, store); err != nil {
			color.Red("  %s", err)
		} else if obj != nil {
			color.Cyan("  %s", obj)
		}
	} */
}

func execute(text string, store *vm.Store) (object.Object, error) {
	var (
		cmp   = compiler.New()
		parse = parser.New(text)
		prog  = parse.Parse()
	)

	if len(parse.Errors) > 0 {
		parse.PrintErrors()
		return nil, nil
	}

	err := cmp.CompileProgram(prog)
	if err != nil {
		return nil, err
	}

	code, err := bytecode.Read(cmp.Bytes)
	if err != nil {
		return nil, err
	}

	graph := graph.New(code, cmp.Constants)
	dot, err := graph.Render()
	if err != nil {
		return nil, err
	}

	fmt.Println(dot)

	store.Names = cmp.Names
	store.Functions.Functions = append(cmp.Functions, store.Functions.Functions...)
	store.Patterns = cmp.Patterns

	machine := vm.New()
	machine.Run(code, store, cmp.Constants)

	if machine.Error != nil {
		return nil, machine.Error
	}

	return machine.ExtractValue(), nil
}
