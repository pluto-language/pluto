package parser

import "github.com/Zac-Garby/pluto/token"

const (
	lowest = iota
	assign
	question
	or
	and
	bitOr
	bitAnd
	equals
	compare
	sum
	product
	exp
	prefix
	methodCall
	index
)

var precedences = map[token.Type]int{
	token.Assign:             assign,
	token.Declare:            assign,
	token.AndEquals:          assign,
	token.BitAndEquals:       assign,
	token.BitOrEquals:        assign,
	token.ExpEquals:          assign,
	token.FloorDivEquals:     assign,
	token.MinusEquals:        assign,
	token.ModEquals:          assign,
	token.OrEquals:           assign,
	token.PlusEquals:         assign,
	token.QuestionMarkEquals: assign,
	token.SlashEquals:        assign,
	token.StarEquals:         assign,
	token.QuestionMark:       question,
	token.Or:                 or,
	token.And:                and,
	token.BitOr:              bitOr,
	token.BitAnd:             bitAnd,
	token.Equal:              equals,
	token.NotEqual:           equals,
	token.LessThan:           compare,
	token.GreaterThan:        compare,
	token.LessThanEq:         compare,
	token.GreaterThanEq:      compare,
	token.Plus:               sum,
	token.Minus:              sum,
	token.Star:               product,
	token.Slash:              product,
	token.Mod:                product,
	token.Exp:                exp,
	token.FloorDiv:           exp,
	token.Bang:               prefix,
	token.Colon:              methodCall,
	token.DoubleColon:        methodCall,
	token.Dot:                index,
	token.LeftSquare:         index,
}

var argBlacklist = []token.Type{
	token.If,
	token.BackSlash,
	token.While,
	token.For,
	token.Match,
	token.Minus,
	token.Plus,
	token.LeftSquare,
}

func isBlacklisted(t token.Type) bool {
	for _, blacklisted := range argBlacklist {
		if t == blacklisted {
			return true
		}
	}

	return false
}
