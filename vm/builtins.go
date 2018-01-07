package vm

import (
	"fmt"

	"github.com/Zac-Garby/pluto/object"

	"github.com/Zac-Garby/pluto/bytecode"
)

func bytePrint(f *Frame, i bytecode.Instruction) {
	fmt.Print(f.stack.pop())
}

func bytePrintln(f *Frame, i bytecode.Instruction) {
	fmt.Println(f.stack.pop())
}

func byteLength(f *Frame, i bytecode.Instruction) {
	top := f.stack.pop()

	if col, ok := top.(object.Collection); ok {
		f.stack.push(&object.Number{
			Value: float64(len(col.Elements())),
		})
	} else {
		f.vm.Error = Errf("cannot get the length of type %s", ErrWrongType, top.Type())
	}
}
