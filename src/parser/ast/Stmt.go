package ast

import "strings"

type Program struct {
	Stmt     `json:"-"`
	NodeType string
	Body     []Stmt
}

func (n Program) GetNodeType() string { return n.NodeType }
func (n Program) ToString() string {
	str := new(strings.Builder)
	for i := 0; i < len(n.Body); i++ {
		str.WriteString(n.Body[i].ToString())
		if i != (len(n.Body)-1) {
			str.WriteString(";")
		}
	}
	return str.String()
}

type ExprStmt struct {
	Stmt       `json:"-"`
	NodeType   string
	Expression Expr
}

func (n ExprStmt) GetNodeType() string { return n.NodeType }
func (n ExprStmt) ToString() string { return n.Expression.ToString() }


type ReturnStmt struct {
	Stmt       `json:"-"`
	NodeType   string
	Expression Expr
}

func (n ReturnStmt) GetNodeType() string { return n.NodeType }
func (n ReturnStmt) ToString() string { return "return "+n.Expression.ToString() }

type FuncStmt struct {
	Stmt     `json:"-"`
	NodeType string
	Name     string
	Args     []Expr
	Body     []Stmt
}

func (n FuncStmt) GetNodeType() string { return n.NodeType }
func (n FuncStmt) ToString() string {
	str := new(strings.Builder)
	str.WriteString("function ")
	str.WriteString(n.Name)
	str.WriteString("(")
	for i := 0; i < len(n.Args); i++ {
		str.WriteString(n.Args[i].ToString())
		if i != (len(n.Args)-1) {
			str.WriteString(",")
		}
	}
	str.WriteString("){")
	for _, v := range n.Body {
		str.WriteString(v.ToString())
		str.WriteString(";")
	}
	str.WriteString("}")
	return str.String()
}

type VarStmt struct {
	Stmt     `json:"-"`
	NodeType string
	Name     string
	Value    Expr
	IsConst  bool
}

func (n VarStmt) GetNodeType() string { return n.NodeType }
func (n VarStmt) ToString() string {
	str := new(strings.Builder)
	if n.IsConst { str.WriteString("const ") } else { str.WriteString("var ") }
	str.WriteString(n.Name)
	str.WriteString("=")
	str.WriteString(n.Value.ToString())
	return str.String()
}

type VarStmtList struct {
	Stmt     `json:"-"`
	NodeType string
	Stmts    []Stmt
}

func (n VarStmtList) GetNodeType() string { return n.NodeType }
func (n VarStmtList) ToString() string {
	str := new(strings.Builder)
	for i := 0; i < len(n.Stmts); i++ {
		str.WriteString(n.Stmts[i].ToString())
		if i != (len(n.Stmts)-1) {
			str.WriteString(";")
		}
	}
	return str.String()
}
type BodyStmt struct {
	Stmt     `json:"-"`
	NodeType string
	Body     []Stmt
}

func (n BodyStmt) GetNodeType() string { return n.NodeType }
func (n BodyStmt) ToString() string {
	str := new(strings.Builder)
	str.WriteString("{")
	for i := 0; i < len(n.Body); i++ {
		str.WriteString(n.Body[i].ToString())
		if i != (len(n.Body)-1) {
			str.WriteString(";")
		}
	}
	str.WriteString("}")
	return str.String()
}

type ForStmt struct {
	Stmt      `json:"-"`
	NodeType  string
	Init      Expr
	Condition Expr
	After     Expr
	Statement Stmt
}

func (n ForStmt) GetNodeType() string { return n.NodeType }
func (n ForStmt) ToString() string {
	str := new(strings.Builder)
	str.WriteString("for (")
	if n.Init != nil { str.WriteString(n.Init.ToString()); str.WriteString(";"); } else { str.WriteString(";") }
	if n.Condition != nil { str.WriteString(n.Condition.ToString()); str.WriteString(";"); } else { str.WriteString(";") }
	if n.After != nil { str.WriteString(n.After.ToString()) }
	str.WriteString(")")
	str.WriteString(n.Statement.ToString())
	return str.String()
}

type ForInStmt struct {
	Stmt      `json:"-"`
	NodeType  string
	Var      string
	IsConst bool
	Object Expr
	Statement Stmt
}

func (n ForInStmt) GetNodeType() string { return n.NodeType }
func (n ForInStmt) ToString() string {
	str := new(strings.Builder)
	str.WriteString("for (")
	if n.IsConst { str.WriteString("const ") } else { str.WriteString("var ") }
	str.WriteString(n.Var); str.WriteString(" in "); str.WriteString(n.Object.ToString())
	str.WriteString(")")
	str.WriteString(n.Statement.ToString())
	return str.String()
}

type IfStmt struct {
	Stmt       `json:"-"`
	NodeType   string
	Condition  Expr
	Consequent Expr
	Other      Expr
}

func (n IfStmt) GetNodeType() string { return n.NodeType }
func (n IfStmt) ToString() string {
	str := new(strings.Builder)
	str.WriteString("if (")
	str.WriteString(n.Condition.ToString())
	str.WriteString(")")
	str.WriteString(n.Consequent.ToString())
	if (n.Consequent.GetNodeType() != "BodyStmt") { str.WriteString(" ;") }
	str.WriteString("else ")
	str.WriteString(n.Other.ToString())
	return str.String()
}

type WhileStmt struct {
	Stmt       `json:"-"`
	NodeType   string
	Condition  Expr
	Body Stmt
}

func (n WhileStmt) GetNodeType() string { return n.NodeType }
func (n WhileStmt) ToString() string {
	str := new(strings.Builder)
	str.WriteString("while (")
	str.WriteString(n.Condition.ToString())
	str.WriteString(")")
	str.WriteString(n.Body.ToString())
	return str.String()
}

type ContinueStmt struct {
	Stmt       `json:"-"`
	NodeType   string
}

func (n ContinueStmt) GetNodeType() string { return n.NodeType }
func (n ContinueStmt) ToString() string {
	str := new(strings.Builder)
	str.WriteString("continue")
	return str.String()
}

type BreakStmt struct {
	Stmt       `json:"-"`
	NodeType   string
}

func (n BreakStmt) GetNodeType() string { return n.NodeType }
func (n BreakStmt) ToString() string {
	str := new(strings.Builder)
	str.WriteString("break")
	return str.String()
}
