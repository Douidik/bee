package main

type Def interface {
	Id() string
}

type Var struct {
	Name string
	Type Type
}

type Typedef struct {
	Name string
	Type Type
}

type Fn struct {
	Name   string
	Return Var
	Params []Var
}

func (v Var) Id() string {
	return v.Name
}

func (fn Fn) Id() string {
	return fn.Name
}

func (td Typedef) Id() string {
	return td.Name
}
