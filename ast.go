package main

type Ast struct {
	Body  []Node
	Scope *Scope
}

type Scope struct {
	Sp    uint64
	Defs  map[string]Def
	Owner *Scope
}

func NewScope(owner *Scope) *Scope {
	return &Scope{Sp: 0, Defs: make(map[string]Def), Owner: owner}
}

func (sc *Scope) Search(id string) Def {
	if def, found := sc.Defs[id]; found {
		return def
	}
	if sc.Owner != nil {
		return sc.Owner.Search(id)
	}
	return nil
}

func (sc *Scope) Add(def Def) Def {
	if v, ok := def.(*Var); ok {
		v.Offset = sc.Sp
		sc.Sp += uint64(v.Type.Size())
	}

	sc.Defs[def.Id()] = def
	return def
}

func (ast *Ast) Asm_x86() Asm_x86 {
	asm := Asm_x86{Scope: ast.Scope}
	asm.Writef("section .text")
	asm.Writef("global _start")
	asm.Labelf("_start")
	for _, node := range ast.Body {
		node.Asm_x86(&asm)
	}
	asm.Writef("mov eax")
	asm.Writef("int 0x80")
	return asm
}
