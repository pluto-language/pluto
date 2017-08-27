package object

import (
	"fmt"
)

/* Structs */
type (
	Number struct {
		Value float64
	}

	Boolean struct {
		Value bool
	}

	String struct {
		Value string
	}

	Char struct {
		Value rune
	}

	Null struct{}
)

/* Type() methods */
func (_ *Number) Type() Type  { return NUMBER }
func (_ *Boolean) Type() Type { return BOOLEAN }
func (_ *String) Type() Type  { return STRING }
func (_ *Char) Type() Type    { return CHAR }
func (_ *Null) Type() Type    { return NULL }

/* Equals() methods */
func (n *Number) Equals(o Object) bool {
	if other, ok := o.(*Number); ok {
		return n.Value == other.Value
	}

	return false
}

func (b *Boolean) Equals(o Object) bool {
	if other, ok := o.(*Boolean); ok {
		return b.Value == other.Value
	}

	return false
}

func (s *String) Equals(o Object) bool {
	if other, ok := o.(*String); ok {
		return s.Value == other.Value
	}

	return false
}

func (c *Char) Equals(o Object) bool {
	if other, ok := o.(*Char); ok {
		return c.Value == other.Value
	}

	return false
}

func (_ *Null) Equals(o Object) bool {
	_, ok := o.(*Null)
	return ok
}

/* Stringer implementations */
func (n *Number) String() string {
	return fmt.Sprintf("%g", n.Value)
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

func (s *String) String() string {
	return s.Value
}

func (c *Char) String() string {
	return string(c.Value)
}

func (_ *Null) String() string {
	return "null"
}

/* Collection implementations */
func (s *String) Elements() []Object {
	chars := make([]Object, len(s.Value))

	for i, ch := range s.Value {
		chars[i] = &Char{ch}
	}

	return chars
}

func (s *String) GetIndex(i int) Object {
	return &Char{Value: rune(s.Value[i])}
}

func (s *String) SetIndex(i int, o Object) {
	fmt.Println(o)

	if ch, ok := o.(*Char); ok {
		bytes := []byte(s.Value)
		bytes[i] = byte(ch.Value)
		s.Value = string(bytes)
	}
}

/* Hasher implementations */
func (n *Number) Hash() string {
	return fmt.Sprintf("number %g", n.Value)
}

func (b *Boolean) Hash() string {
	return fmt.Sprintf("boolean %t", b.Value)
}

func (s *String) Hash() string {
	return fmt.Sprintf("string %s", s.Value)
}

func (c *Char) Hash() string {
	return fmt.Sprintf("char %s", string(c.Value))
}

func (n *Null) Hash() string {
	return "null"
}
