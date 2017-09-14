package bytecode

// 0-9: stack operations
const (
	Pop byte = 0
	Dup byte = 1
	Rot byte = 2
)

// 10-19: load/store
const (
	LoadConst  byte = 10
	LoadName   byte = 11
	StoreName  byte = 12
	LoadField  byte = 13
	StoreField byte = 14
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
	BinaryNotEqual byte = 37
	BinaryLessThan byte = 38
	BinaryMoreThan byte = 39
	BinaryLessEq   byte = 40
	BinaryMoreEq   byte = 41
)

// 50-59: using functions
const (
	Call   byte = 50
	Return byte = 51
)

// 60-89: pseudo-syscalls (i.e. builtin functions?)
const (
	Print   byte = 60
	Println byte = 61
)

// 90-99: control flow
const (
	Jump        byte = 90
	JumpIfTrue  byte = 91
	JumpIfFalse byte = 92
	Break       byte = 93
	Next        byte = 94
	LoopStart   byte = 95
	LoopEnd     byte = 96
)

// 100-109: data constructors
const (
	MakeArray byte = 100
	MakeTuple byte = 101
	MakeMap   byte = 102
)
