package object

// IsTruthy returns true if o is truthy,
// and false otherwise.
func IsTruthy(o Object) bool {
	if o.Equals(NullObj) || o.Equals(FalseObj) {
		return false
	}

	if num, ok := o.(*Number); ok {
		return num.Value != 0
	}

	if col, ok := o.(Collection); ok {
		return len(col.Elements()) != 0
	}

	return true
}

// BoolObj converts a native bool
// value to a Pluto boolean.
func BoolObj(t bool) Object {
	if t {
		return TrueObj
	}

	return FalseObj
}

// MakeCollection creates a collection of
// type t containing the given elements.
func MakeCollection(t Type, elements []Object) (Object, bool) {
	switch t {
	case ArrayType:
		return &Array{Value: elements}, true
	case TupleType:
		return &Tuple{Value: elements}, true
	case StringType:
		var str string

		for _, elem := range elements {
			str += elem.String()
		}

		return &String{Value: str}, true
	default:
		return nil, false
	}
}
