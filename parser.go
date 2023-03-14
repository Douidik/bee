package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Parser struct {
	sn        Scanner
	peekQueue []Token
	ast       Ast
	scope     *Scope
}

func NewParser(sc Scanner) Parser {
	return Parser{sn: sc, peekQueue: make([]Token, 0)}
}

func (ps *Parser) Parse() (*Ast, error) {
	ps.scope = ps.ast.NewScope(nil)
	ps.scope.Add(&Typedef{Name: "bool", Type: &Atom{size: 8}})
	ps.scope.Add(&Typedef{Name: "char", Type: &Atom{size: 8, signed: true}})
	ps.scope.Add(&Typedef{Name: "s8", Type: &Atom{size: 8, signed: true}})
	ps.scope.Add(&Typedef{Name: "s16", Type: &Atom{size: 16, signed: true}})
	ps.scope.Add(&Typedef{Name: "s32", Type: &Atom{size: 32, signed: true}})
	ps.scope.Add(&Typedef{Name: "s64", Type: &Atom{size: 64, signed: true}})
	ps.scope.Add(&Typedef{Name: "u8", Type: &Atom{size: 8}})
	ps.scope.Add(&Typedef{Name: "u16", Type: &Atom{size: 16}})
	ps.scope.Add(&Typedef{Name: "u32", Type: &Atom{size: 32}})
	ps.scope.Add(&Typedef{Name: "u64", Type: &Atom{size: 64}})
	ps.scope.Add(&Typedef{Name: "f32", Type: &Atom{size: 32, float: true}})
	ps.scope.Add(&Typedef{Name: "f64", Type: &Atom{size: 64, float: true}})

	for !ps.sn.Finished() {
		body, err := ps.parseNode(NewLine, Eof)
		if err != nil {
			return nil, err
		}
		ps.ast.Body = append(ps.ast.Body, body...)
	}
	return &ps.ast, nil
}

func (ps *Parser) parseNode(delim, end Trait) ([]Node, error) {
	body := make([]Node, 0)
	var prev Node = nil
	var tok Token

	for !ps.sn.Finished() {
		if tok = ps.token(delim); tok.Ok {
			prev = nil
		}
		if tok = ps.token(end); tok.Ok {
			break
		}
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

	if !tok.Ok {
		return nil, ps.errorf(tok, "Expected <%s> got <%s>", end.Repr(), tok.Trait.Repr())
	}

	return body, nil
}

func (ps *Parser) expectNode(prev Node) (Node, error) {
	for ps.token(NewLine).Ok {
		// Ignore new line when expecting a node
	}

	if ps.token(KwIf).Ok {
		var (
			i   If
			err error
		)
		if i.Conds, err = ps.parseCompound(NewLine, ScopeBegin); err != nil {
			return nil, err
		}
		if i.If, err = ps.parseCompound(NewLine, ScopeEnd); err != nil {
			return nil, err
		}
		if ps.token(KwElse).Ok {
			if i.Else, err = ps.parseCompound(NewLine, ScopeEnd); err != nil {
				return nil, err
			}
		}
		return i, nil
	}

	if ps.token(KwFor).Ok {
		var (
			f   For
			err error
		)
		if f.Conds, err = ps.parseCompound(NewLine, ScopeBegin); err != nil {
			return nil, err
		}
		if f.Body, err = ps.parseCompound(NewLine, ScopeEnd); err != nil {
			return nil, err
		}
		return f, nil
	}

	if id := ps.token(Identifier); id.Ok {
		def := ps.scope.Search(id.Expr)
		if def != nil {
			return Reference{Def: def}, nil
		}

		if cast, isCast := def.(Type); isCast && ps.token(ParenBegin).Ok {
			body, err := ps.parseNode(NewLine, ParenEnd)
			if err != nil {
				return nil, err
			}
			switch {
			case len(body) == 0:
				return nil, ps.errorf(id, "Expected expression in '%s' casting parenthesis", id.Expr)
			case len(body) > 1:
				return nil, ps.errorf(id, "Extraneous expressions in '%s' casting parenthesis", id.Expr)
			case !body[0].Result().Infers(cast):
				return nil, ps.errorf(id, "Cannot cast expression to type '%s'", id.Expr)
			}
			return Cast{body[0], cast}, nil
		}

		if operator := ps.token(Define, Declare); operator.Ok {
			expr, err := ps.expectNode(nil)
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, ps.errorf(operator, "No value given to '%s'", id.Expr)
			}
			def = ps.scope.Add(&Var{Name: id.Expr, Type: expr.Result()})
			switch operator.Trait {
			case Define:
				return DefineExpr{def, expr}, nil
			case Declare:
				return DeclareExpr{def, expr}, nil
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
		return IntExpr{uint64(n), Atom{size: size, signed: true}}, err
	}

	if float := ps.token(Float); float.Ok {
		f, err := strconv.ParseFloat(float.Expr, 64)

		var size uint
		if f > math.MaxFloat32 {
			size = 64
		} else {
			size = 32
		}
		return FloatExpr{float64(f), Atom{size: size, float: true}}, err
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
		if !prev.Result().Infers(next.Result()) {
			return nil, ps.errorf(bin, `Incompatible operands in binary expression`)
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

	if ps.token(ScopeBegin).Ok {
		return ps.parseCompound(NewLine, ScopeEnd)
	}

	if ps.token(ParenBegin).Ok {
		body, err := ps.parseNode(NewLine, ParenEnd)
		return Nest{body}, err
	}

	return nil, nil
}

func (ps *Parser) parseCompound(delim, end Trait) (Compound, error) {
	ps.scope = &Scope{Defs: map[string]Def{}, Owner: ps.scope}
	body, err := ps.parseNode(delim, end)
	compound := Compound{Scope: ps.scope, Body: body}
	ps.scope = ps.scope.Owner
	return compound, err

}

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
	var (
		src = ps.sn.src
		begin int
		end int
		count int
	)
	
	if begin = strings.IndexByte(src[:tok.Index], "\n"); begin < 0 {
		begin = 0
	}
	if end = strings.IndexByte(src[tok.Index:], "\n"); end < 0 {
		
	}
	count := strings.Count(src[:begin], "\n")
	
		
	line := ps.sn.src[tok.Index:]
	line, _, _ = strings.Cut(line, "\n")
	lineNum := strings.Count(line, "\n"
	
	_, line, _ = strings.Cut(line, "\n")
	
	
	return fmt.Errorf(f, args...)
}
