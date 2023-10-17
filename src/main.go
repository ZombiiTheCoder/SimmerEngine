package main

import (
	"fmt"
	"os"
	"simmer_js_engine/parser"
	"simmer_js_engine/parser/lexer"
)

func main() {
	lex := lexer.InitLexer(os.ReadFile(os.Args[1]))
	par := parser.InitParser(lex.Tokenize())
	fmt.Println(par.ProduceAst())
}