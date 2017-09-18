package vm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/compiler"
	"github.com/Zac-Garby/pluto/dir"
	"github.com/Zac-Garby/pluto/object"
	"github.com/Zac-Garby/pluto/parser"
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

			bytecode.LoadConst:  byteLoadConst,
			bytecode.LoadName:   byteLoadName,
			bytecode.StoreName:  byteStoreName,
			bytecode.LoadField:  byteLoadField,
			bytecode.StoreField: byteStoreField,

			bytecode.UnaryInvert: bytePrefix,
			bytecode.UnaryNegate: bytePrefix,
			bytecode.UnaryNoOp:   bytePrefix,

			bytecode.BinaryAdd:      byteInfix,
			bytecode.BinarySubtract: byteInfix,
			bytecode.BinaryMultiply: byteInfix,
			bytecode.BinaryDivide:   byteInfix,
			bytecode.BinaryExponent: byteInfix,
			bytecode.BinaryFloorDiv: byteInfix,
			bytecode.BinaryMod:      byteInfix,
			bytecode.BinaryBitOr:    byteInfix,
			bytecode.BinaryBitAnd:   byteInfix,
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
			bytecode.Break:       byteBreak,
			bytecode.Next:        byteNext,
			bytecode.LoopStart:   byteLoopStart,
			bytecode.LoopEnd:     byteLoopEnd,

			bytecode.MakeArray: byteMakeArray,
			bytecode.MakeTuple: byteMakeTuple,
			bytecode.MakeMap:   byteMakeMap,

			bytecode.Use: byteUse,
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

	f.locals.Define(name, f.stack.top())
}

func byteLoadField(f *Frame, i bytecode.Instruction) {
	field, obj := f.stack.pop(), f.stack.pop()

	var val object.Object

	if col, ok := obj.(object.Collection); ok {
		if index, ok := field.(object.Numeric); ok {
			idx := int(index.Float64())

			val = col.GetIndex(idx)
		} else {
			f.vm.Error = fmt.Errorf("evaluation: non-numeric type %s used to index a collection", field.Type())
			return
		}
	} else if cont, ok := obj.(object.Container); ok {
		val = cont.Get(field)
	} else {
		f.vm.Error = fmt.Errorf("evaluation: cannot index type %s", obj.Type())
	}

	f.stack.push(val)
}

func byteStoreField(f *Frame, i bytecode.Instruction) {
	field, obj, val := f.stack.pop(), f.stack.pop(), f.stack.top()

	if col, ok := obj.(object.Collection); ok {
		if index, ok := field.(object.Numeric); ok {
			idx := int(index.Float64())

			col.SetIndex(idx, val)
		} else {
			f.vm.Error = fmt.Errorf("evaluation: non-numeric type %s used to index a collection", field.Type())
			return
		}
	} else if cont, ok := obj.(object.Container); ok {
		cont.Set(field, val)
	} else {
		f.vm.Error = fmt.Errorf("evaluation: cannot index type %s", obj.Type())
	}
}

func bytePrefix(f *Frame, i bytecode.Instruction) {
	right := f.stack.pop()

	if i.Code == bytecode.UnaryInvert {
		f.stack.push(object.BoolObj(!object.IsTruthy(right)))
		return
	}

	if n, ok := right.(object.Numeric); ok {
		val := n.Float64()
		f.stack.push(numPrefix(i.Code, val))
	} else {
		f.vm.Error = errors.New("evaluation: prefix r-value of invalid type")
	}
}

func numPrefix(opcode byte, val float64) object.Object {
	switch opcode {
	case bytecode.UnaryNegate:
		val = -val
	case bytecode.UnaryNoOp:
		// val = val
	}

	return &object.Number{Value: val}
}

func byteInfix(f *Frame, i bytecode.Instruction) {
	right, left := f.stack.pop(), f.stack.pop()

	if n, ok := left.(object.Numeric); ok {
		if m, ok := right.(object.Numeric); ok {
			f.stack.push(numInfix(f, i.Code, n.Float64(), m.Float64()))
		} else if m, ok := right.(object.Collection); ok {
			f.stack.push(numColInfix(f, i.Code, n.Float64(), m))
		} else {
			f.vm.Error = errors.New("evaluation: infix r-value of invalid type when l-value is <number>")
			return
		}
	} else if n, ok := left.(object.Collection); ok {
		if m, ok := right.(object.Numeric); ok {
			f.stack.push(numColInfix(f, i.Code, m.Float64(), n))
		} else if m, ok := right.(object.Collection); ok {
			f.stack.push(colInfix(f, i.Code, n, m))
		} else {
			f.vm.Error = errors.New("evaluation: infix r-value of invalid type when l-value is a collection")
		}
	} else {
		f.vm.Error = errors.New("evaluation: infix l-value of invalid type")
		return
	}
}

func numInfix(f *Frame, opcode byte, left, right float64) object.Object {
	var val float64

	switch opcode {
	case bytecode.BinaryAdd:
		val = left + right
	case bytecode.BinarySubtract:
		val = left - right
	case bytecode.BinaryMultiply:
		val = left * right
	case bytecode.BinaryDivide:
		val = left / right
	case bytecode.BinaryExponent:
		val = math.Pow(left, right)
	case bytecode.BinaryFloorDiv:
		val = math.Floor(left / right)
	case bytecode.BinaryMod:
		val = math.Mod(left, right)
	case bytecode.BinaryBitOr:
		val = float64(int64(left) | int64(right))
	case bytecode.BinaryBitAnd:
		val = float64(int64(left) & int64(right))
	default:
		op := bytecode.Instructions[opcode].Name[7:]
		f.vm.Error = fmt.Errorf("evaluation: operator %s not supported for two numbers", op)
	}

	return &object.Number{Value: val}
}

func numColInfix(f *Frame, opcode byte, left float64, right object.Collection) object.Object {
	var (
		result   []object.Object
		elements = right.Elements()
	)

	if opcode == bytecode.BinaryMultiply {
		for i := 0; i < int(left); i++ {
			result = append(result, elements...)
		}
	} else {
		op := bytecode.Instructions[opcode].Name[7:]
		f.vm.Error = fmt.Errorf("evaluation: operator %s not supported for a collection and a number", op)
	}

	col, _ := object.MakeCollection(right.Type(), result)
	return col
}

func colInfix(f *Frame, opcode byte, left, right object.Collection) object.Object {
	var (
		lefts  = left.Elements()
		rights = right.Elements()
		elems  []object.Object
	)

	switch opcode {
	case bytecode.BinaryAdd:
		elems = append(lefts, rights...)
	case bytecode.BinarySubtract:
		for _, el := range lefts {
			for _, rel := range rights {
				if el.Equals(rel) {
					goto next
				}
			}

			elems = append(elems, el)
		next:
		}
	case bytecode.BinaryBitOr:
		for _, el := range append(lefts, rights...) {
			unique := true

			for _, rel := range elems {
				if el.Equals(rel) {
					unique = false
					break
				}
			}

			if unique {
				elems = append(elems, el)
			}
		}
	case bytecode.BinaryBitAnd:
		for _, el := range lefts {
			both := false

			for _, rel := range rights {
				if el.Equals(rel) {
					both = true
					break
				}
			}

			if both {
				elems = append(elems, el)
			}
		}
	default:
		op := bytecode.Instructions[opcode].Name[7:]
		f.vm.Error = fmt.Errorf("evaluation: operator %s not supported for two collections", op)
	}

	col, _ := object.MakeCollection(left.Type(), elems)
	return col
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

	fn := f.locals.FunctionStore.SearchString(pattern)
	if fn == nil {
		f.vm.Error = fmt.Errorf("evaluation: function '%s' not found in the current scope", pattern)
		return
	}

	locals := f.locals

	locals.Names = fn.Names
	locals.Patterns = fn.Patterns

	for _, item := range fn.Pattern {
		if param, ok := item.(*ast.Parameter); ok {
			// Found a parameter

			locals.Data[param.Name] = f.stack.pop()
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
	f.vm.runFrame(fnFrame)

	if len(fnFrame.stack.objects) > 0 {
		ret := fnFrame.stack.pop()

		// Push the returned value
		f.stack.push(ret)
	}
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

func byteBreak(f *Frame, i bytecode.Instruction) {
	top := f.breaks[len(f.breaks)-1]
	f.offset = top
}

func byteNext(f *Frame, i bytecode.Instruction) {
	top := f.nexts[len(f.nexts)-1]
	f.offset = top
}

func byteLoopStart(f *Frame, i bytecode.Instruction) {
	f.nexts = append(f.nexts, f.offset+1)

	var o int

	for o = f.offset; f.code[o].Code != bytecode.LoopEnd; o++ {
		// Nothing here
	}

	f.breaks = append(f.breaks, o)
}

func byteLoopEnd(f *Frame, i bytecode.Instruction) {
	f.breaks = f.breaks[:len(f.breaks)-1]
	f.nexts = f.nexts[:len(f.nexts)-1]
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

func byteUse(f *Frame, i bytecode.Instruction) {
	path := f.constants[i.Arg].String()

	sources, err := dir.LocateAnySources(path)
	if err != nil {
		f.vm.Error = err
		return
	}

	for _, source := range sources {
		src, err := ioutil.ReadFile(source)
		if err != nil {
			f.vm.Error = err
		}

		var (
			str   = string(src)
			cmp   = compiler.New()
			parse = parser.New(str, source)
			prog  = parse.Parse()
		)

		if len(parse.Errors) > 0 {
			parse.PrintErrors()
			f.vm.Error = errors.New("use: parse error")
			return
		}

		err = cmp.CompileProgram(prog)
		if err != nil {
			f.vm.Error = err
			return
		}

		code, err := bytecode.Read(cmp.Bytes)
		if err != nil {
			f.vm.Error = err
			return
		}

		store := *f.locals
		store.Names = cmp.Names
		store.Functions = cmp.Functions
		store.Patterns = cmp.Patterns

		machine := New()
		machine.Run(code, &store, cmp.Constants)

		f.locals.Extend(&store)
	}
}
