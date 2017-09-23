package bytecode

// 0-9: stack operations
const (
	// Pop pops the stack
	Pop byte = 0

	// Dup duplicates the top item, so [x, y, z] -> [x, y, z, z]
	Dup byte = 1

	// Rot rotates the top two items, so [x, y, z] -> [x, z, y]
	Rot byte = 2
)

// 10-19: load/store
const (
	// LoadConst loads a constant by index
	LoadConst byte = 10

	// LoadName loads a name by name index
	LoadName byte = 11

	// StoreName stores the top item
	StoreName byte = 12

	// LoadField pops two items, essentially does second[top]
	LoadField byte = 13

	// StoreField pops three items, essentially does second[top] = third
	StoreField byte = 14
)

// 20-39: operators
const (
	// Unary operators pop one item and do something with it
	UnaryInvert byte = 20
	UnaryNegate byte = 21
	UnaryNoOp   byte = 22

	// Binary operators pop two items and do something with them
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

// 50-59: using functions/blocks
const (
	// PushFn pushes the function at pattern index n
	PushFn byte = 50

	// CallFn calls the function at the top of the stack,
	// popping arguments as necessary
	CallFn byte = 51

	// Return skips to the end of the context
	Return byte = 52

	// DoBlock executes the block at the top of the stack,
	// popping arguments off as necessary
	DoBlock byte = 53
)

// 60-89: pseudo-syscalls (i.e. builtin functions?)
const (
	// Print prints the item at the top of the stack
	Print byte = 60

	// Println prints the item at the top of the stack,
	// with a trailing new line
	Println byte = 61
)

// 90-99: control flow
const (
	// Jump unconditionally jumps to the given offset
	Jump byte = 90

	// JumpIfTrue jumps to the given offset if the top item is truthy
	JumpIfTrue byte = 91

	// JumpIfFalse jumps to the given offset if the top item is falsey
	JumpIfFalse byte = 92

	// Break jumps to the LoopEnd instruction of the innermost loop
	Break byte = 93

	// Next jumps to the LoopStart instruction of the innermost loop
	Next byte = 94

	// LoopStart pushes the start and end positions for the loop
	LoopStart byte = 95

	// LoopEnd pops the start and end positions
	LoopEnd byte = 96
)

// 100-109: data constructors
const (
	// MakeArray makes an array object from the top n items
	MakeArray byte = 100

	// MakeTuple makes a tuple from the top n items
	MakeTuple byte = 101

	// MakeMap makes a map from the top n * 2 items.
	// The top n*2 items should be in key, val, ..., key, val order
	MakeMap byte = 102
)

// 110-119: packages
const (
	// Use directly imports the specified sources
	Use byte = 110
)
