package object

import "fmt"

type Type string

const (
	/* Internal Types */
	RETURN_VALUE Type = "<return value>"
	FUNCTION     Type = "<function>"
	NEXT         Type = "<next>"
	BREAK        Type = "<break>"

	/* Normal Types */
	NUMBER   Type = "<number>"
	BOOLEAN  Type = "<boolean>"
	STRING   Type = "<string>"
	CHAR     Type = "<char>"
	ARRAY    Type = "<array>"
	NULL     Type = "<null>"
	BLOCK    Type = "<block>"
	TUPLE    Type = "<tuple>"
	MAP      Type = "<map>"
	CLASS    Type = "<class>"
	INIT     Type = "<init method>"
	METHOD   Type = "<method>"
	INSTANCE Type = "<instance>"
)

type Object interface {
	fmt.Stringer
	Equals(Object) bool
	Type() Type
}

type Collection interface {
	Object
	Elements() []Object
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
