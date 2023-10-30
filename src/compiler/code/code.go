package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte
type Opcode byte

type Definition struct {
	Name          string
	OperandWidths []int
}

const (
	OP_CONST Opcode = iota
	OP_NULL

	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD

	OP_POP

	OP_TRUE
	OP_FALSE

	OP_EQUAL
	OP_NOT_EQUAL
	OP_LESS_THAN
	OP_LESS_THAN_EQUAL
	OP_GREATER_THAN
	OP_GREATER_THAN_EQUAL

	OP_NOT
	OP_BNOT
	OP_UPLUS
	OP_UMINUS
	OP_PRINC
	OP_PRDEC
	OP_TYPEOF
	OP_VOID
	OP_DELETE

	OP_JMP
	OP_JIF
)

var Definitions = map[Opcode]*Definition{
	OP_CONST: {"OP_CONST", []int{2}},
	OP_NULL: {"OP_NULL", []int{}},
	OP_ADD: {"OP_ADD", []int{}},
	OP_SUB: {"OP_SUB", []int{}},
	OP_MUL: {"OP_MUL", []int{}},
	OP_DIV: {"OP_DIV", []int{}},
	OP_POP: {"OP_POP", []int{}},
	OP_MOD: {"OP_MOD", []int{}},
	OP_TRUE: {"OP_TRUE", []int{}},
	OP_FALSE: {"OP_FALSE", []int{}},
	OP_EQUAL: {"OP_EQUAL", []int{}},
	OP_NOT_EQUAL: {"OP_NOT_EQUAL", []int{}},
	OP_LESS_THAN: {"OP_LESS_THAN", []int{}},
	OP_LESS_THAN_EQUAL: {"OP_LESS_THAN_EQUAL", []int{}},
	OP_GREATER_THAN: {"OP_GREATER_THAN", []int{}},
	OP_GREATER_THAN_EQUAL: {"OP_GREATER_THAN_EQUAL", []int{}},

	OP_NOT: {"OP_NOT", []int{}},
	OP_BNOT: {"OP_BNOT", []int{}},
	OP_UPLUS: {"OP_UPLUS", []int{}},
	OP_UMINUS: {"OP_UMINUS", []int{}},
	OP_PRINC: {"OP_PRINC", []int{}},
	OP_PRDEC: {"OP_PRDEC", []int{}},
	OP_TYPEOF: {"OP_TYPEOF", []int{}},
	OP_VOID: {"OP_VOID", []int{}},
	OP_DELETE: {"OP_DELETE", []int{}},

	OP_JMP: {"OP_JMP", []int{2}},
	OP_JIF: {"OP_JIF", []int{2}},

}

func Lookup(op byte) (*Definition, error) {
	def, ok := Definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := Definitions[op]
	if !ok { return []byte{} }
	instrLen := 1
	for _, w := range def.OperandWidths { instrLen+=w }
	instr := make([]byte, instrLen)
	instr[0] = byte(op)
	offs := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2: binary.BigEndian.PutUint16(instr[offs:], uint16(o))
		}
		offs+=width
	}
	return instr
}

func (ins Instructions) String() string { 
	var out bytes.Buffer
	i := 0
	for i < len(ins) {
	def, err := Lookup(ins[i])
	if err != nil {
	fmt.Fprintf(&out, "ERROR: %s\n", err)
	continue
	}
	operands, read := ReadOperands(def, ins[i+1:])
	fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
	i += 1 + read
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
		len(operands), operandCount)
	}
	switch operandCount {
	case 0: return def.Name
	case 1: return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}


func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	ops := make([]int, len(def.OperandWidths))
	offs := 0
	for i, widths := range def.OperandWidths {
		switch widths {
		case 2:
			ops[i] = int(ReadUint16(ins[offs:]))
		}
	}
	return ops, offs
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}