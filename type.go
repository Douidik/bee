package main

import "golang.org/x/exp/slices"

type Type interface {
	Size() uint64
	Cast(as Type) bool
}

// Language type primitive, cannot be de-constructed into simpler types.
// Signedness and size are used to choose the correct CPU Instruction
type Atom struct {
	size   uint64
	signed bool
	float  bool
}

type Void struct{}

type Struct struct {
	Members []Var
}

func (at Atom) Size() uint64 {
	return at.size
}

func (at Atom) Cast(as Type) bool {
	_, same := as.(Atom)
	return same
}

func (s Struct) Size() uint64 {
	size := uint64(0)
	for _, member := range s.Members {
		size += member.Type.Size()
	}
	return size
}

func (s Struct) Cast(as Type) bool {
	if as, same := as.(Struct); same {
		eq := func(a, b Var) bool {
			return a.Type == b.Type && a.Name == b.Name
		}
		return slices.EqualFunc(s.Members, as.Members, eq)
	}
	return false
}

func (v Void) Size() uint64 {
	return 0
}

func (v Void) Cast(as Type) bool {
	_, same := as.(Void)
	return same
}
