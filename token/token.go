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
	"def":     DEF,
	"return":  RETURN,
	"true":    TRUE,
	"yes":     TRUE,
	"false":   FALSE,
	"no":      FALSE,
	"null":    NULL,
	"if":      IF,
	"else":    ELSE,
	"elif":    ELIF,
	"while":   WHILE,
	"for":     FOR,
	"next":    NEXT,
	"break":   BREAK,
	"class":   CLASS,
	"extends": EXTENDS,
	"init":    INIT,
	"match":   MATCH,
	"try":     TRY,
	"catch":   CATCH,
}

func IsKeyword(t Type) bool {
	for _, k := range Keywords {
		if t == k {
			return true
		}
	}

	return false
}
