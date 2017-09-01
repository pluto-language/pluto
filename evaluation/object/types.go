package object

// Type is an object type
type Type string

const (
	/* Internal Types */

	// ReturnValueType is the type of a return value
	ReturnValueType = "<return value>"

	// FunctionType is the type of a function
	FunctionType = "<function>"

	// NextType is the type of the 'next' value
	NextType = "<next>"

	// BreakType is the type of the 'break' value
	BreakType = "<break>"

	/* Special Types */

	// CollectionType is the general type of collections
	CollectionType = "<collection>"

	// ContainerType is the general type of containers
	ContainerType = "<container>"

	// HasherType is the type of any value that can be hashed
	HasherType = "<hasher>"

	// AnyType is any type
	AnyType = "<any>"

	/* Normal Types */

	// NumberType is the type of number objects
	NumberType = "<number>"

	// BooleanType is the type of boolean objects
	BooleanType = "<boolean>"

	// StringType is the type of string objects
	StringType = "<string>"

	// CharType is the type of char objects
	CharType = "<char>"

	// ArrayType is the type of array objects
	ArrayType = "<array>"

	// NullType is the type of the null object
	NullType = "<null>"

	// BlockType is the type of block objects
	BlockType = "<block>"

	// TupleType is the type of tuples
	TupleType = "<tuple>"

	// MapType is the type of map objects
	MapType = "<map>"

	// ClassType is the type of class objects
	ClassType = "<class>"

	// InitType is the type of init methods
	InitType = "<init method>"

	// MethodType is the type of regular methods
	MethodType = "<method>"

	// InstanceType is the type of class instances
	InstanceType = "<instance>"
)

func is(obj Object, t Type) bool {
	if t == AnyType {
		return true
	}

	if t == CollectionType {
		_, ok := obj.(Collection)
		return ok
	}

	if t == ContainerType {
		_, ok := obj.(Container)
		return ok
	}

	if t == HasherType {
		_, ok := obj.(Hasher)
		return ok
	}

	return obj.Type() == t
}
