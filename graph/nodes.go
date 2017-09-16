package graph

import (
	"fmt"

	"github.com/Zac-Garby/pluto/bytecode"
)

func (g *Grapher) addNodes() {
	g.graph.AddNode(name, "title", map[string]string{
		"color": "green",
		"shape": "box",
		"style": `"rounded, bold"`,
		"label": fmt.Sprintf("\"%s\"", name),
	})

	for i, instr := range g.code {
		var label string

		if bytecode.Instructions[instr.Code].HasArg {
			label = fmt.Sprintf("\"%d | %s | %d\"", i, instr.Name, int(instr.Arg))
		} else {
			label = fmt.Sprintf("\"%d | %s\"", i, instr.Name)
		}

		attr := map[string]string{
			"shape": shape,
			"label": label,
		}

		g.graph.AddNode(name, itoa(i), attr)
	}

	g.graph.AddNode(name, itoa(len(g.code)), map[string]string{
		"color": "red",
		"shape": "box",
		"style": `"rounded, bold"`,
		"label": "end",
	})
}
