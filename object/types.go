package object

type Type string

const (
	/* Internal Types */
	RETURN_VALUE Type = "<return value>"
	FUNCTION     Type = "<function>"
	NEXT         Type = "<next>"
	BREAK        Type = "<break>"

	/* Special Types */
	COLLECTION Type = "<collection>"
	CONTAINER  Type = "<container>"
	HASHER     Type = "<hasher>"
	ANY        Type = "<any>"

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

func is(obj Object, t Type) bool {
	if t == ANY {
		return true
	}

	if t == COLLECTION {
		_, ok := obj.(Collection)
		return ok
	}

	if t == CONTAINER {
		_, ok := obj.(Container)
		return ok
	}

	if t == HASHER {
		_, ok := obj.(Hasher)
		return ok
	}

	return obj.Type() == t
}
