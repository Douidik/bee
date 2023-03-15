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

func Precedence(n Node) uint {
	switch n.(type) {
	case UnaryExpr:
		switch n.(UnaryExpr).Order {
		case OrderPrev:
			return 0
		case OrderPost:
			return 1
		}
	case InvokeExpr, IndexExpr:
		return 0
	case BinaryExpr:
		switch n.(BinaryExpr).Operator.Trait {
		case Mul, Div, Mod, BinShiftL, BinShiftR, BinAnd:
			return 1
		case Add, Sub, BinOr, BinXor, BinNot:
			return 2
		case Equal, NotEq, Less, LessEq, Greater, GreaterEq:
			return 3
		case KwAnd:
			return 4
		case KwOr:
			return 5
		}
	}
}

type Order uint

const (
	OrderPrev Order = 0
	OrderPost Order = 1
)

type Reference struct {
	Def Def
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

type If struct {
	Conds Compound
	If    Compound
	Else  Compound
}

type For struct {
	Conds Compound
	Body  Compound
}

func (ref Reference) Result() Type {
	switch def := ref.Def.(type) {
	case *Typedef:
		return def.Type
	case *Var:
		return def.Type
	default:
		return Void{}
	}
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

func (i If) Result() Type {
	return i.If.Result()
}

func (f For) Result() Type {
	return f.Body.Result()
}

func (comp Compound) Result() Type {
	if len(comp.Body) != 0 {
		return comp.Body[len(comp.Body)-1].Result()
	}
	return Void{}
}
