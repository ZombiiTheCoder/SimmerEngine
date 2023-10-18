package ast

type NumberLiteral struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}