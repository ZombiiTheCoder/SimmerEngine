package vm

import (
	"fmt"
	"simmer_js_engine/compiler"
	"simmer_js_engine/compiler/code"
	"simmer_js_engine/compiler/object"
)

var (
	True = object.Boolean{Value: true}
	False = object.Boolean{Value: false}
	Null = object.Null{Value: nil}
)

const StackSize = 2048

type VM struct {
	conts []object.Object
	instr code.Instructions
	stack []object.Object
	stackPtr int
}

func InitVM(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instr: bytecode.Instrs,
		conts: bytecode.Constants,
		stack: make([]object.Object, StackSize),
		stackPtr: 0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPtr == 0 { return nil }
	return vm.stack[vm.stackPtr-1]
}

func (vm *VM) ExecuteBangOp() error {
	operand := vm.pop()

	switch operand {
	case True: return vm.push(False)
	case False: return vm.push(True)
	case Null: return vm.push(True)
	default: return vm.push(False)
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instr); ip++ {
		op := code.Opcode(vm.instr[ip])
		switch op {
		case code.OP_CONST:
 			constIndex := code.ReadUint16(vm.instr[ip+1:])
			err := vm.push(vm.conts[constIndex])
			if err != nil { return err }
			ip+=2
		case code.OP_ADD, code.OP_SUB, code.OP_MUL, code.OP_DIV, code.OP_MOD:
			err := vm.ExecuteBinaryOp(op)
			if err != nil { return err }
		case
		code.OP_LESS_THAN,
		code.OP_LESS_THAN_EQUAL,

		code.OP_GREATER_THAN,
		code.OP_GREATER_THAN_EQUAL,

		code.OP_EQUAL,
		code.OP_NOT_EQUAL:
			err := vm.ExecuteComparison(op)
			if err != nil { return err }
		case code.OP_POP: vm.pop()
		case code.OP_NULL:
			err := vm.push(Null)
			if err != nil { return err }
		case code.OP_TRUE:
			err := vm.push(True)
			if err != nil { return err }
		case code.OP_FALSE:
			err := vm.push(False)
			if err != nil { return err }
		case code.OP_NOT:
			err := vm.ExecuteBangOp()
			if err != nil { return err }
		case code.OP_UPLUS, code.OP_UMINUS, code.OP_BNOT, code.OP_PRINC, code.OP_PRDEC:
			err := vm.ExecutePrefix(op)
			if err != nil { return err }
		case code.OP_JIF:
			pos := int(code.ReadUint16(vm.instr[ip+1:]))
			c := vm.pop()
			ip+=2
			if !object.IsTruthy(c) { ip = pos - 1 }

		case code.OP_JMP:
			pos := int(code.ReadUint16(vm.instr[ip+1:]))
			ip = pos - 1
		}
	}
	fmt.Println(vm.lastPopped().GetValue())
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.stackPtr >= StackSize { return fmt.Errorf("stack overflow") }
	vm.stack[vm.stackPtr] = o
	vm.stackPtr++
	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.stackPtr-1]
	vm.stackPtr--
	return o
}

func (vm *VM) lastPopped() object.Object { return vm.stack[vm.stackPtr] }