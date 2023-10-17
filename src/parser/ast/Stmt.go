package ast

type Program struct {
	Stmt
	Body []Stmt
}

type ExprStmt struct {
	Stmt
	Expression Expr
}
