package main

import (
	"fmt"
	"golang.org/x/exp/slices"
	"math"
	"strconv"
)

type Parser struct {
	sn        Scanner
	peekQueue []Token
	ast       Ast
	scope     *Scope
}

func NewParser(sc Scanner) Parser {
	return Parser{sn: sc, peekQueue: make([]Token, 8)}
}

func (ps *Parser) Parse() (*Ast, error) {
	ps.scope = ps.ast.NewScope(nil)
	ps.scope.Add(&Atom{Name: "bool", Size: 8, Signed: false})
	ps.scope.Add(&Atom{Name: "char", Size: 8, Signed: true})
	ps.scope.Add(&Atom{Name: "s8", Size: 8, Signed: true})
	ps.scope.Add(&Atom{Name: "s16", Size: 16, Signed: true})
	ps.scope.Add(&Atom{Name: "s32", Size: 32, Signed: true})
	ps.scope.Add(&Atom{Name: "s64", Size: 64, Signed: true})
	ps.scope.Add(&Atom{Name: "u8", Size: 8, Signed: false})
	ps.scope.Add(&Atom{Name: "u16", Size: 16, Signed: false})
	ps.scope.Add(&Atom{Name: "u32", Size: 32, Signed: false})
	ps.scope.Add(&Atom{Name: "u64", Size: 64, Signed: false})

	for !ps.sn.Finished() {
		node, err := ps.parseNode(NewLine, Semicolon)
		if err != nil {
			return nil, err
		}
		ps.ast.Body = append(ps.ast.Body, node)
	}
	return &ps.ast, nil
}

func (ps *Parser) parseNode(delim ...Trait) (Node, error) {
	body := make([]Node, 8)
	var prev Node = nil
	var end Token

	for !end.Ok && !ps.sn.Finished() {
		node, err := ps.expectNode(prev)
		if err != nil {
			return nil, err
		}
		if node != nil {
			body = append(body, node)
			prev = node
		} else {
			break
		}
	}

	if !end.Ok {
		expected := make([]byte, 0)
		for i := 0; i < len(delim); i++ {
			if i < len(delim)-1 {
				expected = fmt.Appendf(expected, "%s, ", BeeTraitName(delim[i]))
			} else {
				expected = fmt.Appendf(expected, "%s", BeeTraitName(delim[i]))
			}
		}
		ps.errorf(end, "Expected <%s>", expected)
	}

	switch len(body) {
	case 0:
		return nil, nil
	case 1:
		return body[0], nil
	default:
		return Nested{body}, nil
	}
}

func (ps *Parser) expectNode(prev Node) (Node, error) {
	for ps.token(NewLine).Ok {
		// allow new lines during if node not finished to parse
	}

	if id := ps.token(Identifier); id.Ok {
		def := ps.scope.Search(id.Expr)
		if def != nil {
			return Reference{def}, nil
		}
		if operator := ps.token(Define, Define); operator.Ok {
			expr, err := ps.expectNode(nil)
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, ps.errorf(operator, "No value given to '%s'", id.Expr)
			}
			def = ps.scope.Add()
			switch operator.Trait {
			case Define:
				return DefineExpr{
			}
		}
		return nil, ps.errorf(id, "Use of undeclared identifier")
	}

	if str := ps.token(RawStr, Str); str.Ok {
		var content string
		switch str.Trait {
		case RawStr:
			content = str.Expr[1 : len(str.Expr)-1]
		case Str:
			content = UnescapeStr(str.Expr[1 : len(str.Expr)-1])
		}
		return StrExpr{content}, nil
	}

	// Integers and floats constants infers to 32-64 bits depending on the value size
	if int := ps.token(IntDec, IntBin, IntHex); int.Ok {
		var (
			n    int64
			err  error
			size uint
		)

		switch int.Trait {
		case IntDec:
			n, err = strconv.ParseInt(int.Expr[:], 10, 64)
		case IntBin:
			n, err = strconv.ParseInt(int.Expr[2:], 2, 64)
		case IntHex:
			n, err = strconv.ParseInt(int.Expr[2:], 16, 64)
		}

		if n > math.MaxInt32 {
			size = 64
		} else {
			size = 32
		}
		return IntExpr{uint64(n), size, true}, err
	}

	if float := ps.token(Float); float.Ok {
		f, err := strconv.ParseFloat(float.Expr, 64)

		var size uint
		if f > math.MaxFloat32 {
			size = 64
		} else {
			size = 32
		}
		return FloatExpr{f, size}, err
	}

	if char := ps.token(Char); char.Ok {
		content := UnescapeStr(char.Expr[1 : len(char.Expr)-1])
		switch len(content) {
		case 0:
			return nil, ps.errorf(char, `Empty character constant`)
		case 1:
			return CharExpr{content[0]}, nil
		default:
			return nil, ps.errorf(char, `Character constant too long`)
		}
	}

	if sign := ps.token(Add, Sub); sign.Ok && prev == nil {
		expr, err := ps.expectNode(nil)
		if err != nil {
			return nil, err
		}
		return UnaryExpr{OrderPrev, expr, sign}, nil
	}

	if bin := ps.token(
		Assign,
		KwAnd, KwOr,
		Add, Sub, Mul, Div, Mod,
		BinNot, BinAnd, BinOr, BinXor, BinShiftL, BinShiftR,
		Equal, NotEq, Less, Greater, LessEq, GreaterEq); bin.Ok {
		if prev == nil {
			return nil, ps.errorf(bin, `Missing pre-operand for binary expression`)
		}
		next, err := ps.expectNode(nil)
		if err != nil {
			return nil, err
		}
		if next == nil {
			return nil, ps.errorf(bin, `Missing post-operand for binary expression`)
		}
		return BinaryExpr{[2]Node{prev, next}, bin}, nil
	}

	if incr := ps.token(Increment, Decrement); incr.Ok {
		if prev != nil {
			return UnaryExpr{OrderPost, prev, incr}, nil
		} else {
			next, err := ps.expectNode(nil)
			if err != nil {
				return nil, err
			}
			if next == nil {
				return nil, ps.errorf(incr, `Missing expression for increment`)
			}
			return UnaryExpr{OrderPrev, next, incr}, nil
		}
	}

	if paren := ps.token(ParenBegin); paren.Ok {
		node, err := ps.parseNode(ParenEnd)
		if err != nil {
			return nil, err
		}
		switch node.(type) {
		case nil:
			return nil, nil
		case Nested:
			return node.(Nested), err
		default:
			return Nested{Body: []Node{node}}, err
		}
	}
	return nil, nil
}

// func (ps *Parser) parseStmt() (Node, error) {
// 	switch peek := ps.peekTok(); peek.Trait {
// 	case KwIf:
// 	case KwFor:
// 	case KwReturn:
// 		break

// 	case Identifier:
// 		id := ps.expectTok(Identifier)
// 		op := ps.maybeTok(Define, Declare)
// 		if !op.Ok {
// 			break
// 		}

// 		if ps.maybeTok(KwFn).Ok {

// 		}

// 		if ps.maybeTok(Define).Ok {
// 			return DefineStmt{ps.scope.Add(id), expr}, err
// 		} else if ps.maybeTok(Declare).Ok {

// 		}
// 	}

// 	expr, err := ps.parseExpr(NewLine)
// 	return ExprStmt{expr}, err

// 	// switch peek := p.peekTok(); peek.Trait {
// 	// case Identifier:
// 	// 	name := p.expectTok(Identifier)
// 	// 	op := p.expectTok(Define, Declare)

// 	// 	if !op.Ok {
// 	// 		return nil, p.errorf(op, `Expected operator '::' or ':' after identifier`)
// 	// 	}

// 	// 	// def, err := p.parseDef(name)
// 	// 	// if err != nil {
// 	// 	// 	return nil, err
// 	// 	// }

// 	// 	switch op.Trait {
// 	// 	case Declare:

// 	// 	case Define:

// 	// 	}
// 	// }

// 	return nil, nil
// }

// // func (p *Parser) parseDef(name Token) (Node, error) {
// // 	if tok := p.maybeTok(ParenBegin); tok.Ok {
// // 	}

// // 	if tok := p.maybeTok(Identifier); tok.Ok {
// // 		def := p.scope.Search(tok.Expr)
// // 		if def == nil {
// // 			return nil, p.errorf(tok, `Unknown name`)
// // 		}

// // 	}
// // }

// func (ps *Parser) parseExpr(stop uint) (Node, error) {
// 	nested := NestedExpr{}
// 	var prev Node = nil
// 	var end Token

// 	for {
// 		if ps.maybeTok(stop, End).Ok {
// 			break
// 		}

// 		expr, err := ps.parseNextExpr(prev)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if expr != nil {
// 			nested.Body = append(nested.Body, expr)
// 		} else {
// 			break
// 		}
// 		prev = expr
// 	}

// 	if end.Trait != stop {
// 		return nil, ps.errorf(end, `Unexpected end of source during expression parsing`)
// 	}

// 	switch len(nested.Body) {
// 	case 0:
// 		return nil, nil
// 	case 1:
// 		return nested.Body[0], nil
// 	default:
// 		return nested, nil
// 	}
// }

// func (ps *Parser) parseNextExpr(prev Node) (Node, error) {
// 	if id := ps.maybeTok(Identifier); id.Ok {
// 		def := ps.scope.Search(id.Expr)
// 		if def == nil {
// 			return nil, ps.errorf(id, `Unknown name`)
// 		}

// 		return IdExpr{def}, nil
// 	}

// 	if str := ps.maybeTok(RawStr, Str); str.Ok {
// 		var content string

// 		switch str.Trait {
// 		case RawStr:
// 			content = str.Expr[1 : len(str.Expr)-1]
// 		case Str:
// 			content = UnescapeStr(str.Expr[1 : len(str.Expr)-1])
// 		}

// 		return StrExpr{content}, nil
// 	}

// 	// Integers and floats constants infers to 32-64 bits depending on the size
// 	// Must have an explicit cast in order to access smaller parts of the register

// 	if int := ps.maybeTok(IntDec, IntBin, IntHex); int.Ok {
// 		var (
// 			n    int64
// 			err  error
// 			size uint
// 		)

// 		switch int.Trait {
// 		case IntDec:
// 			n, err = strconv.ParseInt(int.Expr[:], 10, 64)
// 		case IntBin:
// 			n, err = strconv.ParseInt(int.Expr[2:], 2, 64)
// 		case IntHex:
// 			n, err = strconv.ParseInt(int.Expr[2:], 16, 64)
// 		}

// 		if n > math.MaxInt32 {
// 			size = 64
// 		} else {
// 			size = 32
// 		}

// 		return IntExpr{n, size}, err
// 	}

// 	if float := ps.maybeTok(Float); float.Ok {
// 		f, err := strconv.ParseFloat(float.Expr, 64)

// 		var size uint
// 		if f > math.MaxFloat32 {
// 			size = 64
// 		} else {
// 			size = 32
// 		}

// 		return FloatExpr{f, size}, err
// 	}

// 	// Implementation doesn't support multibyte character constants !
// 	if char := ps.maybeTok(Char); char.Ok {
// 		content := UnescapeStr(char.Expr[1 : len(char.Expr)-1])
// 		switch len(content) {
// 		case 0:
// 			return nil, ps.errorf(char, `Empty character constant`)
// 		case 1:
// 			return CharExpr{content[0]}, nil
// 		default:
// 			return nil, ps.errorf(char, `Character constant too long`)
// 		}
// 	}

// 	if sign := ps.maybeTok(Add, Sub); sign.Ok && prev == nil {
// 		expr, err := ps.parseNextExpr(nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return UnaryExpr{OrderPrev, expr, sign}, nil
// 	}

// 	if binaryOp := ps.maybeTok(
// 		Assign,
// 		KwAnd, KwOr,
// 		Add, Sub, Mul, Div, Mod,
// 		BinNot, BinAnd, BinOr, BinXor, BinShiftL, BinShiftR,
// 		Equal, NotEq, Less, Greater, LessEq, GreaterEq); binaryOp.Ok {
// 		if prev == nil {
// 			return nil, ps.errorf(binaryOp, `Missing pre-operand for binary expression`)
// 		}
// 		next, err := ps.parseNextExpr(nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if next == nil {
// 			return nil, ps.errorf(binaryOp, `Missing post-operand for binary expression`)
// 		}

// 		return BinaryExpr{prev, next, binaryOp}, nil
// 	}

// 	if incrementOp := ps.maybeTok(Increment, Decrement); incrementOp.Ok {
// 		if prev != nil {
// 			return UnaryExpr{OrderPost, prev, incrementOp}, nil
// 		} else {
// 			next, err := ps.parseNextExpr(nil)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if next == nil {
// 				return nil, ps.errorf(incrementOp, `Missing expression for increment`)
// 			}
// 			return UnaryExpr{OrderPrev, next, incrementOp}, nil
// 		}
// 	}

// 	if paren := ps.maybeTok(ParenBegin); paren.Ok {
// 		expr, err := ps.parseExpr(')')
// 		return NestedExpr{Body: []Node{expr}}, err
// 	}

// 	return nil, nil
// }

func (ps *Parser) peek() Token {
	tok := ps.sn.Tokenize()
	ps.peekQueue = append(ps.peekQueue, tok)
	return tok
}

func (ps *Parser) token(traits ...Trait) Token {
	var tok Token
	if len(ps.peekQueue) != 0 {
		tok = ps.peekQueue[0]
		ps.peekQueue = ps.peekQueue[1:]
	} else {
		tok = ps.sn.Tokenize()
	}
	tok.Ok = len(traits) == 0 || slices.Contains(traits[:], tok.Trait)
	if !tok.Ok {
		ps.peekQueue = append(ps.peekQueue, tok)
	}
	return tok
}

func (ps *Parser) errorf(tok Token, f string, args ...any) error {
	return fmt.Errorf(f, args...)
}
