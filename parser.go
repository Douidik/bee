package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
)

type Parser struct {
	sn        Scanner
	peekQueue []Token
	ast       Ast
	scope     *Scope
}

func NewParser(sc Scanner) Parser {
	return Parser{sn: sc, peekQueue: make([]Token, 2)}
}

func (p *Parser) Parse() (*Ast, error) {
	p.scope = p.ast.NewScope(nil)

	for !p.sn.Finished() {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		p.ast.head.Body = append(p.ast.head.Body, stmt)
	}

	return &p.ast, nil
}

func (p *Parser) parseStmt() (Node, error) {
	switch peek := p.peekTok(); peek.Trait {
	case KwIf:
	case KwFor:
	case KwReturn:
		break

	case Identifier:
		if def := p.maybeTok(Define); def.Ok {

		} else if decl := p.maybeTok(Declare); decl.Ok {

		} else {
			expr, err := p.parseExpr('\n')
			return ExprStmt{expr}, err
		}
	}

	// switch peek := p.peekTok(); peek.Trait {
	// case Identifier:
	// 	name := p.expectTok(Identifier)
	// 	op := p.expectTok(Define, Declare)

	// 	if !op.Ok {
	// 		return nil, p.errorf(op, `Expected operator '::' or ':' after identifier`)
	// 	}

	// 	// def, err := p.parseDef(name)
	// 	// if err != nil {
	// 	// 	return nil, err
	// 	// }

	// 	switch op.Trait {
	// 	case Declare:

	// 	case Define:

	// 	}
	// }

	return nil, nil
}

// func (p *Parser) parseDef(name Token) (Node, error) {
// 	if tok := p.maybeTok(ParenBegin); tok.Ok {
// 	}

// 	if tok := p.maybeTok(Identifier); tok.Ok {
// 		def := p.scope.Search(tok.Expr)
// 		if def == nil {
// 			return nil, p.errorf(tok, `Unknown name`)
// 		}

// 	}
// }

func (p *Parser) parseExpr(stop uint) (Node, error) {
	nested := NestedExpr{}
	var prev *Node = nil
	var end Token
	
	for {
		end = p.maybeTok(stop, End)
		if end.Ok {
			break
		}
		
		expr, err := p.parseNextExpr(prev)
		if err != nil {
			return nil, err
		}
		if expr != nil {
			nested.Body = append(nested.Body, expr)
		} else {
			break
		}
		prev = &expr
	}

	if end.Trait != stop {
		return nil, p.errorf(end, `Unexpected end of source during expression parsing`)
	}

	return nested, nil
}

func (p *Parser) parseNextExpr(prev *Node) (Node, error) {
	if id := p.maybeTok(Identifier); id.Ok {
		def := p.scope.Search(id.Expr)
		if def == nil {
			return nil, p.errorf(id, `Unknown name`)
		}

		return IdExpr{def}, nil
	}

	if str := p.maybeTok(RawStr, Str); str.Ok {
		content := str.Expr[1 : len(str.Expr) - 1]
		
		var content string

		switch str.Trait {
		case RawStr:
			content = str.Expr[1 : len(str.Expr)-1]
		case Str:
			content = UnescapeStr(str.Expr[1 : len(str.Expr)-1])
		}

		return StrExpr{content}, nil
	}

	// Integers and floats constants infers to 32-64 bits depending on the size
	// Must have an explicit cast in order to access smaller parts of the register
	
	if int := p.maybeTok(IntDec, IntBin, IntHex); int.Ok {
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

		return IntExpr{n, size}, err
	}

	if float := p.maybeTok(Float); float.Ok {
		f, err := strconv.ParseFloat(float.Expr, 64)
		
		var size uint
		if f > math.MaxFloat32 {
			size = 64
		} else {
			size = 32
		}

		return FloatExpr{f, size}, err
	}

	// Implementation doesn't support multibyte character constants !
	if char := p.maybeTok(Char); char.Ok {
		content := UnescapeStr(char.Expr[1 : len(char.Expr)-1])
		switch len(content) {
		case 0:
			return nil, p.errorf(char, `Empty character constant`)
		case 1:
			return CharExpr{content[0]}, nil
		default:
			return nil, p.errorf(char, `Character constant too long`)
		}
	}

	if sign := p.maybeTok(Add, Sub); sign.Ok && prev == nil {
		expr, err := p.parseNextExpr(nil)
		if err != nil {
			return nil, err
		}

		return UnaryExpr{OrderPrev, expr, sign}, nil
	}

	if binaryOp := p.maybeTok(
		Assign,
		KwAnd, KwOr,
		Add, Sub, Mul, Div, Mod,
		BinNot, BinAnd, BinOr, BinXor, BinShiftL, BinShiftR,
		Equal, NotEq, Less, Greater, LessEq, GreaterEq); binaryOp.Ok {
		if prev == nil {
			return nil, p.errorf(binaryOp, `Missing pre-operand for binary expression`)
		}
		next, err := p.parseNextExpr(nil)
		if err != nil {
			return nil, err
		}
		if next == nil {
			return nil, p.errorf(binaryOp, `Missing post-operand for binary expression`)
		}

		return BinaryExpr{prev, next, binaryOp}, nil
	}

	if incrementOp := p.maybeTok(Increment, Decrement); incrementOp.Ok {
		if prev != nil {
			return UnaryExpr{OrderPost, prev, incrementOp}, nil
		} else {
			next, err := p.parseNextExpr(nil)
			if err != nil {
				return nil, err
			}
			if next == nil {
				return nil, p.errorf(incrementOp, `Missing expression for increment`)
			}
			return UnaryExpr{OrderPrev, next, incrementOp}, nil
		}
	}

	if paren := p.maybeTok(ParenBegin); paren.Ok {
		expr, err := p.parseExpr(')')
		return NestedExpr{Body: []Node{expr}}, err
	}

	return nil, nil
}

func (p *Parser) tokenize() Token {
	if len(p.peekQueue) == 0 {
		tok := p.peekQueue[0]
		p.peekQueue = p.peekQueue[1:]
		return tok
	}

	return p.sn.Tokenize()
}

func (p *Parser) peekTok() Token {
	tok := p.tokenize()
	p.peekQueue = append(p.peekQueue, tok)
	return tok
}

func (p *Parser) expectTok(traits ...uint) Token {
	tok := p.tokenize()
	tok.Ok = len(traits) == 0 || slices.Contains(traits[:], tok.Trait)
	return tok
}

func (p *Parser) maybeTok(traits ...uint) Token {
	tok := p.expectTok()

	if !tok.Ok {
		p.peekQueue = append(p.peekQueue, tok)
	}

	return tok
}

func (p *Parser) errorf(tok Token, f string, args ...any) error {
	return fmt.Errorf(f, args...)
}
