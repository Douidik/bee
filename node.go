package main

import (
	"fmt"
	"io"
)

type Asm_x86 struct{}
type Asm_6502 struct{}

func Writef(w io.Writer, f string, args ...any) {
	str := fmt.Sprintf(f, args...)
	w.Write([]byte(str))
	w.Write([]byte{'\n'})
}

type Node interface {
	// Dump(w io.Writer, depth uint)
	// Graph(w io.Writer)
	// 	Asm_x86(ctx *Asm_x86)
	// 	Asm_6502(ctx *Asm_6502)
}

type Order uint

const (
	OrderPrev Order = 0
	OrderPost Order = 1
)

type Def interface {
	Id() string
}

type Var struct {
	Name string
}

// Language type primitive, cannot be de-constructed into simpler types.
// Signedness and size are used to choose the correct CPU Instruction
type Atom struct {
	Name   string
	Size   uint
	Signed bool
}

type Fn struct {
	Name   string
	Return *Var
	Params []*Var
}

func (v *Var) Id() string {
	return v.Name
}

func (fn *Fn) Id() string {
	return fn.Name
}

func (atom *Atom) Id() string {
	return atom.Name
}

type Reference struct {
	Def Def
}

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
	Operand *Fn
	Args    []*Var
}

type Cast struct {
	Operand Node
}

type Nested struct {
	Body []Node
}

type IntExpr struct {
	Value  uint64
	Size   uint
	Signed bool
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

type If struct {
	Conds []Node
	Body  []Node
	Else  []Node
}
type For If

type DefineExpr struct {
	Def  Def
	Expr Node
}
type DeclareExpr DefineExpr
