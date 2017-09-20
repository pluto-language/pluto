package object

import (
	"fmt"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/fatih/color"
)

type (
	/* Collections and collection-likes */

	// Tuple is a tuple value, such as (1, "foo", false)
	Tuple struct {
		Value []Object
	}

	// Array is an array value, such as [1, 5, "baz", true]
	Array struct {
		Value []Object
	}

	// Map is a mapping of keys to values, such as ["x": 2, "y": 5]
	Map struct {
		Values map[string]Object
		Keys   map[string]Object
	}

	/* Others */

	// Block is an anonymous function
	Block struct {
		Params    []ast.Expression
		Body      bytecode.Code
		Constants []Object
		Names     []string
		Patterns  []string
	}

	// Class is a class, containing methods and a parent
	Class struct {
		Name    string
		Parent  Object
		Methods []Object
	}

	// Instance is an instance of a class
	Instance struct {
		Base Object
		Data map[string]Object
	}
)

/* Type() methods */

// Type returns the type of this object
func (t *Tuple) Type() Type { return TupleType }

// Type returns the type of this object
func (a *Array) Type() Type { return ArrayType }

// Type returns the type of this object
func (m *Map) Type() Type { return MapType }

// Type returns the type of this object
func (b *Block) Type() Type { return BlockType }

// Type returns the type of this object
func (c *Class) Type() Type { return ClassType }

// Type returns the type of this object
func (i *Instance) Type() Type { return InstanceType }

/* Equals() methods */

// Equals checks if two objects are equal to each other
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

// Equals checks if two objects are equal to each other
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

// Equals checks if two objects are equal to each other
func (m *Map) Equals(o Object) bool {
	if other, ok := o.(*Map); ok {
		if len(other.Values) != len(m.Values) {
			return false
		}

		for k, v := range m.Values {
			if _, ok := other.Values[k]; !ok {
				return false
			}

			if !v.Equals(other.Values[k]) {
				return false
			}
		}

		return true
	}

	return false
}

// Equals checks if two objects are equal to each other
func (b *Block) Equals(o Object) bool {
	_, ok := o.(*Block)
	return ok
}

// Equals checks if two objects are equal to each other
func (c *Class) Equals(o Object) bool {
	if other, ok := o.(*Class); ok {
		return other.Name == c.Name
	}

	return false
}

// Equals checks if two objects are equal to each other
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
	stringArr := make([]string, len(m.Values))
	i := 0

	for k, v := range m.Values {
		stringArr[i] = fmt.Sprintf(
			"%s: %s",
			m.Keys[k].String(),
			v.String(),
		)

		i++
	}

	return fmt.Sprintf("[%s]", strings.Join(stringArr, ", "))
}

func (b *Block) String() string {
	if len(b.Params) == 0 {
		return "<block>"
	}

	var params []string

	for _, param := range b.Params {
		params = append(params, param.Token().Literal)
	}

	return fmt.Sprintf("<block: %s>", strings.Join(params, ", "))
}

func (c *Class) String() string {
	return c.Name
}

func (i *Instance) String() string {
	if i.Base.(*Class).Name == "Error" {
		tag := color.New(color.FgRed, color.Bold).Sprint(i.Get(&String{"tag"}))
		return fmt.Sprintf("%s\t%s", tag, i.Get(&String{"msg"}))
	}

	stringMethod := i.Base.(*Class).GetMethod("string")

	if stringMethod != nil {
		panic("Implement me :)")
	}

	return fmt.Sprintf("<instance of %s>", i.Base.(*Class).Name)
}

/* Collection implementations */

// Elements returns the elements in a collection
func (t *Tuple) Elements() []Object {
	return t.Value
}

// GetIndex returns the ith element in a collection
func (t *Tuple) GetIndex(i int) Object {
	if i >= len(t.Value) || i < 0 {
		return NullObj
	}

	return t.Value[i]
}

// SetIndex sets the ith element in a collection to o
func (t *Tuple) SetIndex(i int, o Object) {
	if i >= len(t.Value) || i < 0 {
		return
	}

	t.Value[i] = o
}

// Elements returns the elements in a collection
func (a *Array) Elements() []Object {
	return a.Value
}

// GetIndex returns the ith element in a collection
func (a *Array) GetIndex(i int) Object {
	if i >= len(a.Value) || i < 0 {
		return NullObj
	}

	return a.Value[i]
}

// SetIndex sets the ith element in a collection to o
func (a *Array) SetIndex(i int, o Object) {
	if i >= len(a.Value) || i < 0 {
		return
	}

	a.Value[i] = o
}

/* Container implementations */

// Get gets an object at the given key
func (m *Map) Get(key Object) Object {
	hasher, ok := key.(Hasher)

	if !ok {
		return NullObj
	}

	if val, ok := m.Values[hasher.Hash()]; ok {
		return val
	}

	return NullObj
}

// Set sets an object at the given key
func (m *Map) Set(key, value Object) {
	if hasher, ok := key.(Hasher); ok {
		hash := hasher.Hash()
		m.Values[hash] = value
		m.Keys[hash] = key
	}
}

// Get gets an object at the given key
func (i *Instance) Get(key Object) Object {
	strkey, ok := key.(*String)

	if !ok {
		return NullObj
	}

	if val, ok := i.Data[strkey.Value]; ok {
		return val
	}

	return NullObj
}

// Set sets an object at the given key
func (i *Instance) Set(key, value Object) {
	if strkey, ok := key.(*String); ok {
		i.Data[strkey.Value] = value
	}
}

/* Other methods */

// GetMethods returns a slice of methods of a class
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

// GetMethod returns a method matching the given pattern
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
