package visitor

import (
	"fmt"
	"os"
	"simmer_js_engine/parser/ast"
	"simmer_js_engine/parser/tokens"
	opcode "simmer_js_engine/visitor/opcodes"
	"simmer_js_engine/vm"
	"strconv"
)

type Visitor struct {
	Chunk *vm.Chunk
	constants []float64
}

func InitVisitor() *Visitor {
	return &Visitor{
		Chunk: vm.InitChunk(),
		constants: make([]float64, 0),
	}
}

func (v *Visitor) Emit(op byte) {
	fmt.Printf("%v\n", op)
	v.Chunk.Write(op)
}
func (v *Visitor) EmitCONST(op byte, val vm.Value) {
	index := v.Chunk.AddConstant(val)
	if index > 255 {
		fmt.Println("TOO MANY CONSTANTS IN ONE CHUNK")
		os.Exit(-1)
	}
	fmt.Printf("%v %v\n", op, index)
	v.Chunk.Write(op)
	v.Chunk.Write(byte(index))
}

func (v *Visitor) Visit(node ast.Node) {
	switch node.GetNodeType() {
	case "Program": {
		for _, va := range node.(ast.Program).Body {
			v.Visit(va)
		}
	}
	case "ExprStmt": v.Visit(node.(ast.ExprStmt).Expression)

	case "PrefixExpr": v.VisitPrefixExpr(node.(ast.PrefixExpr))

	case "BinaryExpr": v.VisitBinaryExpr(node.(ast.BinaryExpr))

	case "NumberLiteral": {
		n, _ := strconv.ParseFloat(node.(ast.Literal).Value, 64)
		v.EmitCONST(opcode.OP_CONSTANT, vm.Value(n))
	}

	case "BooleanLiteral": {
		n, _ := strconv.ParseBool(node.(ast.Literal).Value)
		if n {
			v.Emit(opcode.OP_TRUE)
		} else {
			v.Emit(opcode.OP_FALSE)
		}
	}

	case "NullLiteral": {
		v.Emit(opcode.OP_NULL)
	}

	default: fmt.Println("ERR", node.GetNodeType())
	}
}

func (v *Visitor) VisitBinaryExpr(node ast.BinaryExpr) {
	v.Visit(node.Left)
	v.Visit(node.Right)
	switch node.Op{
	case tokens.Plus: v.Emit(opcode.OP_ADD)
	case tokens.Minus: v.Emit(opcode.OP_SUB)
	case tokens.Multiply: v.Emit(opcode.OP_MULTIPLY)
	case tokens.Divide: v.Emit(opcode.OP_DIVIDE)
	}
}

func (v *Visitor) VisitPrefixExpr(node ast.PrefixExpr) {
	v.Visit(node.Right)
	switch node.Op{
	case tokens.Plus: return
	case tokens.Minus: v.Emit(opcode.OP_NEGATATE)
	case tokens.BitwiseNot: v.Emit(opcode.OP_NEGATATE_BITWISE)
	case tokens.Increment: v.Emit(opcode.OP_INCREMENT)
	case tokens.Decrement: v.Emit(opcode.OP_DECREMENT)
	case tokens.LogicalNot: v.Emit(opcode.OP_NOT)
	}
}