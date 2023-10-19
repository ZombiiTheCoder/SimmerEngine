package ast

import "simmer_js_engine/parser/tokens"

type AssignExpr struct {
	Expr `json:"-"`
	NodeType string
	Left  Expr
	Right Expr
	Op    tokens.TokenType
}

type TernaryExpr struct {
	Expr `json:"-"`
	NodeType string
	Condition  Expr
	Consequent Expr
	Else Expr
}

type PrefixExpr struct {
	Expr `json:"-"`
	NodeType string
	Right  Expr
	Op    tokens.TokenType
}

type BinaryExpr struct {
	Expr `json:"-"`
	NodeType string
	Left  Expr
	Right Expr
	Op    tokens.TokenType
}

type PostfixExpr struct {
	Expr `json:"-"`
	NodeType string
	Left  Expr
	Op    tokens.TokenType
}

type CallExpr struct {
	Expr `json:"-"`
	NodeType string
	Caller  Expr
	Args []Expr
}