package parser

import "github.com/Zac-Garby/pluto/token"

func (p *Parser) curIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.cur.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) peekIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.peek.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) expect(ts ...token.Type) bool {
	if p.peekIs(ts...) {
		p.next()
		return true
	}

	p.peekErr(ts...)
	return false
}

func (p *Parser) expectCur(ts ...token.Type) bool {
	if p.curIs(ts...) {
		p.next()
		return true
	}

	p.curErr(ts...)
	return false
}
