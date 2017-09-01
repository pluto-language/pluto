package evaluation

// IsErr checks if an object is an instance of Error
func IsErr(o Object) bool {
	if instance, ok := o.(*Instance); ok {
		return instance.Base.(*Class).Name == "Error"
	}

	return false
}
