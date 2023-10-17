package ast

import "simmer_js_engine/parser/tokens"

type AssignExpr struct {
	Expr
	Left  Expr
	Right Expr
	Op    tokens.TokenType
}

type BinaryExpr struct {
	Expr
	Left  Expr
	Right Expr
	Op    tokens.TokenType
}

type CallExpr struct {
	Expr
	Caller  Expr
	Args []Expr
}