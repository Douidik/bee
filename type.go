package main

import "golang.org/x/exp/slices"

type Type interface {
	Size() uint
	Infers(to Type) bool
}

// Language type primitive, cannot be de-constructed into simpler types.
// Signedness and size are used to choose the correct CPU Instruction
type Atom struct {
	size   uint
	signed bool
	float  bool
}

type Void struct{}

type Struct struct {
	Members []Var
}

func (at Atom) Size() uint {
	return at.size
}

func (at Atom) Infers(to Type) bool {
	_, same := to.(Atom)
	return same
}

func (s Struct) Size() uint {
	size := uint(0)
	for _, member := range s.Members {
		size += member.Type.Size()
	}
	return size
}

func (s Struct) Infers(to Type) bool {
	if sto, same := to.(Struct); same {
		equals := func(a, b Var) bool {
			return a.Type == b.Type && a.Name == b.Name
		}
		return slices.EqualFunc(s.Members, sto.Members, equals)
	}
	return false
}

func (v Void) Size() uint {
	return 0
}

func (v Void) Infers(to Type) bool {
	_, same := to.(Void)
	return same
}
