package vm

import (
	"math"
	"strings"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/object"
)

// Effector is a function which performs a particular instruction
type Effector func(f *Frame, i bytecode.Instruction)

var effectors map[byte]Effector

func init() {
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

		bytecode.PushFn:     bytePushFn,
		bytecode.PushQualFn: bytePushQualFn,
		bytecode.CallFn:     byteCall,
		bytecode.Return:     byteReturn,
		bytecode.DoBlock:    byteDoBlock,

		bytecode.Print:   bytePrint,
		bytecode.Println: bytePrintln,
		bytecode.Length:  byteLength,

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
		f.vm.Error = Err("name not defined", ErrInternal)
		return
	}

	val, ok := f.searchName(name)
	if !ok {
		f.vm.Error = Errf("name %s not found in the current scope", ErrNotFound, name)
		return
	}

	f.stack.push(val)
}

func byteStoreName(f *Frame, i bytecode.Instruction) {
	name, ok := f.getName(i.Arg)
	if !ok {
		f.vm.Error = Err("name not defined", ErrInternal)
		return
	}

	f.locals.Define(name, f.stack.top(), true)
}

func byteLoadField(f *Frame, i bytecode.Instruction) {
	field, obj := f.stack.pop(), f.stack.pop()

	var val object.Object

	if col, ok := obj.(object.Collection); ok {
		if index, ok := field.(object.Numeric); ok {
			idx := int(index.Float64())

			val = col.GetIndex(idx)
		} else {
			f.vm.Error = Errf("non-numeric type %s used to index a collection", ErrWrongType, field.Type())
			return
		}
	} else if cont, ok := obj.(object.Container); ok {
		val = cont.Get(field)
	} else {
		f.vm.Error = Errf("cannot index type %s", ErrNotFound, obj.Type())
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
			f.vm.Error = Errf("non-numeric type %s used to index a collection", ErrWrongType, field.Type())
			return
		}
	} else if cont, ok := obj.(object.Container); ok {
		cont.Set(field, val)
	} else {
		f.vm.Error = Errf("cannot index type %s", ErrWrongType, obj.Type())
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
		f.vm.Error = Err("prefix r-value of invalid type", ErrWrongType)
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
			f.vm.Error = Err("infix r-value of invalid type when l-value is <number>", ErrWrongType)
			return
		}
	} else if n, ok := left.(object.Collection); ok {
		if m, ok := right.(object.Numeric); ok {
			f.stack.push(numColInfix(f, i.Code, m.Float64(), n))
		} else if m, ok := right.(object.Collection); ok {
			f.stack.push(colInfix(f, i.Code, n, m))
		} else {
			f.vm.Error = Err("infix r-value of invalid type when l-value is a collection", ErrWrongType)
		}
	} else {
		f.vm.Error = Err("infix l-value of invalid type", ErrWrongType)
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
		f.vm.Error = Errf("operator %s not supported for two numbers", ErrNoOp, op)
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
		f.vm.Error = Errf("operator %s not supported for a collection and a number", ErrNoOp, op)
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
		f.vm.Error = Errf("operator %s not supported for two collections", ErrNoOp, op)
	}

	col, _ := object.MakeCollection(left.Type(), elems)
	return col
}

func bincmp(f *Frame, i bytecode.Instruction) {
	f.byteToInstructionIndex(int(i.Arg))

	b, a := f.stack.pop(), f.stack.pop()

	n, ok := a.(object.Numeric)
	if !ok {
		f.vm.Error = Err("non-numeric value in numeric binary expression", ErrWrongType)
		return
	}

	m, ok := b.(object.Numeric)
	if !ok {
		f.vm.Error = Err("non-numeric value in numeric binary expression", ErrWrongType)
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

func bytePushFn(f *Frame, i bytecode.Instruction) {
	pattern := f.locals.Patterns[i.Arg]

	fn := f.locals.FunctionStore.SearchString(pattern)
	if fn == nil {
		f.vm.Error = Errf("function '%s' not found in the current scope", ErrNotFound, pattern)
		return
	}

	f.stack.push(fn)
}

func bytePushQualFn(f *Frame, i bytecode.Instruction) {
	pattern := strings.Split(f.locals.Patterns[i.Arg], " ")

	baseObj := f.stack.pop()

	if baseObj.Type() != object.MapType {
		f.vm.Error = Errf("cannot call a method of non-map type %s", ErrWrongType, baseObj.Type())
		return
	}

	var (
		base    = baseObj.(*object.Map)
		methods = base.Get(&object.String{Value: "_methods"})
	)

	if methods == nil {
		f.vm.Error = Err("_methods key not found", ErrWrongType)
		return
	}

	if methods.Type() != object.ArrayType {
		f.vm.Error = Err("_methods is not an array", ErrWrongType)
		return
	}

	methArr := methods.(*object.Array)

outer:
	for _, obj := range methArr.Elements() {
		fn, ok := obj.(*object.Function)
		if !ok {
			continue
		}

		fnpat := fn.Pattern

		if len(fnpat) != len(pattern) {
			// Doesn't match
			continue outer
		}

		for i, item := range pattern {
			var (
				fItem = fnpat[i]
				isArg = item[0] == '$'
			)

			if isArg {
				if _, ok := fItem.(*ast.Parameter); !ok {
					// Doesn't match
					continue outer
				}
			} else {
				if id, ok := fItem.(*ast.Identifier); !ok {
					// Doesn't match
					continue outer
				} else if id.Value != item {
					// Doesn't match
					continue outer
				}
			}
		}

		f.stack.push(fn)
		return
	}

	f.vm.Error = Errf("no method was found matching the pattern: '%s'", ErrNotFound, strings.Join(pattern, " "))
}

func byteCall(f *Frame, i bytecode.Instruction) {
	fn, ok := f.stack.pop().(*object.Function)
	if !ok {
		f.vm.Error = Errf("cannot call non-function type: %s", ErrWrongType, fn.Type())
		return
	}

	locals := f.locals

	locals.Names = fn.Names
	locals.Patterns = fn.Patterns

	for _, item := range fn.Pattern {
		if param, ok := item.(*ast.Parameter); ok {
			// Found a parameter

			locals.Define(param.Name, f.stack.pop(), true)
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

func byteDoBlock(f *Frame, i bytecode.Instruction) {
	top := f.stack.pop()

	block, ok := top.(*object.Block)
	if !ok {
		f.vm.Error = Errf("cannot 'do' a non-block. got %s", ErrWrongType, top.Type())
		return
	}

	locals := f.locals

	locals.Names = block.Names
	locals.Patterns = block.Patterns

	for _, item := range block.Params {
		name := item.Token().Literal

		locals.Define(name, f.stack.pop(), true)
	}

	blockFrame := &Frame{
		code:      block.Body,
		constants: block.Constants,
		locals:    locals,
		offset:    0,
		previous:  f,
		stack:     newStack(),
		vm:        f.vm,
	}

	f.vm.runFrame(blockFrame)

	if len(blockFrame.stack.objects) > 0 {
		ret := blockFrame.stack.pop()

		f.stack.push(ret)
	}
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
	if len(f.breaks) < 1 {
		f.vm.Error = Err("break statement found outside loop", ErrSyntax)
		return
	}

	top := f.breaks[len(f.breaks)-1]
	f.offset = top
}

func byteNext(f *Frame, i bytecode.Instruction) {
	if len(f.nexts) < 1 {
		f.vm.Error = Err("next statement found outside loop", ErrSyntax)
		return
	}

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
			f.vm.Error = Errf("non-hashable type as map key: %s", ErrWrongType, key.Type())
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

	f.Use(path)
}
