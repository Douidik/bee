package main

import (
	"io"
)

type HeadStmt struct {
	Body []Node
}

type CompoundStmt struct {
	Body []Node
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

func (h *HeadStmt) Graph(w io.Writer) {
	Writef(w, `strict digraph {`)

	for _, stmt := range h.Body {
		stmt.Graph(w)
	}

	Writef(w, `}`)
}

func (c *CompoundStmt) Graph(w io.Writer) {
	Writef(w, `subgraph cluster_%p {`, c)
	
	for _, stmt := range c.Body {
		stmt.Graph(w)
	}

	Writef(w, `}`)
}
