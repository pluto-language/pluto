package graph

import "github.com/Zac-Garby/pluto/bytecode"

func byteToInstructionIndex(code bytecode.Code, b int) int {
	var index, counter int

	for _, instr := range code {
		if bytecode.Instructions[instr.Code].HasArg {
			counter += 3
		} else {
			counter++
		}

		if counter >= b {
			return index
		}

		index++
	}

	return index
}
