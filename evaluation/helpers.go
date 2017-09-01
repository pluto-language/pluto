package evaluation

func unwrapReturnValue(o Object) Object {
	if ret, ok := o.(*ReturnValue); ok {
		return ret.Value
	}

	return o
}

func isTruthy(o Object) bool {
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

func boolObj(t bool) Object {
	if t {
		return TrueObj
	}

	return FalseObj
}
