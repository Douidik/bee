package main

type Scanner struct {
	src string
	cur int
	sm  SyntaxMap
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
		return Token{sn.src[len(sn.src)-1:], End, false}
	}

	for _, pt := range sn.sm {
		match := pt.Regex.Match(sn.src[sn.cur:])
		if match != -1 {
			expr := sn.src[sn.cur : sn.cur+match]
			sn.cur += match
			return Token{expr, pt.Trait, true}
		}
	}

	return Token{`<unreachable>`, None, false}
}
