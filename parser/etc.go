package parser

import "github.com/Zac-Garby/pluto/token"

const (
	LOWEST = iota
	ASSIGN
	QUESTION
	OR
	AND
	BIT_OR
	BIT_AND
	EQUALS
	COMPARE
	SUM
	PRODUCT
	EXP
	PREFIX
	METHOD_CALL
	INDEX
)

var precedences = map[token.Type]int{
	token.ASSIGN:   ASSIGN,
	token.DECLARE:  ASSIGN,
	token.A_AND:    ASSIGN,
	token.A_B_AND:  ASSIGN,
	token.A_B_OR:   ASSIGN,
	token.A_EXP:    ASSIGN,
	token.A_F_DIV:  ASSIGN,
	token.A_MINUS:  ASSIGN,
	token.A_MOD:    ASSIGN,
	token.A_OR:     ASSIGN,
	token.A_PLUS:   ASSIGN,
	token.A_Q_MARK: ASSIGN,
	token.A_SLASH:  ASSIGN,
	token.A_STAR:   ASSIGN,
	token.Q_MARK:   QUESTION,
	token.OR:       OR,
	token.AND:      AND,
	token.B_OR:     BIT_OR,
	token.B_AND:    BIT_AND,
	token.EQ:       EQUALS,
	token.N_EQ:     EQUALS,
	token.LT:       COMPARE,
	token.GT:       COMPARE,
	token.LTE:      COMPARE,
	token.GTE:      COMPARE,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.STAR:     PRODUCT,
	token.SLASH:    PRODUCT,
	token.MOD:      PRODUCT,
	token.EXP:      EXP,
	token.F_DIV:    EXP,
	token.BANG:     PREFIX,
	token.COLON:    METHOD_CALL,
	token.DOT:      INDEX,
	token.LSQUARE:  INDEX,
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
