package lexer

import (
	"strings"

	"github.com/Zac-Garby/pluto/token"
)

type transformer func(token.Type, string, string) (token.Type, string, string)
type handler func([]string) (token.Type, string, string)

func lexemeHandler(t token.Type, group int, transformer transformer) handler {
	return func(m []string) (token.Type, string, string) {
		return transformer(t, m[group], m[0])
	}
}

func none(t token.Type, literal, whole string) (token.Type, string, string) {
	return t, literal, whole
}

func stringTransformer(t token.Type, literal, whole string) (token.Type, string, string) {
	escapes := map[string]string{
		`\n`: "\n",
		`"`:  "\"",
		`\a`: "\a",
		`\b`: "\b",
		`\f`: "\f",
		`\r`: "\r",
		`\t`: "\t",
		`\v`: "\v",
	}

	for k, v := range escapes {
		literal = strings.Replace(literal, k, v, -1)
	}

	return t, literal, whole
}

func idTransformer(t token.Type, literal, whole string) (token.Type, string, string) {
	if t, ok := token.Keywords[literal]; ok {
		return t, literal, whole
	}

	return t, literal, whole
}

type lexicalPair struct {
	regex   string
	handler handler
}

var lexicalDictionary = []lexicalPair{
	// Literals
	{regex: `^\d+(?:\.\d+)?`, handler: lexemeHandler(token.NUM, 0, none)},
	{regex: `^"((\\"|[^"])*)"`, handler: lexemeHandler(token.STR, 1, stringTransformer)},
	{regex: "^`([^`]*)`", handler: lexemeHandler(token.STR, 1, none)},
	{regex: `^'([^']|\w)'`, handler: lexemeHandler(token.CHAR, 1, stringTransformer)},
	{regex: `^\w+`, handler: lexemeHandler(token.ID, 0, idTransformer)},
	{regex: `^\$(\w+)`, handler: lexemeHandler(token.PARAM, 1, none)},

	// Punctuation
	{regex: `^->`, handler: lexemeHandler(token.ARROW, 0, none)},
	{regex: `^\+`, handler: lexemeHandler(token.PLUS, 0, none)},
	{regex: `^-`, handler: lexemeHandler(token.MINUS, 0, none)},
	{regex: `^\*\*`, handler: lexemeHandler(token.EXP, 0, none)},
	{regex: `^\*`, handler: lexemeHandler(token.STAR, 0, none)},
	{regex: `^\/\/`, handler: lexemeHandler(token.F_DIV, 0, none)},
	{regex: `^\/`, handler: lexemeHandler(token.SLASH, 0, none)},
	{regex: `^\\`, handler: lexemeHandler(token.BSLASH, 0, none)},
	{regex: `^\(`, handler: lexemeHandler(token.LPAREN, 0, none)},
	{regex: `^\)`, handler: lexemeHandler(token.RPAREN, 0, none)},
	{regex: `^<=`, handler: lexemeHandler(token.LTE, 0, none)},
	{regex: `^>=`, handler: lexemeHandler(token.GTE, 0, none)},
	{regex: `^<`, handler: lexemeHandler(token.LT, 0, none)},
	{regex: `^>`, handler: lexemeHandler(token.GT, 0, none)},
	{regex: `^{`, handler: lexemeHandler(token.LBRACE, 0, none)},
	{regex: `^}`, handler: lexemeHandler(token.RBRACE, 0, none)},
	{regex: `^\[`, handler: lexemeHandler(token.LSQUARE, 0, none)},
	{regex: `^]`, handler: lexemeHandler(token.RSQUARE, 0, none)},
	{regex: `^;`, handler: lexemeHandler(token.SEMI, 0, none)},
	{regex: `^==`, handler: lexemeHandler(token.EQ, 0, none)},
	{regex: `^!=`, handler: lexemeHandler(token.N_EQ, 0, none)},
	{regex: `^\|\|`, handler: lexemeHandler(token.OR, 0, none)},
	{regex: `^&&`, handler: lexemeHandler(token.AND, 0, none)},
	{regex: `^\|`, handler: lexemeHandler(token.B_OR, 0, none)},
	{regex: `^&`, handler: lexemeHandler(token.B_AND, 0, none)},
	{regex: `^=>`, handler: lexemeHandler(token.F_ARROW, 0, none)},
	{regex: `^=`, handler: lexemeHandler(token.ASSIGN, 0, none)},
	{regex: `^:=`, handler: lexemeHandler(token.DECLARE, 0, none)},
	{regex: `^\,`, handler: lexemeHandler(token.COMMA, 0, none)},
	{regex: `^:`, handler: lexemeHandler(token.COLON, 0, none)},
	{regex: `^%`, handler: lexemeHandler(token.MOD, 0, none)},
	{regex: `^\?`, handler: lexemeHandler(token.Q_MARK, 0, none)},
	{regex: `^\.`, handler: lexemeHandler(token.DOT, 0, none)},
	{regex: `^!`, handler: lexemeHandler(token.BANG, 0, none)},
}
