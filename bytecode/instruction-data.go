package bytecode

type data struct {
	Name   string
	HasArg bool
}

// Instructions stores data about the different instruction types
var Instructions = map[byte]data{
	Pop: {Name: "POP"},
	Dup: {Name: "DUP"},
	Rot: {Name: "ROT"},

	LoadConst: {Name: "LOAD_CONST", HasArg: true},
	LoadName:  {Name: "LOAD_NAME", HasArg: true},
	StoreName: {Name: "STORE_NAME", HasArg: true},

	UnaryInvert: {Name: "UNARY_INVERT"},
	UnaryNegate: {Name: "UNARY_NEGATE"},
	UnaryNoOp:   {Name: "UNARY_NO_OP"},

	BinaryAdd:      {Name: "BINARY_ADD"},
	BinarySubtract: {Name: "BINARY_SUBTRACT"},
	BinaryMultiply: {Name: "BINARY_MULTIPLY"},
	BinaryDivide:   {Name: "BINARY_DIVIDE"},
	BinaryExponent: {Name: "BINARY_EXPONENT"},
	BinaryFloorDiv: {Name: "BINARY_FLOOR_DIV"},
	BinaryMod:      {Name: "BINARY_MOD"},
	BinaryOr:       {Name: "BINARY_OR"},
	BinaryAnd:      {Name: "BINARY_AND"},
	BinaryBitOr:    {Name: "BINARY_BIT_OR"},
	BinaryBitAnd:   {Name: "BINARY_BIT_AND"},
	BinaryEquals:   {Name: "BINARY_EQUALS"},
	Print:          {Name: "PRINT"},

	Call:   {Name: "CALL_FN", HasArg: true},
	Return: {Name: "RETURN_FN"},

	Jump:        {Name: "JUMP", HasArg: true},
	JumpIfTrue:  {Name: "JUMP_IF_TRUE", HasArg: true},
	JumpIfFalse: {Name: "JUMP_IF_FALSE", HasArg: true},

	MakeArray: {Name: "MAKE_ARRAY", HasArg: true},
	MakeTuple: {Name: "MAKE_TUPLE", HasArg: true},
	MakeMap:   {Name: "MAKE_MAP", HasArg: true},
}
