package object

func unwrapReturnValue(o Object) Object {
	if ret, ok := o.(*ReturnValue); ok {
		return ret.Value
	}

	return o
}

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
