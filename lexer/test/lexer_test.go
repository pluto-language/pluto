package lexer

import (
	"testing"

	"github.com/Zac-Garby/pluto/token"

	. "github.com/Zac-Garby/pluto/lexer"
)

func TestLexer(t *testing.T) {
	cases := map[string][]token.Type{
		`@ " ' $`:                           {token.Illegal, token.Illegal, token.Illegal, token.Illegal},
		`1 2 3`:                             {token.Number, token.Number, token.Number},
		`"Hello" "foo bar" "x\"y"`:          {token.String, token.String, token.String},
		`'a' ' '`:                           {token.Char, token.Char},
		`hello_world fooBar123`:             {token.ID, token.ID},
		`$param $under_score $123`:          {token.Param, token.Param, token.Param},
		`+ - * ** / // %`:                   {token.Plus, token.Minus, token.Star, token.Exp, token.Slash, token.FloorDiv, token.Mod},
		`\ ( ) { } [ ]`:                     {token.BackSlash, token.LeftParen, token.RightParen, token.LeftBrace, token.RightBrace, token.LeftSquare, token.RightSquare},
		`< > <= >= == !=`:                   {token.LessThan, token.GreaterThan, token.LessThanEq, token.GreaterThanEq, token.Equal, token.NotEqual},
		`|| && | &`:                         {token.Or, token.And, token.BitOr, token.BitAnd},
		`= += -= *= **= /= //= %=`:          {token.Assign, token.PlusEquals, token.MinusEquals, token.StarEquals, token.ExpEquals, token.SlashEquals, token.FloorDivEquals, token.ModEquals},
		`||= &&= |= &= ?=`:                  {token.OrEquals, token.AndEquals, token.BitOrEquals, token.BitAndEquals, token.QuestionMarkEquals},
		`true false null`:                   {token.True, token.False, token.Null},
		`def return import use`:             {token.Def, token.Return, token.Import, token.Use},
		`if else elif while for next break`: {token.If, token.Else, token.Elif, token.While, token.For, token.Next, token.Break},
	}

	for in, out := range cases {
		var (
			out = append(out, token.Semi)
			l   = Lexer(in, "<test suite>")
			tok token.Token
		)

		for _, exp := range out {
			tok = l()

			if tok.Type == token.EOF {
				break
			}

			if tok.Type != exp {
				t.Errorf("wrong token type (%s). expected %s, got %s", tok.Literal, exp, tok.Type)
			}
		}

		if tt := l().Type; tt != token.EOF {
			t.Errorf("too many tokens scanned (wanted %d, last: %s)", len(out), tt)
		}
	}
}
