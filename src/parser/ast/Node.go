package ast

type Node interface{}

type Expr interface{ Node }

type Stmt interface{ Node }