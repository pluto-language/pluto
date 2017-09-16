package graph

import (
	"strconv"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"

	gv "github.com/awalterschulze/gographviz"
)

var itoa = strconv.Itoa

const (
	name  = "<main>"
	shape = "record"
)

// Grapher emits dot code to render some bytecode.
type Grapher struct {
	code  bytecode.Code
	graph *gv.Graph

	breaks, nexts []int
	constants     []object.Object
}

// New initialises a new grapher.
func New(code bytecode.Code, constants []object.Object) *Grapher {
	return &Grapher{
		code:   code,
		breaks: []int{},
		nexts:  []int{},
	}
}

// Render generates some dot code to render a
// graph of the given bytecode.
func (g *Grapher) Render() (string, error) {
	g.graph = gv.NewGraph()
	g.graph.SetName(name)
	g.graph.SetDir(true)

	g.graph.AddEdge("title", "0", true, nil)

	g.addNodes()
	g.addEdges()

	return g.graph.String(), nil
}
