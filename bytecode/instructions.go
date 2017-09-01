package bytecode

// 0-9: stack operations
const (
	Pop = 0
	Dup = 1
	Rot = 2
	//
	//
	//
	// reserved
	//
	//
	//
)

// 10-19: load/store
const (
	LoadConst = 10
	LoadName  = 11
	StoreName = 12
	//
	//
	//
	// reserved
	//
	//
	//
)

// 20-39: operators
const (
	UnaryInvert = 20
	UnaryNegate = 21
	UnaryNoOp   = 22
	// reserved
	//
	BinaryAdd      = 25
	BinarySubtract = 26
	BinaryMultiply = 27
	BinaryDivide   = 28
	BinaryExponent = 29
	BinaryFloorDiv = 30
	BinaryMod      = 31
	BinaryOr       = 32
	BinaryAnd      = 33
	BinaryBitOr    = 34
	BinaryBitAnd   = 35
	BinaryEqual    = 36
	//
	// reserved
	//
)
