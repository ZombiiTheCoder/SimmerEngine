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
	tokens := lex.Tokenize()
	par := parser.InitParser(tokens)
	// for _, v := range tokens {
	// 	fmt.Println(v)
	// }
	by, _ := json.MarshalIndent(par.ProduceAst(), "", "	",)
	fmt.Println(string(by))
}