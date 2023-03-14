package main

import (
	"fmt"
)

type Pattern struct {
	Trait Trait
	Regex Regex
}

type SyntaxMap []Pattern

func NewBeeSyntax() SyntaxMap {
	def := func(trait Trait, src string) Pattern {
		rx, err := NewRegex(src)
		if err != nil {
			fmt.Println(err)
			return Pattern{}
		}
		return Pattern{trait, rx}
	}

	return SyntaxMap{
		def(NewLine, "'\n'"),
		def(Blank, `_+`),
		def(Comment, `'//' {^  ~ '\n'}'`),
		def(Directive, `'#' {^  ~ '\n'}'`),

		def(KwStruct, `'struct'/!a`),
		def(KwEnum, `'enum'/!a`),
		def(KwUnion, `'union'/!a`),
		def(KwUnderscore, `'_'/!a`),
		def(KwSelf, `'$'`),
		def(KwBreak, `'break'/!a`),
		def(KwCase, `'case'/!a`),
		def(KwContinue, `'continue'/!a`),
		def(KwElse, `'else'/!a`),
		def(KwEach, `'each'/!a`),
		def(KwFor, `'for'/!a`),
		def(KwIf, `'if'/!a`),
		def(KwReturn, `'return'/!a`),
		def(KwSwitch, `'switch'/!a`),
		def(KwAnd, `'and'/!a`),
		def(KwOr, `'or'/!a`),
		def(KwOr, `'fn'/!a`),

		def(ParenBegin, `'('`),
		def(ParenEnd, `')'`),
		def(ScopeBegin, `'{'`),
		def(ScopeEnd, `'}'`),
		def(CrochetBegin, `'['`),
		def(CrochetEnd, `']'`),

		def(Arrow, `'->'`),
		def(Increment, `'++'`),
		def(Decrement, `'--'`),
		def(Add, `'+'`),
		def(Sub, `'-'`),

		def(Float, `{[0-9]+ '.' [0-9]*} | {[0-9]* '.' [0-9]+}`),
		def(IntDec, `[0-9]+`),
		def(IntBin, `'0b' [0-1]+`),
		def(IntHex, `'0x' {[0-9]|[a-f]|[A-F]}+`),

		def(RawStr, "Q^Q"),
		def(Str, "q {{{'\\'^}|^} ~ /{q|'\n'}} ? {q|'\n'}"),
		def(Char, "'`' {{{'\\'^}|^} ~ /{'`'|'\n'}} ? {'`'|'\n'}"),
		def(Identifier, `{a|'_'} {a|'_'|n}*`),

		def(Declare, `'::'`),
		def(Define, `':'`),
		def(BinNot, `'~'`),
		def(BinOr, `'|'`),
		def(BinXor, `'^'`),
		def(BinShiftL, `'<<'`),
		def(BinShiftR, `'>>'`),
		def(Div, `'/'`),
		def(Mod, `'%'`),
		def(Equal, `'=='`),
		def(NotEq, `'!='`),
		def(LessEq, `'<='`),
		def(GreaterEq, `'>='`),
		def(Less, `'<'`),
		def(Greater, `'>'`),
		def(Not, `'!'`),
		def(Assign, `'='`),
		def(Ref, `'&'`),
		def(Deref, `'*'`),
		def(Dot, `'.'`),
		def(Comma, `','`),
		def(Semicolon, `';'`),

		def(None, `^~_`),
	}
}
