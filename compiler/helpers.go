package compiler

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
