package main

import (
	"encoding/json"
	"fmt"
	"os"
	"simmer_js_engine/parser"
	"simmer_js_engine/parser/lexer"
)

func main() {
	lex := lexer.InitLexer(os.ReadFile(os.Args[1]))
	par := parser.InitParser(lex.Tokenize())
	by, _ := json.MarshalIndent(par.ProduceAst(), "", "	",)
	fmt.Println(string(by))
}