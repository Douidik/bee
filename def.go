package main

type DefType uint

const (
	VariableType DefType = iota
	FuncType
	StrConstantType
)

type Def interface {
	DefName() string
	DefType() DefType
}

type Variable struct {
	Type Def
	Name string
}

type Func struct {
	Name   string
	Return *Variable
	Args   []*Variable
}

// type StrConstant struct {
// 	Content string
// }

func (v *Variable) DefName() string {
	return v.Name
}

func (v *Variable) DefType() DefType {
	return VariableType
}

func (f *Func) DefName() string {
	return f.Name
}

func (f *Func) DefType() DefType {
	return FuncType
}

// func (s *StrConstant) DefName() string {
// 	return s.Content
// }

// func (v *StrConstant) DefType() DefType {
// 	return StrConstantType
// }
