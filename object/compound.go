package object

import (
	"fmt"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
)

/* Structs */
type (
	/* Collections and collection-likes */
	Tuple struct {
		Value []Object
	}

	Array struct {
		Value []Object
	}

	Map struct {
		Pairs map[Object]Object
	}

	/* Others */
	Block struct {
		Params []ast.Expression
		Body   ast.Statement
	}

	Class struct {
		Name    string
		Parent  Object
		Methods []Object
	}

	Instance struct {
		Base Object
		Data map[string]Object
	}
)

/* Type() methods */
func (_ *Tuple) Type() Type    { return TUPLE }
func (_ *Array) Type() Type    { return ARRAY }
func (_ *Map) Type() Type      { return MAP }
func (_ *Block) Type() Type    { return BLOCK }
func (_ *Class) Type() Type    { return CLASS }
func (_ *Instance) Type() Type { return INSTANCE }

/* Equals() methods */
func (t *Tuple) Equals(o Object) bool {
	if other, ok := o.(*Tuple); ok {
		if len(other.Value) != len(t.Value) {
			return false
		}

		for i, elem := range t.Value {
			if !elem.Equals(other.Value[i]) {
				return false
			}
		}

		return true
	}

	return false
}

func (a *Array) Equals(o Object) bool {
	if other, ok := o.(*Array); ok {
		if len(other.Value) != len(a.Value) {
			return false
		}

		for i, elem := range a.Value {
			if !elem.Equals(other.Value[i]) {
				return false
			}
		}

		return true
	}

	return false
}

func (m *Map) Equals(o Object) bool {
	if other, ok := o.(*Map); ok {
		if len(other.Pairs) != len(m.Pairs) {
			return false
		}

		for k, v := range m.Pairs {
			if !v.Equals(other.Pairs[k]) {
				return false
			}
		}

		return true
	}

	return false
}

func (_ *Block) Equals(o Object) bool {
	_, ok := o.(*Block)
	return ok
}

func (c *Class) Equals(o Object) bool {
	if other, ok := o.(*Class); ok {
		return other.Name == c.Name
	}

	return false
}

func (i *Instance) Equals(o Object) bool {
	if other, ok := o.(*Instance); ok {
		if !other.Base.Equals(i.Base) {
			return false
		}

		for k, v := range i.Data {
			if !v.Equals(other.Data[k]) {
				return false
			}
		}

		return true
	}

	return false
}

/* Stringer implementations */
func join(arr []Object, sep string) string {
	stringArr := make([]string, len(arr))

	for i, elem := range arr {
		stringArr[i] = elem.String()
	}

	return strings.Join(stringArr, ", ")
}

func (t *Tuple) String() string {
	return fmt.Sprintf("(%s)", join(t.Value, ", "))
}

func (a *Array) String() string {
	return fmt.Sprintf("[%s]", join(a.Value, ", "))
}

func (m *Map) String() string {
	stringArr := make([]string, len(m.Pairs))
	i := 0

	for k, v := range m.Pairs {
		stringArr[i] = fmt.Sprintf(
			"%s: %s",
			k.String(),
			v.String(),
		)

		i += 1
	}

	return fmt.Sprintf("[%s]", strings.Join(stringArr, ", "))
}

func (b *Block) String() string {
	return "<block>"
}

func (c *Class) String() string {
	return c.Name
}

func (i *Instance) String() string {
	return fmt.Sprintf("<instance of %s>", i.Base.(*Class).Name)
}

/* Collection implementations */
func (t *Tuple) Elements() []Object {
	return t.Value
}

func (a *Array) Elements() []Object {
	return a.Value
}

/* Container implementations */
func (m *Map) Get(key Object) Object {
	return m.Pairs[key]
}

func (m *Map) Set(key, value Object) {
	m.Pairs[key] = value
}

/* Other methods */
func (c *Class) GetMethods() []Method {
	var methods []Method

	if c.Parent != nil {
		methods = c.Parent.(*Class).GetMethods()
	}

	for _, m := range c.Methods {
		if method, ok := m.(*Method); ok {
			methods = append(methods, *method)
		}
	}

	return methods
}

func (c *Class) GetMethod(pattern string) *Method {
	fnPattern := strings.Split(pattern, " ")

	for _, method := range c.GetMethods() {
		methodPattern := method.Fn.Pattern

		if len(fnPattern) != len(methodPattern) {
			continue
		}

		isMatch := true
		for i, mPatItem := range methodPattern {
			_, isParam := mPatItem.(*ast.Parameter)
			_, isIdent := mPatItem.(*ast.Identifier)

			if !(fnPattern[i] == "$" && isParam || isIdent && fnPattern[i] == methodPattern[i].Token().Literal) {
				isMatch = false
			}
		}

		if isMatch {
			return &method
		}
	}

	return nil
}
