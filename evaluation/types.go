package evaluation

type Type = string

const (
	/* Internal Types */
	RETURN_VALUE = "<return value>"
	FUNCTION     = "<function>"
	NEXT         = "<next>"
	BREAK        = "<break>"

	/* Special Types */
	COLLECTION = "<collection>"
	CONTAINER  = "<container>"
	HASHER     = "<hasher>"
	ANY        = "<any>"

	/* Normal Types */
	NUMBER   = "<number>"
	BOOLEAN  = "<boolean>"
	STRING   = "<string>"
	CHAR     = "<char>"
	ARRAY    = "<array>"
	NULL     = "<null>"
	BLOCK    = "<block>"
	TUPLE    = "<tuple>"
	MAP      = "<map>"
	CLASS    = "<class>"
	INIT     = "<init method>"
	METHOD   = "<method>"
	INSTANCE = "<instance>"
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
