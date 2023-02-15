package main

type Ast struct {
	head   HeadStmt
	scopes []Scope
}

type Scope struct {
	Defs  []Def
	Owner *Scope
}

func (ast *Ast) NewScope(owner *Scope) *Scope {
	ast.scopes = append(ast.scopes, Scope{Owner: owner})
	return &ast.scopes[len(ast.scopes)-1]
}

func (sc *Scope) Search(name string) Def {
	for _, def := range sc.Defs {
		if def.DefName() == name {
			return def
		}
	}

	if sc.Owner != nil {
		return sc.Owner.Search(name)
	}
	return nil
}

func (sc *Scope) Define(def Def) {
	sc.Defs = append(sc.Defs, def)
}
