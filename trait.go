package main

type Trait uint

const (
	None Trait = iota

	NewLine
	Empty
	Blank
	Eof
	Comment
	Directive

	KwStruct
	KwEnum
	KwUnion
	KwUnderscore
	KwSelf
	KwArrow
	KwBreak
	KwCase
	KwContinue
	KwElse
	KwEach
	KwFor
	KwIf
	KwReturn
	KwSwitch
	KwAnd
	KwOr
	KwFn

	Identifier

	Float
	IntDec
	IntBin
	IntHex
	Str
	RawStr
	Char

	Increment
	Decrement
	ParenBegin
	ParenEnd
	ScopeBegin
	ScopeEnd
	CrochetBegin
	CrochetEnd
	Declare
	Define
	Assign
	Arrow
	Not
	Add
	Sub
	Mul
	Div
	Mod
	BinNot
	BinAnd
	BinOr
	BinXor
	BinShiftL
	BinShiftR
	Equal
	NotEq
	Less
	Greater
	LessEq
	GreaterEq
	Ref
	Deref
	Dot
	Comma
	Semicolon
)

func BeeTraitName(trait Trait) string {
	switch trait {

	case NewLine:
		return "NewLine"
	case Empty:
		return "Empty"
	case Blank:
		return "Blank"
	case Eof:
		return "Eof"
	case Comment:
		return "Comment"
	case Directive:
		return "Directive"

	case KwStruct:
		return "Struct"
	case KwEnum:
		return "Enum"
	case KwUnion:
		return "Union"
	case KwUnderscore:
		return "Underscore"
	case KwSelf:
		return "Self"
	case KwArrow:
		return "Arrow"
	case KwBreak:
		return "Break"
	case KwCase:
		return "Case"
	case KwContinue:
		return "Continue"
	case KwElse:
		return "Else"
	case KwEach:
		return "Each"
	case KwFor:
		return "For"
	case KwIf:
		return "If"
	case KwReturn:
		return "Return"
	case KwSwitch:
		return "Switch"
	case KwAnd:
		return "And"
	case KwOr:
		return "Or"
	case KwFn:
		return "Fn"

	case Identifier:
		return "Identifier"

	case Float:
		return "Float"
	case IntDec:
		return "IntDec"
	case IntBin:
		return "IntBin"
	case IntHex:
		return "IntHex"
	case Str:
		return "Str"
	case Char:
		return "Char"

	case Increment:
		return "++"
	case Decrement:
		return "--"
	case ParenBegin:
		return "("
	case ParenEnd:
		return ")"
	case ScopeBegin:
		return "{"
	case ScopeEnd:
		return "}"
	case CrochetBegin:
		return "["
	case CrochetEnd:
		return "]"
	case Declare:
		return ":"
	case Define:
		return "::"
	case Assign:
		return "="
	case Arrow:
		return "->"
	case Not:
		return "!"
	case Add:
		return "+"
	case Sub:
		return "-"
	case Mul:
		return "*"
	case Div:
		return "/"
	case Mod:
		return "*"
	case BinNot:
		return "~"
	case BinAnd:
		return "&"
	case BinOr:
		return "|"
	case BinXor:
		return "^"
	case BinShiftL:
		return "<<"
	case BinShiftR:
		return ">>"
	case Equal:
		return "=="
	case NotEq:
		return "!="
	case Less:
		return "<"
	case Greater:
		return ">"
	case LessEq:
		return "<="
	case GreaterEq:
		return ">="
	case Ref:
		return "&"
	case Deref:
		return "*"
	case Dot:
		return "."
	case Comma:
		return ","
	case Semicolon:
		return ";"

	default:
		return "?"
	}
}
