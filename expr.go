package main

type Order uint

const (
	OrderPrev Order = 0
	OrderPost Order = 1
)

type IdExpr struct {
	Def Def
}

type UnaryExpr struct {
	Order Order
	Expr  Node
	Op    Token
}

type BinaryExpr struct {
	Lhs Node
	Rhs Node
	Op  Token
}

type IndexExpr struct {
	Expr Node
}

type InvokeExpr struct {
	Func *Func
	Args []*Variable
}

type CastExpr struct {
	Type *Def
	Expr Node
}

type NestedExpr struct {
	Body []Node
}

type IntExpr struct {
	Value int64
	Size  uint
}

type FloatExpr struct {
	Value float64
	Size  uint
}

type StrExpr struct {
	Value string
}

type CharExpr struct {
	Value byte
}
