package main

import (
	"io"
	"fmt"
)

type Asm_x86 struct {}
type Asm_6502 struct {}

func Writef(w io.Writer, f string, args ...any) {
	str := fmt.Sprintf(f, args...)
	w.Write([]byte(str))
	w.Write([]byte{'\n'})
}

type Node interface {
	Graph(w io.Writer)
// 	Asm_x86(ctx *Asm_x86)
// 	Asm_6502(ctx *Asm_6502)
}
