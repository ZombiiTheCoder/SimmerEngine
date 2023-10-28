package ast

type Identifier struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}

func (n Identifier) GetNodeType() string { return n.NodeType }
func (n Identifier) ToString() string    { return n.Value }

type Literal struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}

func (n Literal) GetNodeType() string { return n.NodeType }
func (n Literal) ToString() string    { return n.Value }