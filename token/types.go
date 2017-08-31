package token

// Type is the type of a token
type Type string

const (
	// EOF is at the end of every file
	EOF = "EOF"

	// Illegal is any non-recognized character
	Illegal = "illegal"

	// Number is a number literal (123.456)
	Number = "number"

	// String is a string literal ("foo")
	String = "string"

	// Char is a character literal ('x')
	Char = "char"

	// ID is an identifier (bar)
	ID = "identifier"

	// Param is a parameter ($baz)
	Param = "parameter"

	// Plus is the addition operator (+)
	Plus = "plus"

	// Minus is the subtraction operator (-)
	Minus = "minus"

	// Star is the multiplication operator (*)
	Star = "star"

	// Exp is the exponentiation operator (**)
	Exp = "exponent"

	// Slash is the division operator (/)
	Slash = "slash"

	// FloorDiv is the floor division operator (//)
	FloorDiv = "floor-div"

	// Mod is the modulo operator (%)
	Mod = "modulo"

	// BackSlash is a backslash (\)
	BackSlash = "backslash"

	// LeftParen is a left paren '('
	LeftParen = "lparen"

	// RightParen is a right paren ')'
	RightParen = "rparen"

	// LessThan is the less than operator (<)
	LessThan = "less-than"

	// GreaterThan is the greater than operator (>)
	GreaterThan = "greater-than"

	// LessThanEq is the less than or equal to operator (<=)
	LessThanEq = "less-than-or-equal"

	// GreaterThanEq is the greater than or equal to operator (>=)
	GreaterThanEq = "greater-than-or-equal"

	// LeftBrace is a left brace ({)
	LeftBrace = "left-brace"

	// RightBrace is a right brace (})
	RightBrace = "right-brace"

	// LeftSquare is a left square bracket ([)
	LeftSquare = "left-square"

	// RightSquare is a right square bracket (])
	RightSquare = "right-square"

	// Semi is a semi-colon, either in the source or ASI'd
	Semi = "semi"

	// Equal is the equality test operator (==)
	Equal = "equal"

	// NotEqual is the inverted equality test operator (!=)
	NotEqual = "not-equal"

	// Or is the or operator (||)
	Or = "or"

	// And is the and operator (&&)
	And = "and"

	// BitOr is the bitwise or operator (|)
	BitOr = "bitwise-or"

	// BitAnd is the bitwise and operator (&)
	BitAnd = "bitwise-and"

	// Assign is the assign operator (=)
	Assign = "assign"

	// Declare is the declare operator (:=)
	Declare = "declare"

	// Comma is a comma (,)
	Comma = "comma"

	// Arrow is a right-facing arrow (->)
	Arrow = "arrow"

	// Colon is a colon (:)
	Colon = "colon"

	// QuestionMark is the question-mark operator (?)
	QuestionMark = "question-mark"

	// Dot is the dot-access operator (.)
	Dot = "dot"

	// FatArrow is the fat arrow (=>)
	FatArrow = "fat-arrow"

	// Bang is an exclaimation mark (!)
	Bang = "bang"

	// DoubleColon is two colons (::)
	DoubleColon = "double-colon"

	// PlusEquals is the addition-assignment operator (+=)
	PlusEquals = "assign-plus"

	// MinusEquals is the subtraction-assignment operator (-=)
	MinusEquals = "assign-minus"

	// StarEquals is the multiplication-assignment operator (*=)
	StarEquals = "assign-star"

	// ExpEquals is the exponentiation-assignment operator (**=)
	ExpEquals = "assign-exponent"

	// SlashEquals is the division-assignment operator (/=)
	SlashEquals = "assign-slash"

	// FloorDivEquals is the floor-division-assignment operator (//=)
	FloorDivEquals = "assign-floor-div"

	// ModEquals is the modulo-assignment operator (%=)
	ModEquals = "assign-modulo"

	// OrEquals is the or-assignment operator (||=)
	OrEquals = "assign-or"

	// AndEquals is the and-assignment operator (&&=)
	AndEquals = "assign-and"

	// BitOrEquals is the bitwise-or-assignment operator (|=)
	BitOrEquals = "assign-bitwise-or"

	// BitAndEquals is the bitwise-and-assignment operator (&=)
	BitAndEquals = "assign-bitwise-and"

	// QuestionMarkEquals is the question-mark-assignment operator (?=)
	QuestionMarkEquals = "assign-question-mark"

	// Def is the 'def' keyword
	Def = "def"

	// Return is the 'return' keyword
	Return = "return"

	// True is the 'true' keyword
	True = "true"

	// False is the 'false' keyword
	False = "false"

	// Null is the 'null' keyword
	Null = "null"

	// If is the 'if' keyword
	If = "if"

	// Else is the 'else' keyword
	Else = "else"

	// Elif is the 'elif' keyword
	Elif = "elif"

	// While is the 'while' keyword
	While = "while"

	// For is the 'for' keyword
	For = "for"

	// Next is the 'next' keyword
	Next = "next"

	// Break is the 'break' keyword
	Break = "break"

	// Class is the 'class' keyword
	Class = "class"

	// Extends is the 'extends' keyword
	Extends = "extends"

	// Init is the 'init' keyword
	Init = "init"

	// Match is the 'match' keyword
	Match = "match"

	// Try is the 'try' keyword
	Try = "try"

	// Catch is the 'catch' keyword
	Catch = "catch"

	// Import is the 'import' keyword
	Import = "import"

	// Use is the 'use' keyword
	Use = "use"
)
