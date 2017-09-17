package token

import "fmt"

// Position is a token's position in the source
type Position struct {
	Line, Column int
}

func (p *Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}
