package main

type Scanner struct {
	src  string
	cur  int
	smap SyntaxMap
}

func NewScanner(src string, sm SyntaxMap) Scanner {
	return Scanner{src, 0, sm}
}

func (sn *Scanner) Finished() bool {
	return sn.cur >= len(sn.src)
}

func (sn *Scanner) Tokenize() Token {
	tok := sn.match()

	if tok.Trait != Blank {
		return tok
	} else {
		return sn.Tokenize()
	}
}

func (sn *Scanner) match() Token {
	if sn.Finished() {
		return Token{len(sn.src) - 1, sn.src[len(sn.src)-1:], Eof, false}
	}

	for _, pt := range sn.smap {
		match := pt.Regex.Match(sn.src[sn.cur:])
		if match != -1 {
			index := sn.cur
			expr := sn.src[index : sn.cur+match]
			sn.cur += match
			return Token{index, expr, pt.Trait, true}
		}
	}

	return Token{0, `<unreachable>`, None, false}
}
