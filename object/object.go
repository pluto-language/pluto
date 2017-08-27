package object

import "fmt"

type Object interface {
	fmt.Stringer
	Equals(Object) bool
	Type() Type
}

type Collection interface {
	Object
	Elements() []Object
	GetIndex(int) Object
	SetIndex(int, Object)
}

type Container interface {
	Object
	Get(Object) Object
	Set(Object, Object)
}

type Hasher interface {
	Object
	Hash() string
}
