package main

import (
	"io"
	"math"
	"strings"
)

type CompoundStmt struct {
	Body  []Node
	Depth uint
}

type IfStmt struct {
	Keyword Token
	Cond    []Node
	Stmt    CompoundStmt
}

type ForStmt struct {
	Keyword Token
	Cond    []Node
	Stmt    CompoundStmt
}

type FuncStmt struct {
	Func     Func
	Compound CompoundStmt
}

type DeclareStmt struct {
	Def Def
}

type DefineStmt struct {
	Def  Def
	Expr Node
}

type ExprStmt struct {
	Expr Node
}

const IntensityFac = 0.1

func shade(depth uint, color uint32) uint32 {
	shaded := uint32(0x00)

	for off := 0; off < 6; off += 2 {
		channel := ((color >> off) & 0xff) * IntensityFac * depth
		if channel > 0xff {
			channel = 0xff
		}
		shaded |= channel << off
	}

	return shaded
}

func (c *CompoundStmt) Graph(w io.Writer) {
	Writef(w, `subgraph cluster_%p {`, c)
	color := shade(c.Depth, 0xFEF3BD)
	desc := strings.NewBuilder()
	desc.WriteString(``)
	// # FFF9C4, each channel is has an inten

	for n, _ := range c.Body {

	}

	for n, stmt := range c.Body {

		stmt.Graph(w)
	}

	Writef(w, `}`)
}
