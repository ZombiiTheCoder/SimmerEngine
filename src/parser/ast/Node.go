package ast

type Node interface {
	GetNodeType() string
	ToString() string
}

type Expr interface {
	Node
}

type Stmt interface {
	Node
}

type StmtExpr interface {
	Stmt
	Expr
}
