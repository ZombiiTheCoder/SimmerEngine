package vm

import (
	"fmt"
	"io/fs"
	"os"
	opcode "simmer_js_engine/visitor/opcodes"
	"strconv"
)

type Chunk struct {
	Count     int
	Code      []byte
	Constants []Value
}

func InitChunk() *Chunk {
	chunk := Chunk{}
	chunk.Count = 0
	chunk.Code = make([]byte, 0)
	chunk.Constants = make([]Value, 0)
	return &chunk
}

func (c *Chunk) Write(Byte byte) {
	c.Code = append(c.Code, Byte)
	c.Count++
}

func (c *Chunk) Disassemble(name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < c.Count; {
		offset = c.DisassembleInstruction(offset)
	}
}

func (c *Chunk) DisassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)
	instruction := c.Code[offset]
	switch (instruction) {
	case opcode.OP_CONSTANT: return c.ConstantInstruction("OP_CONSTANT", offset)
	case opcode.OP_NULL: return c.SimpleInstruction("OP_NULL", offset)
	case opcode.OP_TRUE: return c.SimpleInstruction("OP_TRUE", offset)
	case opcode.OP_FALSE: return c.SimpleInstruction("OP_FALSE", offset)

	case opcode.OP_ADD: return c.SimpleInstruction("OP_ADD", offset)
	case opcode.OP_SUB: return c.SimpleInstruction("OP_SUB", offset)
	case opcode.OP_MULTIPLY: return c.SimpleInstruction("OP_MULTIPLY", offset)
	case opcode.OP_DIVIDE: return c.SimpleInstruction("OP_DIVIDE", offset)
	
	case opcode.OP_NEGATATE: return c.SimpleInstruction("OP_NEGATATE", offset)
	case opcode.OP_INCREMENT: return c.SimpleInstruction("OP_INCREMENT", offset)
	case opcode.OP_DECREMENT: return c.SimpleInstruction("OP_DECREMENT", offset)
	case opcode.OP_NEGATATE_BITWISE: return c.SimpleInstruction("OP_NEGATATE_BITWISE", offset)
	case opcode.OP_NOT: return c.SimpleInstruction("OP_NOT", offset)

	case opcode.OP_RETURN: return c.SimpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset+1
	}
}

func (c *Chunk) SimpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset+1
}

func (c *Chunk) ConstantInstruction(name string, offset int) int {
	constant := c.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, offset)
	fmt.Print(c.Constants[constant])
	fmt.Print("'\n")
	return offset+2
}

func (c *Chunk) AddConstant(value Value) int {
	c.Constants = append(c.Constants, value)
	return len(c.Constants)-1
}

func (c *Chunk) Out(name string) {
	_, e := os.Stat(name+".sbc")
	if (e == nil) { os.Remove(name+".sbc") }
	bc := ""
	for i := 0; i < len(c.Code); i++ {
		bc+=strconv.Itoa(int(c.Code[i]))
		if (i < len(c.Code)-1) { bc+="," }
	}
	co := ""
	for i := 0; i < len(c.Constants); i++ {
		co+=fmt.Sprintf("%v", c.Constants[i])
		if (i < len(c.Constants)-1) { co+="," }
	}
	os.WriteFile(name+".sbc", []byte(`{"CONSTANTS":[`+co+`],"BYTECODE":[`+fmt.Sprint(bc)+`]}`), fs.FileMode(os.O_RDWR))
}