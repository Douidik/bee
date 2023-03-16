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

func (un UnaryExpr) Asm_x86(asm *Asm_x86) {

}

func (bin BinaryExpr) Asm_x86(asm *Asm_x86) {
	bin.Operands[0].Asm_x86(asm)
	bin.Operands[1].Asm_x86(asm)

	asm.Writef("pop rbx")
	asm.Writef("pop rax")

	switch bin.Operator.Trait {
	case Sub:
		asm.Writef("sub rax, rbx")
	case Add:
		asm.Writef("add rax, rbx")
	default:
		panic("todo!")
	}

	asm.Writef("push rax")
}

func (ind IndexExpr) Asm_x86(asm *Asm_x86) {
}

func (inv InvokeExpr) Asm_x86(asm *Asm_x86) {
}

func (decl DeclareExpr) Asm_x86(asm *Asm_x86) {
}

func (def DefineExpr) Asm_x86(asm *Asm_x86) {
	switch d := def.Def.(type) {
	case Var:
		def.Expr.Asm_x86(asm)
		asm.Writef("pop rax")
		asm.Writef("mov dword ptr [rbp - %d], eax", d.Offset)
	case Fn:
		panic("todo!")
	}
}
