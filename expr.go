package main

type UnaryExpr struct {
	Order    Order
	Operand  Node
	Operator Token
}

type BinaryExpr struct {
	Operands [2]Node
	Operator Token
}

type IndexExpr struct {
	Operand Node
}

type InvokeExpr struct {
	Operand Fn
	Args    []Var
}

type DefineExpr struct {
	Def  Def
	Expr Node
}
type DeclareExpr DefineExpr

func (un UnaryExpr) Result() Type {
	return un.Operand.Result()
}

func (bin BinaryExpr) Result() Type {
	return bin.Operands[0].Result()
}

func (ind IndexExpr) Result() Type {
	return ind.Operand.Result()
}

func (inv InvokeExpr) Result() Type {
	return inv.Operand.Return.Type
}

func (decl DeclareExpr) Result() Type {
	return decl.Expr.Result()
}

func (def DefineExpr) Result() Type {
	return def.Expr.Result()
}
