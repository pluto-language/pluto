package bytecode

// 0-9: stack operations
const (
	/*  0x00  */ Pop byte = 0
	/*  0x01  */ Dup byte = 1
	/*  0x02  */ Rot byte = 2
)

// 10-19: load/store
const (
	/*  0x0a  */ LoadConst byte = 10
	/*  0x0b  */ LoadName byte = 11
	/*  0x0c  */ StoreName byte = 12
)

// 20-39: operators
const (
	UnaryInvert    byte = 20
	UnaryNegate    byte = 21
	UnaryNoOp      byte = 22
	BinaryAdd      byte = 25
	BinarySubtract byte = 26
	BinaryMultiply byte = 27
	BinaryDivide   byte = 28
	BinaryExponent byte = 29
	BinaryFloorDiv byte = 30
	BinaryMod      byte = 31
	BinaryOr       byte = 32
	BinaryAnd      byte = 33
	BinaryBitOr    byte = 34
	BinaryBitAnd   byte = 35
	BinaryEquals   byte = 36
)

// 40-49: using functions
const (
	Call   byte = 40
	Return byte = 41
)

// 50-59: pseudo-syscalls (i.e. builtin functions?)
const (
	Print   byte = 50
	Println byte = 51
)
