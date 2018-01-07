package token

import "fmt"

// Token is a lexical token used in parsing
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

// Keywords maps keyword literals to their
// corresponding token types
var Keywords = map[string]Type{
	"def":    Def,
	"return": Return,
	"true":   True,
	"yes":    True,
	"false":  False,
	"no":     False,
	"null":   Null,
	"if":     If,
	"else":   Else,
	"elif":   Elif,
	"while":  While,
	"for":    For,
	"next":   Next,
	"break":  Break,
	"use":    Use,
}

// IsKeyword checks if a token of type t is a
// keyword
func IsKeyword(t Type) bool {
	for _, k := range Keywords {
		if t == k {
			return true
		}
	}

	return false
}
