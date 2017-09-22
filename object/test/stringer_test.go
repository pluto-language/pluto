package object

import (
	"testing"

	. "github.com/Zac-Garby/pluto/object"
)

func TestStringer(t *testing.T) {
	var (
		n1 = &Number{Value: 0}
		n2 = &Number{Value: 1}

		s1 = &String{Value: "foo"}
		s2 = &String{Value: "bar"}

		c1 = &Char{Value: 'x'}
		c2 = &Char{Value: 'y'}

		a0 = &Array{Value: []Object{}}
		a1 = &Array{Value: []Object{n1}}
		a2 = &Array{Value: []Object{n1, n2}}

		t0 = &Tuple{Value: []Object{}}
		t1 = &Tuple{Value: []Object{n1}}
		t2 = &Tuple{Value: []Object{n1, n2}}

		m0 = &Map{Keys: make(map[string]Object), Values: make(map[string]Object)}
		m1 = &Map{Keys: make(map[string]Object), Values: make(map[string]Object)}
		m2 = &Map{Keys: make(map[string]Object), Values: make(map[string]Object)}
	)

	m1.Set(s1, n1)

	m2.Set(s1, n1)
	m2.Set(s2, n2)

	cases := []struct {
		obj Object
		str string
	}{
		{n1, "0"},
		{n2, "1"},
		{s1, "foo"},
		{s2, "bar"},
		{c1, "x"},
		{c2, "y"},
		{a0, "[]"},
		{a1, "[0]"},
		{a2, "[0, 1]"},
		{t0, "()"},
		{t1, "(0)"},
		{t2, "(0, 1)"},
		{m0, "[:]"},
		{m1, "[foo: 0]"},
		{m2, "[foo: 0, bar: 1]"},
	}

	for _, pair := range cases {
		str := pair.obj.String()

		if str != pair.str {
			t.Errorf("wrong string representation for %s, should be %s", str, pair.str)
		}
	}
}
