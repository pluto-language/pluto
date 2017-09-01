package object

var (
	// NullObj is the null value
	NullObj = &Null{}

	// TrueObj is the true boolean
	TrueObj = &Boolean{Value: true}

	// FalseObj is the false boolean
	FalseObj = &Boolean{Value: false}
)
