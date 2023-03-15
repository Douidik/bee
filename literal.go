package main

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
	return Atom{size: 1, signed: true}
}

func (int IntExpr) Asm_x86(asm *Asm_x86) {
	asm.Writef("push %x")
}
