package evaluator

import (
	"github.com/Zac-Garby/pluto/object"
)

const (
	NEXT  = new(object.Next)
	BREAK = new(object.Break)

	NULL  = new(object.Null)
	TRUE  = &object.Boolean{true}
	FALSE = &object.Boolean{false}
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
