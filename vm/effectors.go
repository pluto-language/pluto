package vm

import (
	"errors"
	"fmt"
	"math"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

type effector func(f *Frame, i bytecode.Instruction)

var effectors = map[byte]effector{
	bytecode.Pop: bytePop,
	bytecode.Dup: byteDup,

	bytecode.LoadConst: byteLoadConst,
	bytecode.LoadName:  byteLoadName,

	bytecode.UnaryNegate: nunop,
	bytecode.UnaryNoOp:   nunop,

	bytecode.BinaryAdd:      nbinop,
	bytecode.BinarySubtract: nbinop,
	bytecode.BinaryMultiply: nbinop,
	bytecode.BinaryDivide:   nbinop,
	bytecode.BinaryExponent: nbinop,
	bytecode.BinaryFloorDiv: nbinop,
	bytecode.BinaryMod:      nbinop,
	bytecode.BinaryBitOr:    nbinop,
	bytecode.BinaryBitAnd:   nbinop,
	bytecode.BinaryEquals:   byteEquals,
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

func byteLoadName(f *Frame, i bytecode.Instruction) {
	name, ok := f.getName(i.Arg)
	if !ok {
		f.vm.lastError = errors.New("evaluation: internal: name not found")
		return
	}

	val, ok := f.searchName(name)
	if !ok {
		f.vm.lastError = fmt.Errorf("evaluation: name %s not found in the current scope", name)
		return
	}

	f.stack.push(val)
}

func nunop(f *Frame, i bytecode.Instruction) {
	a := f.stack.pop()

	n, ok := a.(object.Numeric)
	if !ok {
		f.vm.lastError = errors.New("evaluation: non-numeric value in numeric unary expression")
		return
	}

	v := n.Float64()

	var val float64

	switch i.Code {
	case bytecode.UnaryNegate:
		val = -v
	case bytecode.UnaryNoOp:
		val = v
	}

	f.stack.push(&object.Number{Value: val})
}

func nbinop(f *Frame, i bytecode.Instruction) {
	b, a := f.stack.pop(), f.stack.pop()

	n, ok := a.(object.Numeric)
	if !ok {
		f.vm.lastError = errors.New("evaluation: non-numeric value in numeric binary expression")
		return
	}

	m, ok := b.(object.Numeric)
	if !ok {
		f.vm.lastError = errors.New("evaluation: non-numeric value in numeric binary expression")
		return
	}

	lval := n.Float64()
	rval := m.Float64()

	var val float64

	switch i.Code {
	case bytecode.BinaryAdd:
		val = lval + rval
	case bytecode.BinarySubtract:
		val = lval - rval
	case bytecode.BinaryMultiply:
		val = lval * rval
	case bytecode.BinaryDivide:
		val = lval / rval
	case bytecode.BinaryExponent:
		val = math.Pow(lval, rval)
	case bytecode.BinaryFloorDiv:
		val = math.Floor(lval / rval)
	case bytecode.BinaryMod:
		val = math.Mod(lval, rval)
	case bytecode.BinaryBitOr:
		val = float64(int64(lval) | int64(rval))
	case bytecode.BinaryBitAnd:
		val = float64(int64(lval) & int64(rval))
	}

	f.stack.push(&object.Number{Value: val})
}

func byteEquals(f *Frame, i bytecode.Instruction) {
	right, left := f.stack.pop(), f.stack.pop()
	eq := left.Equals(right)

	f.stack.push(object.BoolObj(eq))
}
