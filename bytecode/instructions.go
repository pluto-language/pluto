package bytecode

// 0-9: stack operations
const (
	Pop uint16 = 0
	Dup uint16 = 1
	Rot uint16 = 2
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
	LoadConst uint16 = 10
	LoadName  uint16 = 11
	StoreName uint16 = 12
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
	UnaryInvert uint16 = 20
	UnaryNegate uint16 = 21
	UnaryNoOp   uint16 = 22
	// reserved
	//
	BinaryAdd      uint16 = 25
	BinarySubtract uint16 = 26
	BinaryMultiply uint16 = 27
	BinaryDivide   uint16 = 28
	BinaryExponent uint16 = 29
	BinaryFloorDiv uint16 = 30
	BinaryMod      uint16 = 31
	BinaryOr       uint16 = 32
	BinaryAnd      uint16 = 33
	BinaryBitOr    uint16 = 34
	BinaryBitAnd   uint16 = 35
	BinaryEquals   uint16 = 36
	//
	// reserved
	//
)
