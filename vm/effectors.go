package vm

import (
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

type effector func(f *Frame, i bytecode.Instruction)

var effectors = map[byte]effector{
	bytecode.Pop: bytePop,
	bytecode.Dup: byteDup,

	bytecode.LoadConst: byteLoadConst,

	bytecode.BinaryAdd: byteBinaryAdd,
}

func bytePop(f *Frame, i bytecode.Instruction) {
	f.stack.pop()
}

func byteDup(f *Frame, i bytecode.Instruction) {
	f.stack.dup()
}

func byteLoadConst(f *Frame, i bytecode.Instruction) {
	f.stack.push(f.constants[i.Arg])
}

func byteBinaryAdd(f *Frame, i bytecode.Instruction) {
	b, a := f.stack.pop(), f.stack.pop()

	n, _ := a.(*object.Number)
	m, _ := b.(*object.Number)

	f.stack.push(&object.Number{Value: n.Value + m.Value})
}
