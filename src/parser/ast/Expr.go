package ast

import (
	"simmer_js_engine/parser/tokens"
	"strings"
)

type AssignExpr struct {
	Expr `json:"-"`
	NodeType string
	Left  Expr
	Right Expr
	Op    tokens.TokenType
}
func (n AssignExpr) GetNodeType() string { return n.NodeType }
func (n AssignExpr) ToString() string { return n.Left.ToString() + " " + string(n.Op) + " " + n.Right.ToString()}


type TernaryExpr struct {
	Expr `json:"-"`
	NodeType string
	Condition  Expr
	Consequent Expr
	Other Expr
}
func (n TernaryExpr) GetNodeType() string { return n.NodeType }
func (n TernaryExpr) ToString() string { return n.Condition.ToString() + "?" + n.Consequent.ToString() + ":" + n.Other.ToString()}

type PrefixExpr struct {
	Expr `json:"-"`
	NodeType string
	Right  Expr
	Op    tokens.TokenType
}
func (n PrefixExpr) GetNodeType() string { return n.NodeType }
func (n PrefixExpr) ToString() string { return string(n.Op) + n.Right.ToString()}

type BinaryExpr struct {
	Expr `json:"-"`
	NodeType string
	Left  Expr
	Right Expr
	Op    tokens.TokenType
}
func (n BinaryExpr) GetNodeType() string { return n.NodeType }
func (n BinaryExpr) ToString() string { return "("+n.Left.ToString() + " " + string(n.Op) + " " + n.Right.ToString()+")"}

type PostfixExpr struct {
	Expr `json:"-"`
	NodeType string
	Left  Expr
	Op    tokens.TokenType
}
func (n PostfixExpr) GetNodeType() string { return n.NodeType }
func (n PostfixExpr) ToString() string { return n.Left.ToString() + string(n.Op)}

type CallExpr struct {
	Expr `json:"-"`
	NodeType string
	Caller  Expr
	Args []Expr
}
func (n CallExpr) GetNodeType() string { return n.NodeType }
func (n CallExpr) ToString() string {
	str := new(strings.Builder)
	str.WriteString(n.Caller.ToString())
	str.WriteString("(")
	for i := 0; i < len(n.Args); i++ {
		str.WriteString(n.Args[i].ToString())
		if i != (len(n.Args)-1) {
			str.WriteString(",")
		}
	}
	str.WriteString(")")
	return str.String()
}