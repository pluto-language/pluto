package graph

import (
	"github.com/Zac-Garby/pluto/bytecode"
)

func (g *Grapher) addEdges() {
	type edger func(int, bytecode.Instruction)

	edges := map[byte]edger{
		bytecode.JumpIfFalse: g.condJumpEdge,
		bytecode.JumpIfTrue:  g.condJumpEdge,
		bytecode.Jump:        g.jumpEdge,
		bytecode.Break:       g.breakEdge,
		bytecode.Next:        g.nextEdge,
		bytecode.LoopStart:   g.loopStartEdge,
		bytecode.LoopEnd:     g.loopEndEdge,
	}

	for i, instr := range g.code {
		if i+1 >= len(g.code) {
			break
		}

		if edger, ok := edges[instr.Code]; ok {
			edger(i, instr)
		} else {
			g.graph.AddEdge(itoa(i), itoa(i+1), true, nil)
		}
	}

	g.graph.AddEdge(itoa(len(g.code)-1), itoa(len(g.code)), true, nil)
}

func (g *Grapher) condJumpEdge(i int, instr bytecode.Instruction) {
	pos := byteToInstructionIndex(g.code, int(instr.Arg))

	g.graph.AddEdge(itoa(i), itoa(pos+1), true, map[string]string{
		"color": "green",
		"label": `" yes"`,
	})

	g.graph.AddEdge(itoa(i), itoa(i+1), true, map[string]string{
		"color": "red",
		"label": `" no"`,
	})
}

func (g *Grapher) jumpEdge(i int, instr bytecode.Instruction) {
	pos := byteToInstructionIndex(g.code, int(instr.Arg))

	g.graph.AddEdge(itoa(i), itoa(pos+1), true, nil)
}

func (g *Grapher) breakEdge(i int, instr bytecode.Instruction) {
	top := g.breaks[len(g.breaks)-1]
	g.graph.AddEdge(itoa(i), itoa(top), true, nil)
}

func (g *Grapher) nextEdge(i int, instr bytecode.Instruction) {
	top := g.nexts[len(g.nexts)-1]
	g.graph.AddEdge(itoa(i), itoa(top), true, nil)
}

func (g *Grapher) loopStartEdge(i int, instr bytecode.Instruction) {
	g.nexts = append(g.nexts, i+1)

	var o int

	for o = i; g.code[o].Code != bytecode.LoopEnd; o++ {
		// Nothing here
	}

	g.breaks = append(g.breaks, o)

	g.graph.AddEdge(itoa(i), itoa(i+1), true, nil)
}

func (g *Grapher) loopEndEdge(i int, instr bytecode.Instruction) {
	g.breaks = g.breaks[:len(g.breaks)-1]
	g.nexts = g.nexts[:len(g.nexts)-1]

	g.graph.AddEdge(itoa(i), itoa(i+1), true, nil)
}
