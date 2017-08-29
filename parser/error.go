package parser

import (
	"fmt"
	"strings"

	"github.com/Zac-Garby/pluto/token"
)

type Error struct {
	Message    string
	Start, End token.Position
}

func (p *Parser) err(msg string, start, end token.Position) {
	err := Error{
		Message: msg,
		Start:   start,
		End:     end,
	}

	p.Errors = append(p.Errors, err)
}

func (p *Parser) defaultErr(msg string) {
	err := Error{
		Message: msg,
		Start:   p.cur.Start,
		End:     p.cur.End,
	}

	p.Errors = append(p.Errors, err)
}

func (p *Parser) peekErr(ts ...token.Type) {
	if len(ts) > 1 {
		msg := "expected either "

		for i, t := range ts {
			msg += string(t)

			if i+1 < len(ts) {
				msg += ", "
			} else if i < len(ts) {
				msg += ", or "
			}
		}

		msg += ", but got " + string(p.peek.Type)

		p.err(msg, p.peek.Start, p.peek.End)
	} else if len(ts) == 1 {
		msg := fmt.Sprintf("expected %s, but got %s", ts[0], p.peek.Type)
		p.err(msg, p.peek.Start, p.peek.End)
	}
}

func (p *Parser) curErr(ts ...token.Type) {
	if len(ts) > 1 {
		msg := "expected either "

		for i, t := range ts {
			msg += string(t)

			if i+1 < len(ts) {
				msg += ", "
			} else if i < len(ts) {
				msg += ", or "
			}
		}

		msg += ", but got " + string(p.cur.Type)

		p.err(msg, p.cur.Start, p.cur.End)
	} else if len(ts) == 1 {
		msg := fmt.Sprintf("expected %s, but got %s", ts[0], p.cur.Type)
		p.err(msg, p.cur.Start, p.cur.End)
	}
}

func (p *Parser) unexpectedTokenErr(t token.Type) {
	msg := fmt.Sprintf("unexpected token: %s", t)
	p.defaultErr(msg)
}

func (p *Parser) printError(index int) {
	err := p.Errors[index]

	fmt.Printf("%s → %s\t%s\n", err.Start.String(), err.End.String(), err.Message)
}

func (p *Parser) printErrorVerbose(index int) {
	err := p.Errors[index]
	lines := strings.Split(p.text, "\n")

	fmt.Printf("    %d| %s\n", err.Start.Line, lines[err.Start.Line-1])
	fmt.Printf(
		"    %s %s%s\n",
		strings.Repeat(" ", len(fmt.Sprintf("%d", err.Start.Line))),
		strings.Repeat(" ", err.Start.Column),
		strings.Repeat("^", err.End.Column-err.Start.Column+1),
	)

	fmt.Printf("%s → %s\t%s\n\n", err.Start.String(), err.End.String(), err.Message)
}

func (p *Parser) PrintErrors() {
	for i := range p.Errors {
		if i == 0 {
			p.printErrorVerbose(i)
		} else {
			p.printError(i)
		}
	}
}
