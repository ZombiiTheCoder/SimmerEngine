package main

import (
	"fmt"
	"os"
	"simmer_js_engine/parser"
	"simmer_js_engine/parser/lexer"
)

func main() {
	lex := lexer.InitLexer(os.ReadFile("F:/Projects/SimmerJsEngine/tests/numbers.js"))
	par := parser.InitParser(lex.Tokenize())
	fmt.Println(par.ProduceAst())
}