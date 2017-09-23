package test

import (
	"testing"

	. "github.com/Zac-Garby/pluto/parser"
)

func TestPrecedence(t *testing.T) {
	cases := map[string]string{
		// Assignment
		"x = y + z":  "x = (y + z)",
		"x = y = 5":  "x = (y = 5)",
		"x += y + z": "x += (y + z)",

		// Question-mark operator
		"x ? y + z": "x ? (y + z)",
		"x ? y ? z": "(x ? y) ? z",

		// Logic
		"x || y && z": "x || (y && z)",
		"x && y || z": "(x && y) || z",
		"x | y & z":   "x | (y & z)",
		"x & y | z":   "(x & y) | z",

		// Comparisions
		"x == y < z":  "x == (y < z)",
		"x < y == z":  "(x < y) == z",
		"x != y >= z": "x != (y >= z)",
		"x >= y != z": "(x >= y) != z",

		// Mathematical operators
		"x + y * z":  "x + (y * z)",
		"x * y + z":  "(x * y) + z",
		"x - y / z":  "x - (y / z)",
		"x / y - z":  "(x / y) - z",
		"x * y ** z": "x * (y ** z)",
		"x ** y * z": "(x ** y) * z",

		// Prefixes
		"-x + y": "(-x) + (y)",
		"--x":    "-(-(x))",

		// Method calling
		"x:y z + foo": "(x:y z) + foo",
		"foo + x:y z": "foo + (x:y z)",

		// Indexing
		"x[y + z] - foo": "(x[(y + z)]) - foo",
		"x - y[z + foo]": "x - (y[(z + foo)])",
	}

	for left, right := range cases {
		var (
			lparser = New(left, "<test suite>")
			rparser = New(right, "<test suite>")
		)

		last := lparser.Parse()
		if len(lparser.Errors) > 0 {
			lparser.PrintErrors()
			t.Errorf("a parse error occured when parsing '%s'", left)
		}

		rast := rparser.Parse()
		if len(rparser.Errors) > 0 {
			rparser.PrintErrors()
			t.Errorf("a parse error occured when parsing '%s'", right)
		}

		var (
			ltree = last.Tree()
			rtree = rast.Tree()
		)

		if ltree != rtree {
			t.Errorf("'%s' and '%s' do not product identical ASTs", left, right)
		}
	}
}
