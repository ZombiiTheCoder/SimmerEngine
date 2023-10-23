package ast

type Identifier struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}

func (n Identifier) GetNodeType() string { return n.NodeType }
func (n Identifier) ToString() string    { return n.Value }

type NumberLiteral struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}

func (n NumberLiteral) GetNodeType() string { return n.NodeType }
func (n NumberLiteral) ToString() string    { return n.Value }

type StringLiteral struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}

func (n StringLiteral) GetNodeType() string { return n.NodeType }
func (n StringLiteral) ToString() string    { return n.Value }
