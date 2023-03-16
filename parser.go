package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Parser struct {
	name      string
	sn        Scanner
	peekQueue []Token
	ast       Ast
	scope     *Scope
}

func NewParser(name string, sn Scanner) Parser {
	return Parser{name: name, sn: sn, peekQueue: make([]Token, 0)}
}

func (ps *Parser) Parse() (*Ast, error) {
	ps.ast.Scope = NewScope(nil)
	ps.scope = ps.ast.Scope
	// ps.scope.Add(&Typedef{Name: "bool", Type: &Atom{size: 1}})
	// ps.scope.Add(&Typedef{Name: "char", Type: &Atom{size: 1, signed: true}})
	// ps.scope.Add(&Typedef{Name: "s8", Type: &Atom{size: 1, signed: true}})
	// ps.scope.Add(&Typedef{Name: "s16", Type: &Atom{size: 2, signed: true}})
	// ps.scope.Add(&Typedef{Name: "s32", Type: &Atom{size: 4, signed: true}})
	// ps.scope.Add(&Typedef{Name: "s64", Type: &Atom{size: 8, signed: true}})
	// ps.scope.Add(&Typedef{Name: "u8", Type: &Atom{size: 1}})
	// ps.scope.Add(&Typedef{Name: "u16", Type: &Atom{size: 2}})
	// ps.scope.Add(&Typedef{Name: "u32", Type: &Atom{size: 4}})
	// ps.scope.Add(&Typedef{Name: "u64", Type: &Atom{size: 8}})
	// ps.scope.Add(&Typedef{Name: "f32", Type: &Atom{size: 4, float: true}})
	// ps.scope.Add(&Typedef{Name: "f64", Type: &Atom{size: 8, float: true}})

	for !ps.sn.Finished() {
		node, err := ps.parseNode(NewLine)
		if err != nil {
			return nil, err
		}
		ps.ast.Body = append(ps.ast.Body, node)
	}
	return &ps.ast, nil
}

func (ps *Parser) parseNode(delim Trait) (Node, error) {
	var head Node = nil
	var last Token

	for !ps.sn.Finished() {
		if last = ps.token(delim); last.Ok {
			break
		}
		node, err := ps.expectNode(head, delim)
		if err != nil {
			return nil, err
		}
		if node != nil {
			head = node
		} else {
			break
		}
	}

	if !last.Ok {
		return nil, ps.errorf(last, "Expected <%s> got <%s>", delim.Repr(), last.Trait.Repr())
	}
	if head != nil {
		return head, nil
	} else {
		return ps.parseNode(delim)
	}
}

func (ps *Parser) expectNode(prev Node, delim Trait) (Node, error) {
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

		if cast, casted := def.(Type); casted && ps.token(ParenBegin).Ok {
			node, err := ps.parseNode(ParenEnd)
			if err != nil {
				return nil, err
			}
			if !node.Result().Cast(cast) {
				return nil, ps.errorf(id, "Cannot cast expression to type '%s'", id.Expr)
			}
			return Cast{node, cast}, nil
		}

		if init := ps.token(Define, Declare); init.Ok {
			expr, err := ps.parseNode(delim)
			if err != nil {
				return nil, err
			}
			def = ps.scope.Add(&Var{Name: id.Expr, Type: expr.Result()})
			switch init.Trait {
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
			size = 8
		} else {
			size = 4
		}
		return IntExpr{uint64(n), Atom{size: uint64(size), signed: true}}, err
	}

	if float := ps.token(Float); float.Ok {
		f, err := strconv.ParseFloat(float.Expr, 64)

		var size uint64
		if f > math.MaxFloat32 {
			size = 8
		} else {
			size = 4
		}
		return FloatExpr{float64(f), Atom{size: size, float: true}}, err
	}

	if char := ps.token(Char); char.Ok {
		content := UnescapeStr(char.Expr[1 : len(char.Expr)-1])
		switch len(content) {
		case 0:
			return nil, ps.errorf(char, "Empty character constant")
		case 1:
			return CharExpr{content[0]}, nil
		default:
			return nil, ps.errorf(char, "Character constant too long")
		}
	}

	if prev == nil {
		if sign := ps.token(Add, Sub); sign.Ok {
			expr, err := ps.expectNode(nil, delim)
			if err != nil {
				return nil, err
			}
			return UnaryExpr{OrderPrev, expr, sign}, nil
		}
	}

	if bin := ps.token(
		Assign,
		KwAnd, KwOr,
		Add, Sub, Mul, Div, Mod,
		BinNot, BinAnd, BinOr, BinXor, BinShiftL, BinShiftR,
		Equal, NotEq, Less, Greater, LessEq, GreaterEq); bin.Ok {
		if prev == nil {
			return nil, ps.errorf(bin, "Missing pre-operand for binary expression")
		}
		next, err := ps.expectNode(nil, delim)
		if err != nil {
			return nil, err
		}
		if next == nil {
			return nil, ps.errorf(bin, "Missing post-operand for binary expression")
		}
		if !prev.Result().Cast(next.Result()) {
			return nil, ps.errorf(bin, "Incompatible operands in binary expression")
		}
		return BinaryExpr{[2]Node{prev, next}, bin}, nil
	}

	if incr := ps.token(Increment, Decrement); incr.Ok {
		if prev != nil {
			return UnaryExpr{OrderPost, prev, incr}, nil
		} else {
			next, err := ps.expectNode(nil, delim)
			if err != nil {
				return nil, err
			}
			if next == nil {
				return nil, ps.errorf(incr, "Missing expression for increment")
			}
			return UnaryExpr{OrderPrev, next, incr}, nil
		}
	}

	if ps.token(ScopeBegin).Ok {
		return ps.parseCompound(NewLine, ScopeEnd)
	}

	if ps.token(ParenBegin).Ok {
		nest := Nest{Body: make([]Node, 0)}

		for !ps.sn.Finished() && !ps.token(ParenEnd).Ok {
			node, err := ps.parseNode(Comma)
			if err != nil {
				return nest, err
			}
			nest.Body = append(nest.Body, node)
		}
		return nest, nil
	}

	return nil, nil
}

func (ps *Parser) parseCompound(delim, end Trait) (Compound, error) {
	ps.scope = NewScope(ps.scope)
	compound := Compound{Scope: ps.scope, Body: make([]Node, 0)}

	for !ps.sn.Finished() && !ps.token(end).Ok {
		node, err := ps.parseNode(delim)
		if err != nil {
			return compound, err
		}
		compound.Body = append(compound.Body, node)
	}

	ps.scope = ps.scope.Owner
	return compound, nil
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

//	Example: from 'basic.bee':24 > foo :: fn () -> {
//	                                               ^ Function return type expected in signature after '->'

func (ps *Parser) errorf(tok Token, f string, args ...any) error {
	trim := func(str string, index, end int) int {
		step := func(a, b int) int {
			d := b - a
			if d > 0 {
				return +1
			} else {
				return -1
			}
		}
		// Search new line in the str[index : end] interval
		newLine := index
		for newLine != end && str[newLine] != '\n' {
			newLine += step(index, end)
		}
		// Search first non-blank character in the str[newLine : index] interval
		begin := newLine
		for begin != index && strings.IndexByte(" \t\n\v\f\r", str[begin]) != -1 {
			begin += step(newLine, index)
		}
		return begin
	}

	// todo: trim function not required, just take a slice and rebase the index
	src := ps.sn.src
	begin := trim(src, tok.Index, 0)
	end := begin + strings.IndexByte(src[begin:], '\n')
	line := 1 + strings.Count(src[:begin], "\n")
	location := fmt.Sprintf("from '%s':%d > ", ps.name, line)
	snippet := src[begin:end]
	cursor := len(location) + (tok.Index - begin + 1)
	reason := fmt.Sprintf(f, args...)
	return fmt.Errorf("%s%s\n%*c %s", location, snippet, cursor, '^', reason)
}

// func (ps *Parser) errorf(tok Token, f string, args ...any) error {
// 	var (
// 		src   = ps.sn.src
// 		begin int
// 		end   int
// 		line  int
// 	)

// 	if begin = strings.LastIndexByte(src[:tok.Index], '\n'); begin < 0 {
// 		begin = 0
// 	}
// 	if end = strings.IndexByte(src[:tok.Index], '\n'); end < 0 {
// 		end = len(src)
// 	}
// 	line = strings.Count(src[:begin], "\n")

// 	location := fmt.Sprintf("%s:%d > ", ps.name, line)
// 	snippet := strings.Trim(src[begin:end], " \t\n\v\f\r")
// 	offset := len(location) + (tok.Index - begin + 1)
// 	reason := fmt.Sprintf(f, args...)
// 	return fmt.Errorf("%s%s\n%*c %s", location, snippet, offset, '^', reason)
// }
