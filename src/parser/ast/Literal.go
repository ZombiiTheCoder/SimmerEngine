package ast

type Identifer struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}

type NumberLiteral struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}

type StringLiteral struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}