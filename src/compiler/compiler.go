package compiler

import (
	"fmt"
	"simmer_js_engine/compiler/code"
	"simmer_js_engine/compiler/object"
	"simmer_js_engine/parser/ast"
	"simmer_js_engine/parser/tokens"
	"strconv"
	"strings"
)

type Compiler struct {
	instrs code.Instructions
	constants []object.Object
	lastInstruction EmittedInstruction
	previousInstruction EmittedInstruction
}

func InitCompiler() *Compiler {
	return &Compiler{
		instrs: code.Instructions{},
		constants: []object.Object{},
		lastInstruction: EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instrs: c.instrs,
		Constants: c.constants,
	}
}

type Bytecode struct {
	Instrs code.Instructions
	Constants []object.Object
}

type EmittedInstruction struct {
	Opcode code.Opcode
	Position int
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	// Literals
	case ast.NumberLiteral:
		num := object.Number{Value: node.Value}
		c.emit(code.OP_CONST, c.addConst(num))

	case ast.BooleanLiteral:
		if node.Value {
			c.emit(code.OP_TRUE)
		} else {
			c.emit(code.OP_FALSE)
		}

	case ast.IfStmt:
		err := c.Compile(node.Condition)
		if err != nil { return err }
		jumpFPos := c.emit(code.OP_JIF, 69)
		err = c.Compile(node.Consequent)
		if err != nil { return err }
		// if c.lastInstruction.Opcode == code.OP_POP {
		// 	c.instrs = c.instrs[:c.lastInstruction.Position]
		// 	c.lastInstruction = c.previousInstruction
		// }
		if node.Other != nil {
			jumpPos := c.emit(code.OP_JMP, 69)
			
			c.changeOperand(jumpFPos, len(c.instrs))

			err = c.Compile(node.Other)
			if err != nil { return err }
			afterAlt := len(c.instrs)
			c.changeOperand(jumpPos, afterAlt)
			return nil
		}
		c.changeOperand(jumpFPos, len(c.instrs))
		
	case ast.BodyStmt:
		for _, s := range node.Body {
			err := c.Compile(s)
			if err != nil { return err }
		}

	case ast.Program:
		for _, s := range node.Body {
			err := c.Compile(s)
			if err != nil { return err }
		}
	case ast.ExprStmt:
		err := c.Compile(node.Expression)
		if err != nil { return err }
		c.emit(code.OP_POP)
	case ast.BinaryExpr:
		err := c.Compile(node.Left)
		if err != nil { return err }
		err = c.Compile(node.Right)
		if err != nil { return err }
			
		switch node.Op {
		case tokens.Plus: c.emit(code.OP_ADD)
		case tokens.Minus: c.emit(code.OP_SUB)
		case tokens.Multiply: c.emit(code.OP_MUL)
		case tokens.Divide: c.emit(code.OP_DIV)
		case tokens.Modulo: c.emit(code.OP_MOD)

		case tokens.Equals: c.emit(code.OP_EQUAL)
		case tokens.NotEquals: c.emit(code.OP_NOT_EQUAL)
		case tokens.BiggerThan: c.emit(code.OP_GREATER_THAN)
		case tokens.BiggerThanEqualTo: c.emit(code.OP_GREATER_THAN_EQUAL)
		case tokens.SmallerThan: c.emit(code.OP_LESS_THAN)
		case tokens.SmallerThanEqualTo: c.emit(code.OP_LESS_THAN_EQUAL)
		default: return fmt.Errorf("UNKNOWN OP %s", node.Op)
		}
	case ast.PrefixExpr:
		err := c.Compile(node.Right)
		if err != nil { return err }
		switch node.Op {
			case tokens.Plus: c.emit(code.OP_UPLUS)
			case tokens.Minus: c.emit(code.OP_UMINUS)
			case tokens.LogicalNot: c.emit(code.OP_NOT)
			case tokens.BitwiseNot: c.emit(code.OP_BNOT)

			case tokens.Increment: c.emit(code.OP_PRINC)
			case tokens.Decrement: c.emit(code.OP_PRDEC)
			case tokens.Typeof: c.emit(code.OP_TYPEOF)
			case tokens.Void: c.emit(code.OP_VOID)
			case tokens.Delete: c.emit(code.OP_DELETE)
		default: return fmt.Errorf("UNKNOWN OP %s", node.Op)
		}
	}
	return nil
}

func (c *Compiler) addConst(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants)-1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstr(ins)

	prev := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	c.previousInstruction = prev
	c.lastInstruction = last
	
	if len(operands) > 0 {
		var p []string
		for _, v := range operands { p = append(p, strconv.FormatInt(int64(v), 10)) }
		fmt.Println(pos, code.Definitions[op].Name, strings.Join(p, " "))
	} else {
		fmt.Println(pos, code.Definitions[op].Name)
	}
	return pos
}

func (c *Compiler) replaceInstr(pos int, new []byte) {
	for i := 0; i < len(new); i++ {
		c.instrs[pos+i] = new[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instrs[opPos])
	new := code.Make(op, operand)
	c.replaceInstr(opPos, new)
}

func (c *Compiler) addInstr(ins []byte) int {
	posNewIns := len(c.instrs)
	c.instrs = append(c.instrs, ins...)
	return posNewIns
}