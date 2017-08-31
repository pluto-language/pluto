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
	token.ASSIGN:   assign,
	token.DECLARE:  assign,
	token.A_AND:    assign,
	token.A_B_AND:  assign,
	token.A_B_OR:   assign,
	token.A_EXP:    assign,
	token.A_F_DIV:  assign,
	token.A_MINUS:  assign,
	token.A_MOD:    assign,
	token.A_OR:     assign,
	token.A_PLUS:   assign,
	token.A_Q_MARK: assign,
	token.A_SLASH:  assign,
	token.A_STAR:   assign,
	token.Q_MARK:   question,
	token.OR:       or,
	token.AND:      and,
	token.B_OR:     bitOr,
	token.B_AND:    bitAnd,
	token.EQ:       equals,
	token.N_EQ:     equals,
	token.LT:       compare,
	token.GT:       compare,
	token.LTE:      compare,
	token.GTE:      compare,
	token.PLUS:     sum,
	token.MINUS:    sum,
	token.STAR:     product,
	token.SLASH:    product,
	token.MOD:      product,
	token.EXP:      exp,
	token.F_DIV:    exp,
	token.BANG:     prefix,
	token.COLON:    methodCall,
	token.D_COLON:  methodCall,
	token.DOT:      index,
	token.LSQUARE:  index,
}

var argBlacklist = []token.Type{
	token.IF,
	token.BSLASH,
	token.WHILE,
	token.FOR,
	token.MATCH,
	token.MINUS,
	token.PLUS,
	token.LSQUARE,
}

func isBlacklisted(t token.Type) bool {
	for _, blacklisted := range argBlacklist {
		if t == blacklisted {
			return true
		}
	}

	return false
}
