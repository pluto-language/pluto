package token

import "fmt"

type Token struct {
	Type       Type
	Literal    string
	Start, End Position
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"%s `%s` from %s → %s",
		t.Type,
		t.Literal,
		t.Start.String(),
		t.End.String(),
	)
}

var Keywords = map[string]Type{
	"def":     Def,
	"return":  Return,
	"true":    True,
	"yes":     True,
	"false":   False,
	"no":      False,
	"null":    Null,
	"if":      If,
	"else":    Else,
	"elif":    Elif,
	"while":   While,
	"for":     For,
	"next":    Next,
	"break":   Break,
	"class":   Class,
	"extends": Extends,
	"init":    Init,
	"match":   Match,
	"try":     Try,
	"catch":   Catch,
	"import":  Import,
	"use":     Use,
}

func IsKeyword(t Type) bool {
	for _, k := range Keywords {
		if t == k {
			return true
		}
	}

	return false
}
