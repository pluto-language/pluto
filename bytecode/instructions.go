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
	/*  0x14  */ UnaryInvert byte = 20
	/*  0x15  */ UnaryNegate byte = 21
	/*  0x16  */ UnaryNoOp byte = 22
	/*  0x19  */ BinaryAdd byte = 25
	/*  0x1a  */ BinarySubtract byte = 26
	/*  0x1b  */ BinaryMultiply byte = 27
	/*  0x1c  */ BinaryDivide byte = 28
	/*  0x1d  */ BinaryExponent byte = 29
	/*  0x1e  */ BinaryFloorDiv byte = 30
	/*  0x1f  */ BinaryMod byte = 31
	/*  0x20  */ BinaryOr byte = 32
	/*  0x21  */ BinaryAnd byte = 33
	/*  0x22  */ BinaryBitOr byte = 34
	/*  0x23  */ BinaryBitAnd byte = 35
	/*  0x24  */ BinaryEquals byte = 36
	/*  0x25  */ Print byte = 37
)

// 40-49: using functions
const (
	/*  0x28  */ Call byte = 40
	/*  0x29  */ Return byte = 41
)
