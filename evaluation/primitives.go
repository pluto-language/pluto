package evaluation

import (
	"fmt"
)

/* Structs */
type (
	// Number is a number object
	Number struct {
		Value float64
	}

	// Boolean is a boolean object
	Boolean struct {
		Value bool
	}

	// String is a string object
	String struct {
		Value string
	}

	// Char is a character object
	Char struct {
		Value rune
	}

	// Null is the null object
	Null struct{}
)

/* Type() methods */

// Type returns the type of this object
func (n *Number) Type() Type { return NumberType }

// Type returns the type of this object
func (b *Boolean) Type() Type { return BooleanType }

// Type returns the type of this object
func (s *String) Type() Type { return StringType }

// Type returns the type of this object
func (c *Char) Type() Type { return CharType }

// Type returns the type of this object
func (n *Null) Type() Type { return NullType }

/* Equals() methods */

// Equals checks if two objects are equal to each other
func (n *Number) Equals(o Object) bool {
	if other, ok := o.(*Number); ok {
		return n.Value == other.Value
	}

	return false
}

// Equals checks if two objects are equal to each other
func (b *Boolean) Equals(o Object) bool {
	if other, ok := o.(*Boolean); ok {
		return b.Value == other.Value
	}

	return false
}

// Equals checks if two objects are equal to each other
func (s *String) Equals(o Object) bool {
	if other, ok := o.(*String); ok {
		return s.Value == other.Value
	}

	return false
}

// Equals checks if two objects are equal to each other
func (c *Char) Equals(o Object) bool {
	if other, ok := o.(*Char); ok {
		return c.Value == other.Value
	}

	return false
}

// Equals checks if two objects are equal to each other
func (n *Null) Equals(o Object) bool {
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

func (n *Null) String() string {
	return "null"
}

/* Collection implementations */

// Elements returns the elements in a collection
func (s *String) Elements() []Object {
	chars := make([]Object, len(s.Value))

	for i, ch := range s.Value {
		chars[i] = &Char{ch}
	}

	return chars
}

// GetIndex returns the ith element in a collection
func (s *String) GetIndex(i int) Object {
	if i >= len(s.Value) || i < 0 {
		return NullObj
	}

	return &Char{Value: rune(s.Value[i])}
}

// SetIndex sets the ith element in a collection to o
func (s *String) SetIndex(i int, o Object) {
	if i >= len(s.Value) || i < 0 {
		return
	}

	if ch, ok := o.(*Char); ok {
		bytes := []byte(s.Value)
		bytes[i] = byte(ch.Value)
		s.Value = string(bytes)
	}
}

/* Hasher implementations */

// Hash returns a string unique to the current state of the object
func (n *Number) Hash() string {
	return fmt.Sprintf("number %g", n.Value)
}

// Hash returns a string unique to the current state of the object
func (b *Boolean) Hash() string {
	return fmt.Sprintf("boolean %t", b.Value)
}

// Hash returns a string unique to the current state of the object
func (s *String) Hash() string {
	return fmt.Sprintf("string %s", s.Value)
}

// Hash returns a string unique to the current state of the object
func (c *Char) Hash() string {
	return fmt.Sprintf("char %s", string(c.Value))
}

// Hash returns a string unique to the current state of the object
func (n *Null) Hash() string {
	return "null"
}
