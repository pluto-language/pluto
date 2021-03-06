package compiler

import (
	"fmt"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

const maxRune = 1 << 16

func runeToBytes(x rune) (byte, byte) {
	var (
		low  = byte(x & 0xff)
		high = byte((x >> 8) & 0xff)
	)

	return low, high
}

func bytesToRune(low, high byte) rune {
	return (rune(high) << 8) | rune(low)
}

func (c *Compiler) addConst(val object.Object) (rune, error) {
	for i, cst := range c.Constants {
		if val.Equals(cst) {
			return rune(i), nil
		}
	}

	c.Constants = append(c.Constants, val)
	index := len(c.Constants) - 1

	if index >= maxRune {
		return 0, fmt.Errorf("compiler: constant index %d greater than 0xFFFF (maximum uint16)", index)
	}

	return rune(index), nil
}

func (c *Compiler) loadConst(index rune) {
	low, high := runeToBytes(index)
	c.push(bytecode.LoadConst, high, low)
}

func (c *Compiler) loadName(index rune) {
	low, high := runeToBytes(index)

	c.push(bytecode.LoadName, high, low)
}

func (c *Compiler) push(bytes ...byte) {
	c.Bytes = append(c.Bytes, bytes...)
}
