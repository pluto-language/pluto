package evaluation

import (
	"github.com/Zac-Garby/pluto/ast"
)

var (
	// NextObj is a predefined instance of Next
	NextObj = new(Next)

	// BreakObj is a predefined instance of Break
	BreakObj = new(Break)

	// NullObj is a predefined instance of Null
	NullObj = new(Null)

	// TrueObj is a predefined instance of Boolean, whose value is true
	TrueObj = &Boolean{Value: true}

	// FalseObj is a predefined instance of Boolean, whose value is false
	FalseObj = &Boolean{Value: false}
)

var (
	infixOverloads = map[string]string{
		"+":  "__plus $",
		"-":  "__minus $",
		"*":  "__times $",
		"/":  "__divide $",
		"**": "__exp $",
		"//": "__f_div $",
		`%`:  "__mod $",
		"==": "__eq $",
		"||": "__or $",
		"&&": "__and $",
		"|":  "__b_or $",
		"&":  "__b_and $",
		".":  "__get $",
	}

	prefixOverloads = map[string]string{
		"+": "__no_op",
		"-": "__negate",
		"!": "__invert",
	}
)

// EvaluateProgram evaluates a program in the given context
func EvaluateProgram(prog ast.Program) Object {
	return NullObj
}
