package main

import (
	"fmt"
	"strings"
)

type Asm_x86 struct {
	Stream strings.Builder
	Scope  *Scope
	label  uint32
}
type Asm_6502 Asm_x86

func (asm *Asm_x86) PushLabel() uint32 {
	label := asm.label
	asm.label += 1
	return label
}

func (asm *Asm_x86) Labelf(f string, args ...interface{}) string {
	name := fmt.Sprintf(f, args...)
	asm.Stream.WriteString(fmt.Sprintf("%s:\n", name))
	return name
}

func (asm *Asm_x86) Writef(f string, args ...interface{}) {
	asm.Stream.WriteByte('\t')
	asm.Stream.WriteString(fmt.Sprintf(f, args...))
	asm.Stream.WriteByte('\n')
}
