package main

type Ast struct {
	Body   []Node
	scopes []Scope
}

type Scope struct {
	Defs  map[string]Def
	Owner *Scope
}

func (ast *Ast) NewScope(owner *Scope) *Scope {
	ast.scopes = append(ast.scopes, Scope{Defs: make(map[string]Def), Owner: owner})
	return &ast.scopes[len(ast.scopes)-1]
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

func (sc *Scope) Add(def Def) {
	sc.Defs[def.Id()] = def
}
