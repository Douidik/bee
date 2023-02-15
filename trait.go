package main

const (
	None uint = iota

	NewLine
	Empty
	Blank
	End
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

func BeeTraitName(trait uint) string {
	switch trait {

	case NewLine:
		return "NewLine"
	case Empty:
		return "Empty"
	case Blank:
		return "Blank"
	case End:
		return "End"
	case Comment:
		return "Comment"
	case Directive:
		return "Directive"

	case KwStruct:
		return "KwStruct"
	case KwEnum:
		return "KwEnum"
	case KwUnion:
		return "KwUnion"
	case KwUnderscore:
		return "KwUnderscore"
	case KwSelf:
		return "KwSelf"
	case KwArrow:
		return "KwArrow"
	case KwBreak:
		return "KwBreak"
	case KwCase:
		return "KwCase"
	case KwContinue:
		return "KwContinue"
	case KwDo:
		return "KwDo"
	case KwElse:
		return "KwElse"
	case KwEach:
		return "KwEach"
	case KwFor:
		return "KwFor"
	case KwIf:
		return "KwIf"
	case KwReturn:
		return "KwReturn"
	case KwSwitch:
		return "KwSwitch"
	case KwAnd:
		return "KwAnd"
	case KwOr:
		return "KwOr"

	case Identifier:
		return "Identifier"

	case Float:
		return "Float"
	case Int:
		return "Int"
	case Str:
		return "Str"
	case Char:
		return "Char"

	case Increment:
		return "Increment"
	case Decrement:
		return "Decrement"
	case ParenBegin:
		return "ParenBegin"
	case ParenEnd:
		return "ParenEnd"
	case ScopeBegin:
		return "ScopeBegin"
	case ScopeEnd:
		return "ScopeEnd"
	case CrochetBegin:
		return "CrochetBegin"
	case CrochetEnd:
		return "CrochetEnd"
	case Declare:
		return "Declare"
	case Define:
		return "Define"
	case Assign:
		return "Assign"
	case Arrow:
		return "Arrow"
	case Not:
		return "Not"
	case Add:
		return "Add"
	case Sub:
		return "Sub"
	case Mul:
		return "Mul"
	case Div:
		return "Div"
	case Mod:
		return "Mod"
	case BinNot:
		return "BinNot"
	case BinAnd:
		return "BinAnd"
	case BinOr:
		return "BinOr"
	case BinXor:
		return "BinXor"
	case BinShiftL:
		return "BinShiftL"
	case BinShiftR:
		return "BinShiftR"
	case Equal:
		return "Equal"
	case NotEq:
		return "NotEq"
	case Less:
		return "Less"
	case Greater:
		return "Greater"
	case LessEq:
		return "LessEq"
	case GreaterEq:
		return "GreaterEq"
	case Ref:
		return "Ref"
	case Deref:
		return "Deref"
	case Dot:
		return "Dot"
	case Comma:
		return "Comma"
	case Semicolon:
		return "Semicolon"

	default:
		return "?"
	}
}
