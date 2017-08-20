package lexer

import (
	"regexp"
	"unicode"

	"github.com/Zac-Garby/pluto/token"
)

var lineEndings = []token.Type{
	token.ID,
	token.STR,
	token.CHAR,
	token.NUM,
	token.TRUE,
	token.FALSE,
	token.NULL,
	token.PARAM,
	token.BREAK,
	token.NEXT,
	token.RETURN,
	token.RPAREN,
	token.RSQUARE,
	token.RBRACE,
}

func Lexer(str string) func() token.Token {
	var (
		index = 0
		col   = 1
		line  = 1
		ch    = make(chan token.Token)
	)

	go func() {
		for {
			if index < len(str) {
				foundSpace := false

				for index < len(str) && (unicode.IsSpace(rune(str[index])) || str[index] == '#') {
					if unicode.IsSpace(rune(str[index])) {
						index += 1
						col += 1

						if str[index-1] == '\n' {
							col = 1
							line += 1
						}

						foundSpace = true
					} else {
						for index < len(str) && str[index] != '\n' {
							index += 1
						}

						col = 1
					}
				}

				if foundSpace {
					continue
				}

				found := false

				remainingSubstring := str[index:]

				for _, pair := range lexicalDictionary {
					var (
						regex   = pair.regex
						handler = pair.handler
						pattern = regexp.MustCompile(regex)
						match   = pattern.FindStringSubmatch(remainingSubstring)
					)

					if len(match) > 0 {
						found = true
						t, literal, whole := handler(match)
						l := len(whole)

						ch <- token.Token{
							Type:    t,
							Literal: literal,
							Start:   token.Position{line, col},
							End:     token.Position{line, col + l - 1},
						}

						index += l
						col += l

						for index < len(str) && unicode.IsSpace(rune(str[index])) && str[index] != '\n' {
							index += 1
						}

						if index < len(str) && str[index] == '#' {
							for index < len(str) && str[index] != '\n' {
								index += 1
							}
						}

						isLineEnding := false

						for _, ending := range lineEndings {
							if t == ending {
								isLineEnding = true
							}
						}

						if (isLineEnding && index < len(str) && (str[index] == '\n' || str[index] == '}')) || index >= len(str) {
							ch <- token.Token{
								Type:    token.SEMI,
								Literal: ";",
								Start:   token.Position{line, col},
								End:     token.Position{line, col},
							}
						}

						break
					}
				}

				if !found {
					ch <- token.Token{
						Type:    token.ILLEGAL,
						Literal: string(str[index]),
						Start:   token.Position{line, col},
						End:     token.Position{line, col},
					}

					index += 1
					col += 1
				}
			} else {
				index += 1
				col += 1

				ch <- token.Token{
					Type:    token.EOF,
					Literal: "",
					Start:   token.Position{line, col},
					End:     token.Position{line, col},
				}
			}
		}
	}()

	return func() token.Token {
		return <-ch
	}
}
