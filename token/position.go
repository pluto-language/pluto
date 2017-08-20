package token

import "fmt"

type Position struct {
	Line, Column int
}

func (p *Position) String() string {
	return fmt.Sprintf("%s:%s", p.Line, p.Column)
}
