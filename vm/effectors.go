package vm

import (
	"errors"
	"fmt"
	"math"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

// Effector is a function which performs a particular instruction
type Effector func(f *Frame, i bytecode.Instruction)

var effectors map[byte]Effector

// Effectors returns a map of bytecodes and their
// effectors. Wrapped in a function to stop an
// initialization loop.
func Effectors() map[byte]Effector {
	if len(effectors) == 0 {
		effectors = map[byte]Effector{
			bytecode.Pop: bytePop,
			bytecode.Dup: byteDup,

			bytecode.LoadConst: byteLoadConst,
			bytecode.LoadName:  byteLoadName,
			bytecode.StoreName: byteStoreName,

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
			bytecode.BinaryNotEqual: byteNotEqual,
			bytecode.BinaryLessThan: bincmp,
			bytecode.BinaryMoreThan: bincmp,
			bytecode.BinaryLessEq:   bincmp,
			bytecode.BinaryMoreEq:   bincmp,

			bytecode.Call:   byteCall,
			bytecode.Return: byteReturn,

			bytecode.Print:   bytePrint,
			bytecode.Println: bytePrintln,

			bytecode.Jump:        byteJump,
			bytecode.JumpIfTrue:  byteJumpIfTrue,
			bytecode.JumpIfFalse: byteJumpIfFalse,

			bytecode.MakeArray: byteMakeArray,
			bytecode.MakeTuple: byteMakeTuple,
			bytecode.MakeMap:   byteMakeMap,
		}
	}

	return effectors
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
		f.vm.Error = errors.New("evaluation: internal: name not found")
		return
	}

	val, ok := f.searchName(name)
	if !ok {
		f.vm.Error = fmt.Errorf("evaluation: name %s not found in the current scope", name)
		return
	}

	f.stack.push(val)
}

func byteStoreName(f *Frame, i bytecode.Instruction) {
	name, ok := f.getName(i.Arg)
	if !ok {
		f.vm.Error = errors.New("evaluation: internal: name not found")
		return
	}

	f.locals.Define(name, f.stack.pop())
}

func nunop(f *Frame, i bytecode.Instruction) {
	a := f.stack.pop()

	n, ok := a.(object.Numeric)
	if !ok {
		f.vm.Error = errors.New("evaluation: non-numeric value in numeric unary expression")
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
		f.vm.Error = errors.New("evaluation: non-numeric value in numeric binary expression")
		return
	}

	m, ok := b.(object.Numeric)
	if !ok {
		f.vm.Error = errors.New("evaluation: non-numeric value in numeric binary expression")
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

func bincmp(f *Frame, i bytecode.Instruction) {
	f.byteToInstructionIndex(int(i.Arg))

	b, a := f.stack.pop(), f.stack.pop()

	n, ok := a.(object.Numeric)
	if !ok {
		f.vm.Error = errors.New("evaluation: non-numeric value in numeric binary expression")
		return
	}

	m, ok := b.(object.Numeric)
	if !ok {
		f.vm.Error = errors.New("evaluation: non-numeric value in numeric binary expression")
		return
	}

	lval := n.Float64()
	rval := m.Float64()

	var result bool

	switch i.Code {
	case bytecode.BinaryLessThan:
		result = lval < rval
	case bytecode.BinaryMoreThan:
		result = lval > rval
	case bytecode.BinaryLessEq:
		result = lval <= rval
	case bytecode.BinaryMoreEq:
		result = lval >= rval
	}

	f.stack.push(&object.Boolean{Value: result})
}

func byteEquals(f *Frame, i bytecode.Instruction) {
	right, left := f.stack.pop(), f.stack.pop()
	eq := left.Equals(right)

	f.stack.push(object.BoolObj(eq))
}

func byteNotEqual(f *Frame, i bytecode.Instruction) {
	right, left := f.stack.pop(), f.stack.pop()
	eq := left.Equals(right)

	f.stack.push(object.BoolObj(!eq))
}

func byteCall(f *Frame, i bytecode.Instruction) {
	pattern := f.locals.Patterns[i.Arg]

	fn := f.locals.Functions.SearchString(pattern)
	if fn == nil {
		f.vm.Error = fmt.Errorf("evaluation: function '%s' not found in the current scope", pattern)
		return
	}

	locals := f.locals

	for _, item := range fn.Pattern {
		if param, ok := item.(*ast.Parameter); ok {
			// Found a parameter

			locals.Define(param.Name, f.stack.pop())
		}
	}

	// Create the function's frame
	fnFrame := &Frame{
		code:      fn.Body,
		constants: fn.Constants,
		locals:    locals,
		offset:    0,
		previous:  f,
		stack:     newStack(),
		vm:        f.vm,
	}

	// Push and execute the function's frame
	f.vm.pushFrame(fnFrame)
	f.vm.runFrame(fnFrame)

	// Push the returned value
	f.stack.push(fnFrame.stack.pop())
}

func byteReturn(f *Frame, i bytecode.Instruction) {
	f.offset = len(f.code) - 1
}

func bytePrint(f *Frame, i bytecode.Instruction) {
	fmt.Print(f.stack.pop())
}

func bytePrintln(f *Frame, i bytecode.Instruction) {
	fmt.Println(f.stack.pop())
}

func byteJump(f *Frame, i bytecode.Instruction) {
	f.offset = f.byteToInstructionIndex(int(i.Arg))
}

func byteJumpIfTrue(f *Frame, i bytecode.Instruction) {
	obj := f.stack.pop()

	if object.IsTruthy(obj) {
		f.offset = f.byteToInstructionIndex(int(i.Arg))
	}
}

func byteJumpIfFalse(f *Frame, i bytecode.Instruction) {
	obj := f.stack.pop()

	if !object.IsTruthy(obj) {
		f.offset = f.byteToInstructionIndex(int(i.Arg))
	}
}

func byteMakeArray(f *Frame, i bytecode.Instruction) {
	elems := make([]object.Object, i.Arg)

	for n := int(i.Arg) - 1; n >= 0; n-- {
		elems[n] = f.stack.pop()
	}

	f.stack.push(&object.Array{
		Value: elems,
	})
}

func byteMakeTuple(f *Frame, i bytecode.Instruction) {
	elems := make([]object.Object, i.Arg)

	for n := int(i.Arg) - 1; n >= 0; n-- {
		elems[n] = f.stack.pop()
	}

	f.stack.push(&object.Tuple{
		Value: elems,
	})
}

func byteMakeMap(f *Frame, i bytecode.Instruction) {
	keys := make(map[string]object.Object, i.Arg)
	values := make(map[string]object.Object, i.Arg)

	for n := 0; n < int(i.Arg); n++ {
		val, key := f.stack.pop(), f.stack.pop()

		hasher, ok := key.(object.Hasher)
		if !ok {
			f.vm.Error = fmt.Errorf("evaluation: non-hashable type as map key: %s", key.Type())
			return
		}

		hash := hasher.Hash()

		keys[hash] = key
		values[hash] = val
	}

	obj := &object.Map{
		Keys:   keys,
		Values: values,
	}

	f.stack.push(obj)
}
