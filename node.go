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
	Result() Type
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
	Type Type
}

type Typedef struct {
	Name string
	Type Type
}

type Fn struct {
	Name   string
	Return *Var
	Params []*Var
}

func (v Var) Id() string {
	return v.Name
}

func (fn Fn) Id() string {
	return fn.Name
}

func (td Typedef) Id() string {
	return td.Name
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
	Type    Type
}

type Nest struct {
	Body []Node
}

type Compound struct {
	Scope *Scope
	Body  []Node
}

type IntExpr struct {
	Value uint64
	Type  Atom
}

type FloatExpr struct {
	Value float64
	Type  Atom
}

type StrExpr struct {
	Value string
}

type CharExpr struct {
	Value byte
}

type If struct {
	Conds Compound
	If    Compound
	Else  Compound
}

type For struct {
	Conds Compound
	Body  Compound
}

type DefineExpr struct {
	Def  Def
	Expr Node
}
type DeclareExpr DefineExpr

func (ref Reference) Result() Type {
	switch def := ref.Def.(type) {
	case Typedef:
		return def.Type
	case Var:
		return def.Type
	default:
		return Void{}
	}
}

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

func (cast Cast) Result() Type {
	return cast.Type
}

func (nest Nest) Result() Type {
	if len(nest.Body) != 0 {
		return nest.Body[len(nest.Body)-1].Result()
	}
	return Void{}
}

func (int IntExpr) Result() Type {
	return int.Type
}

func (fl FloatExpr) Result() Type {
	return fl.Type
}

func (str StrExpr) Result() Type {
	return Void{}
}

func (char CharExpr) Result() Type {
	return Atom{size: 8, signed: true}
}

func (i If) Result() Type {
	return i.If.Result()
}

func (f For) Result() Type {
	return f.Body.Result()
}

func (cd Compound) Result() Type {
	if len(cd.Body) != 0 {
		return cd.Body[len(cd.Body)-1].Result()
	}
	return Void{}
}

func (decl DeclareExpr) Result() Type {
	return decl.Expr.Result()
}

func (def DefineExpr) Result() Type {
	return def.Expr.Result()
}
