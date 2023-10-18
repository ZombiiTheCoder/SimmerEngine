package ast

type Program struct {
	Stmt     `json:"-"`
	NodeType string
	Body     []Stmt
}

type ExprStmt struct {
	Stmt       `json:"-"`
	NodeType   string
	Expression Expr
}
