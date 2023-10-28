package vm

import (
	"fmt"
	opcode "simmer_js_engine/visitor/opcodes"
)

type Value float64

var STACK_MAX = 256

func addConst(chunk *Chunk, value Value) {
	constant := chunk.AddConstant(value)
	chunk.Write(opcode.OP_CONSTANT)
	chunk.Write(byte(constant))
}

func main() {
	chunk := InitChunk()

	// addConst(chunk, 3.4)
	// addConst(chunk, 5.6)
	// chunk.Write(opcode.OP_DIVIDE)
	// addConst(chunk, 1.2)
	// chunk.Write(opcode.OP_ADD)

	addConst(chunk, 1.2)
	addConst(chunk, 3.4)
	// chunk.Write(opcode.OP_NEGATATE)
	chunk.Write(opcode.OP_ADD)
	addConst(chunk, 5.6)
	chunk.Write(opcode.OP_DIVIDE)

	chunk.Write(opcode.OP_RETURN)
	chunk.Disassemble("E")
	chunk.Out("E")
	// File := os.Args[0];
	vm := InitVM()
	fmt.Println(vm.Interpret(chunk))
}

type VM struct {
	chunk    *Chunk
	code     []byte
	ip       int
	stack    []Value
	stackPtr int
}

type InterpretResult string

const (
	INTERPRET_OK            InterpretResult = "INTERPRET_OK"
	INTERPRET_COMPILE_ERROR                 = "INTERPRET_COMPILE_ERROR"
	INTERPRET_RUNTIME_ERROR                 = "INTERPRET_RUNTIME_ERROR"
)

func InitVM() *VM {
	vm := VM{}
	return &vm
}

func (v *VM) Interpret(chunk *Chunk) InterpretResult {
	v.chunk = chunk
	v.code = v.chunk.Code
	v.ip = 0
	v.stack = make([]Value, 0)
	v.stackPtr = 0
	return v.run()
}

func (v *VM) READ_BYTE() byte {
	b := v.code[v.ip]
	v.ip++
	return b
}

func (v *VM) READ_CONSTANT() Value {
	return v.chunk.Constants[v.READ_BYTE()]
}


func (v *VM) push(value Value) {
	newStack := make([]Value, 0)
	newStack = append(newStack, value)
	v.stack = append(newStack, v.stack...)
	v.stackPtr++
}

func (v *VM) pop() Value {
	v.stackPtr--
	item := v.stack[v.stackPtr]
	v.stack = v.stack[:1]
	return item
}

func (v *VM) run() InterpretResult {
	for {
		instruction := v.READ_BYTE()
		switch instruction {
		case opcode.OP_CONSTANT:
			{
				constant := v.READ_CONSTANT()
				v.push(constant)
				break
			}
		// case opcode.OP_FALSE: v.push(false)
		// case opcode.OP_TRUE: v.push(true)
		// case opcode.OP_NULL: v.push(nil)
		case opcode.OP_ADD:
			{
				a := v.pop()
				b := v.pop()
				v.push(a + b)
			}
		case opcode.OP_SUB:
			{
				a := v.pop()
				b := v.pop()
				v.push(a - b)
			}
		case opcode.OP_MULTIPLY:
			{
				a := v.pop()
				b := v.pop()
				v.push(a * b)
			}
		case opcode.OP_DIVIDE:
			{
				a := v.pop()
				b := v.pop()
				fmt.Println(a, b)
				v.push(a / b)
			}
		case opcode.OP_NEGATATE:
			{
				v.push(-v.pop())
				break
			}
		case opcode.OP_RETURN:
			{
				fmt.Println(v.pop())
				return INTERPRET_OK
			}
		}
	}
}