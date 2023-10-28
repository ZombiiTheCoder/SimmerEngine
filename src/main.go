package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"simmer_js_engine/parser"
	"simmer_js_engine/parser/ast"
	"simmer_js_engine/parser/lexer"
	"simmer_js_engine/visitor"
	opcode "simmer_js_engine/visitor/opcodes"
	"simmer_js_engine/vm"
	"strings"
)

func prod(text []byte, e error) ast.Stmt {
	lex := lexer.InitLexer(text, e)
	tokens := lex.Tokenize()
	par := parser.InitParser(tokens)
	return par.ProduceAst()
}

func main() {
	ast := prod(os.ReadFile(os.Args[1]))
	by, _ := json.MarshalIndent(ast, "", "	",)
	fmt.Println(string(by))
	os.WriteFile(strings.ReplaceAll(os.Args[1], ".js", "_ast.json"), by, fs.FileMode(os.O_CREATE))
	ast = prod([]byte(ast.ToString()), nil)
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "-compact":
			os.WriteFile(os.Args[3], []byte(ast.ToString()), fs.FileMode(os.O_CREATE))
		}
	}
	fmt.Println("-------------------------------")
	v := visitor.InitVisitor()
	v.Visit(ast)
	v.Emit(opcode.OP_RETURN)
	fmt.Println("-------------------------------")
	v.Chunk.Disassemble("Program")
	fmt.Println("-------------------------------")
	vm := vm.InitVM()
	vm.Interpret(v.Chunk)
}