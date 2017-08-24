package object

import "fmt"

type Type string

const (
	/* Internal Types */
	RETURN_VALUE Type = "<return value>"
	FUNCTION          = "<function>"
	NEXT              = "<next>"
	BREAK             = "<break>"

	/* Normal Types */
	NUMBER   Type = "<number>"
	BOOLEAN       = "<boolean>"
	STRING        = "<string>"
	CHAR          = "<char>"
	ARRAY         = "<array>"
	NULL          = "<null>"
	BLOCK         = "<block>"
	TUPLE         = "<tuple>"
	MAP           = "<map>"
	CLASS         = "<class>"
	INIT          = "<init method>"
	METHOD        = "<method>"
	INSTANCE      = "<instance>"
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
