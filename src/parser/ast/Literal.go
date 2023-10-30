package ast

import "fmt"

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
	Value    float64
}

func (n NumberLiteral) GetNodeType() string { return n.NodeType }
func (n NumberLiteral) ToString() string    { return fmt.Sprint(n.Value) }

type StringLiteral struct {
	Expr     `json:"-"`
	NodeType string
	Value    string
}

func (n StringLiteral) GetNodeType() string { return n.NodeType }
func (n StringLiteral) ToString() string    { return "\"" + n.Value + "\"" }

type NullLiteral struct {
	Expr     `json:"-"`
	NodeType string
	Value    interface{}
}

func (n NullLiteral) GetNodeType() string { return n.NodeType }
func (n NullLiteral) ToString() string    { return "nil" }

type BooleanLiteral struct {
	Expr     `json:"-"`
	NodeType string
	Value    bool
}

func (n BooleanLiteral) GetNodeType() string { return n.NodeType }
func (n BooleanLiteral) ToString() string    { return fmt.Sprint(n.Value) }
