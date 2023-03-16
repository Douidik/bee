package main

type Node interface {
	Result() Type

	// Dump(w io.Writer, depth uint)
	// Graph(w io.Writer)
	Asm_x86(asm *Asm_x86)
	// Asm_6502(asm *Asm_6502)
}

func Precedence(n Node) uint {
	switch node := n.(type) {
	case UnaryExpr:
		switch node.Order {
		case OrderPrev:
			return 0
		case OrderPost:
			return 1
		}
	case InvokeExpr, IndexExpr:
		return 0
	case BinaryExpr:
		switch n.(BinaryExpr).Operator.Trait {
		case Mul, Div, Mod, BinShiftL, BinShiftR, BinAnd:
			return 1
		case Add, Sub, BinOr, BinXor, BinNot:
			return 2
		case Equal, NotEq, Less, LessEq, Greater, GreaterEq:
			return 3
		case KwAnd:
			return 4
		case KwOr:
			return 5
		}
	default:
		return 1
	}
	return 0
}

type Order uint

const (
	OrderPrev Order = 0
	OrderPost Order = 1
)

type Reference struct {
	Def Def
}

type Cast struct {
	Operand Node
	Type    Type
}

type Nest struct {
	Body []Node
}

type Compound struct {
	Scope *Scope
	Body  []Node
}

type If struct {
	Conds Compound
	If    Compound
	Else  Compound
}

type For struct {
	Conds Compound
	Body  Compound
}

func (ref Reference) Result() Type {
	switch def := ref.Def.(type) {
	case *Typedef:
	case *Var:
		return def.Type
	default:
		return Void{}
	}
	return Void{}
}

func (cast Cast) Result() Type {
	return cast.Type
}

func (nest Nest) Result() Type {
	if len(nest.Body) != 0 {
		return nest.Body[len(nest.Body)-1].Result()
	}
	return Void{}
}

func (i If) Result() Type {
	return i.If.Result()
}

func (f For) Result() Type {
	return f.Body.Result()
}

func (comp Compound) Result() Type {
	if len(comp.Body) != 0 {
		return comp.Body[len(comp.Body)-1].Result()
	}
	return Void{}
}

func (ref Reference) Asm_x86(asm *Asm_x86) {

}

func (cast Cast) Asm_x86(asm *Asm_x86) {
}

func (nest Nest) Asm_x86(asm *Asm_x86) {
}

func (i If) Asm_x86(asm *Asm_x86) {
}

func (f For) Asm_x86(asm *Asm_x86) {
}

func (comp Compound) Asm_x86(asm *Asm_x86) {
	asm.Scope = comp.Scope
	asm.Writef("push rbp")
	asm.Writef("mov rbp, rsp")
	for _, node := range comp.Body {
		node.Asm_x86(asm)
	}
	asm.Labelf("L%d", asm.PushLabel())
	asm.Writef("pop rbp")
	asm.Scope = comp.Scope.Owner
}
