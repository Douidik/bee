package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

type stateTag uint
type nodes []*node
type stack []*node

const (
	epsilon stateTag = iota
	anything
	none
	not
	dash
	text
	set
	scope
)

type state struct {
	Tag stateTag
	str string
	seq *node
	a   byte
	b   byte
}

type node struct {
	state state
	index int
	edges nodes
}

type parser struct {
	sr    *strings.Reader
	stack stack
}

type Regex struct {
	Src  string
	Head *node
}

type RegexGraph struct {
	sb   strings.Builder
	head *node
	Name string
}

func NewRegex(src string) (Regex, error) {
	p := newRegexParser(src)
	head, err := p.Parse()

	if err != nil {
		return Regex{}, err
	}

	return Regex{Src: src, Head: head}, nil
}

func (rx *Regex) Match(expr string) int {
	return rx.Head.Submit(expr, 0)
}

func (rx *Regex) Graph(name string) string {
	graph := NewRegexGraph(rx, name)
	return graph.Document()
}

func (s nodes) Len() int {
	return len(s)
}

func (s nodes) Less(i, j int) bool {
	return s[i].index < s[j].index
}

func (s nodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s nodes) Append(n *node) (nodes, *node) {
	s_ := append(s, n)
	sort.Sort(s_)
	return s_, n
}

func (s nodes) Max() *node {
	sort.Sort(s)
	if s.Len() < 1 {
		return nil
	}
	return s[s.Len()-1]
}

func newEpsilon() state {
	return state{Tag: epsilon}
}

func newAnything() state {
	return state{Tag: anything}
}

func newNone() state {
	return state{Tag: none}
}

func newNot(node *node) state {
	return state{Tag: not, seq: node}
}

func newDash(node *node) state {
	return state{Tag: dash, seq: node}
}

func newText(str string) state {
	return state{Tag: text, str: str}
}

func newSet(str string) state {
	return state{Tag: set, str: str}
}

func newScope(a byte, b byte) state {
	return state{Tag: scope, a: a, b: b}
}

func newNode(state state) *node {
	return &node{
		state: state,
		index: 0,
		edges: make(nodes, 0, 16),
	}
}

func (s *state) Submit(expr string, index int) int {
	if s.Tag != epsilon && index >= len(expr) {
		return -1
	}

	switch s.Tag {
	case epsilon:
		return index

	case anything:
		return index + 1

	case not:
		if match := s.seq.Submit(expr, index); match == -1 {
			return index + 1
		}

	case dash:
		if match := s.seq.Submit(expr, index); match != -1 {
			return index
		}

	case text:
		length := min(len(s.str), len(expr)-index)

		if s.str[:length] == expr[index:index+length] {
			return index + length
		}

	case set:
		if strings.IndexByte(s.str, expr[index]) != -1 {
			return index + 1
		}

	case scope:
		if s.a <= expr[index] && expr[index] <= s.b {
			return index + 1
		}
	}

	return -1
}

func (n *node) Submit(expr string, index int) int {
	match := n.state.Submit(expr, index)

	if match != -1 {
		branch := n.Branch()
		if !branch && match >= len(expr) {
			return match
		}

		for _, edge := range n.edges {
			matchFwd := edge.Submit(expr, match)
			if matchFwd != -1 {
				return matchFwd
			}
		}

		if !branch {
			return match
		}
	}

	return -1
}

func (n *node) makeMembers(membs *nodes) *nodes {
	*membs = append(*membs, n)

	for _, edge := range n.edges {
		if edge.index > n.index {
			edge.makeMembers(membs)
		}
	}

	return membs
}

func (n *node) Membs() nodes {
	membs := make(nodes, 0, 32)
	return *n.makeMembers(&membs)
}

func (n *node) NextIndex() int {
	// max := n.Membs().Max()
	// if max != nil {
	// 	return max.index + 1
	// } else {

	// }

	return n.Membs().Max().index + 1
}

func (n *node) Branch() bool {
	return n.edges.Len() > 0 && n.edges.Max().index > n.index
}

func (n *node) Push(edge *node) *node {
	edge.Scope(n.NextIndex())
	n.edges, _ = n.edges.Append(edge)
	return edge
}

func (n *node) Concat(edge *node) *node {
	for _, member := range n.Membs() {
		if !member.Branch() {
			member.edges, _ = member.edges.Append(edge)
		}
	}
	return edge
}

func (n *node) Scope(base int) {
	for _, member := range n.Membs() {
		member.index += base
	}
}

func (n *node) Merge(seq *node) {
	seq.Scope(n.NextIndex())
	n.Concat(seq)
}

func (s stack) Empty() bool {
	return len(s) == 0
}

func (s stack) Push(node *node) (stack, *node) {
	return append(s, node), node
}

func (s stack) Pop() (stack, *node) {
	if s.Empty() {
		return nil, nil
	}
	return s[:len(s)-1], s[len(s)-1]
}

func newRegexParser(src string) parser {
	return parser{sr: strings.NewReader(src)}
}

func (p *parser) Parse() (*node, error) {
	for p.sr.Len() != 0 {
		seq, err := p.nextToken()
		if err != nil {
			return nil, err
		}
		if seq != nil {
			p.stack, _ = p.stack.Push(seq)
		}
	}

	if len(p.stack) != 0 {
		for i := 1; i < len(p.stack); i++ {
			p.stack[0].Merge(p.stack[i])
		}
		return p.stack[0], nil
	}
	return nil, nil
}

func (p *parser) nextToken() (*node, error) {
	tok, err := p.sr.ReadByte()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	switch tok {
	case ' ', '\f', '\n', '\r', '\t', '\v':
		return p.nextToken()

	case '_':
		return p.parseSet(" \n\v\b\f\t")
	case 'a':
		return p.parseSet("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	case 'o':
		return p.parseSet("!#$%&()*+,-./:;<=>?@[\\]^`{|}~")
	case 'n':
		return p.parseSet("0123456789")
	case 'Q':
		return p.parseSet("\"")
	case 'q':
		return p.parseSet("'")

	case '!':
		return p.parseNot()
	case '/':
		return p.parseDash()

	case '[':
		return p.parseScope()
	case '^':
		return p.parseAnything()
	case '\'':
		return p.parseText('\'')
	case '`':
		return p.parseText('`')
	case '{':
		return p.parseSequence()
	case '|':
		return p.parseOr()
	case '?':
		return p.parseQuest()
	case '+':
		return p.parsePlus()
	case '*':
		return p.parseStar()
	case '~':
		return p.parseWave()

	case '}':
		return nil, p.errorf("Unmatched sequence brace, missing <{> operator")
	case ']':
		return nil, p.errorf("Unmatched scope brace, missing <[> operator")
	default:
		return nil, p.errorf("'%c': Unrecognized token in regex, none of [_aonQq^'{}!|?*+~]", tok)
	}
}

func (p *parser) parsePreOp(op byte) (*node, error) {
	var seq *node

	if p.stack, seq = p.stack.Pop(); seq == nil {
		return seq, p.errorf("Missing pre-operand here for <%c>", op)
	}

	return seq, nil
}

func (p *parser) parsePostOp(op byte) (*node, error) {
	seq, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if seq == nil {
		return nil, p.errorf("Missing post-operand here for <%c>", op)
	}
	return seq, err
}

func (p *parser) parseBinOp(op byte) (*node, *node, error) {
	var pre, post *node
	var err error

	pre, err = p.parsePreOp(op)
	if err != nil {
		return nil, nil, err
	}
	post, err = p.parsePostOp(op)
	if err != nil {
		return nil, nil, err
	}

	return pre, post, err
}

func (p *parser) parseSet(str string) (*node, error) {
	return newNode(newSet(str)), nil
}

// NOTE: Cannot because of initialization cycle error, may come up with a solution later

// var scopeFormat = NewRegex("'['? ^ '-' ^ ']'")

// func (p *parser) parseScope() (*node, error) {
// 	buf := make([]byte, 0, 5)
// 	p.sr.Read(buf)

// 	if scopeFormat.Match(string(buf)) != -1 {
// 		return nil, p.errorf("Scope does not match the format: %s", scopeFormat.Src)
// 	}

// 	return newNode(newScope(buf[0], buf[2])), nil
// }

// sad way
func (p *parser) parseScope() (*node, error) {
	// [ ^ - ^ ]
	//-1 0 1 2 3
	a, _ := p.sr.ReadByte()
	if hyphen, _ := p.sr.ReadByte(); hyphen != '-' {
		return nil, p.errorf("Expected <-> token between scope boundaries")
	}
	b, _ := p.sr.ReadByte()
	if end, _ := p.sr.ReadByte(); end != ']' {
		return nil, p.errorf("Expected <]> token at the end of the scope declaration")
	}

	if a > b {
		return nil, p.errorf(`Scope is not matchable because the interval of characters is empty`)
	}

	return newNode(newScope(a, b)), nil
}

func (p *parser) parseAnything() (*node, error) {
	return newNode(newAnything()), nil
}

func (p *parser) parseText(stop byte) (*node, error) {
	readText := func(buf *[]byte) error {
		for {
			c, err := p.sr.ReadByte()
			switch {
			case err != nil:
				return err
			case c == stop:
				return nil
			default:
				*buf = append(*buf, c)
			}
		}

	}

	buf := make([]byte, 0)
	if err := readText(&buf); err != nil {
		return nil, err
	}
	str := string(buf)
	return newNode(newText(str)), nil
}

func (p *parser) readSequenceSrc() (string, error) {
	depth := 1
	buf := make([]byte, 0)

	for {
		tok, err := p.sr.ReadByte()
		if err != nil {
			return "", p.errorf("Unmatched sequence brace, missing <}> token")
		}

		switch tok {
		case '{':
			depth++
		case '}':
			depth--
		}

		if depth < 1 {
			break
		} else {
			buf = append(buf, tok)
		}
	}

	return string(buf), nil
}

func (p *parser) parseSequence() (*node, error) {
	src, err := p.readSequenceSrc()
	if err != nil {
		return nil, err
	}

	sp := newRegexParser(src)
	return sp.Parse()
}

func (p *parser) parseDash() (*node, error) {
	seq, err := p.parsePostOp('/')
	if err != nil {
		return nil, err
	}

	return newNode(newDash(seq)), err
}

func (p *parser) parseNot() (*node, error) {
	seq, err := p.parsePostOp('!')
	if err != nil {
		return nil, err
	}

	return newNode(newNot(seq)), err
}

func (p *parser) parseOr() (*node, error) {
	pre, post, err := p.parseBinOp('|')
	if err != nil {
		return nil, err
	}

	or := newNode(newEpsilon())
	or.Push(pre)
	or.Push(post)

	return or, err
}

func (p *parser) parseQuest() (*node, error) {
	pre, err := p.parsePreOp('?')
	if err != nil {
		return nil, err
	}

	quest := newNode(newEpsilon())
	quest.Push(pre)
	quest.Push(newNode(newEpsilon()))

	return quest, err
}

func (p *parser) parseStar() (*node, error) {
	pre, err := p.parsePreOp('*')
	if err != nil {
		return nil, err
	}

	star := newNode(newEpsilon())
	star.Merge(pre)
	star.Concat(star)
	star.Push(newNode(newEpsilon()))

	return star, err
}

func (p *parser) parsePlus() (*node, error) {
	plus, err := p.parsePreOp('+')
	if err != nil {
		return nil, err
	}

	plus.Concat(plus)
	return plus, err
}

func (p *parser) parseWave() (*node, error) {
	pre, post, err := p.parseBinOp('~')
	if err != nil {
		return nil, err
	}

	wave := newNode(newEpsilon())
	wave.Push(post)
	wave.Push(pre)
	pre.Concat(wave)
	pre.Merge(newNode(newNone()))

	return wave, err
}

func (p *parser) errorf(desc string, args ...interface{}) error {
	return fmt.Errorf(desc, args...)
}

func NewRegexGraph(rx *Regex, name string) RegexGraph {
	return RegexGraph{
		sb:   strings.Builder{},
		head: rx.Head,
		Name: name,
	}
}

func (rg *RegexGraph) write(f string, args ...interface{}) {
	rg.sb.WriteString(fmt.Sprintf(f, args...))
}

func (rg *RegexGraph) writeln(f string, args ...interface{}) {
	rg.write(f, args...)
	rg.sb.WriteByte('\n')
}

func (rg *RegexGraph) Document() string {
	rg.sb.Reset()
	rg.writeln(`strict digraph {`)

	if rg.head != nil {
		rg.writeln(`rankdir=LR;bgcolor="#F9F9F9";compound=true`)
		rg.writeln(`"%s" [shape="none"]`, rg.Name)
		rg.writeln(`"%s" -> "%p" [label="%s"]`, rg.Name, rg.head, rg.makeState(rg.head))

		for _, member := range rg.head.Membs() {
			rg.format(member)
		}
	}

	rg.writeln(`}`)
	return rg.sb.String()
}

func (rg *RegexGraph) define(n *node) {
	shape := "?"
	if n.Branch() {
		shape = "square"
	} else {
		shape = "circle"
	}
	rg.writeln(`"%p" [shape="%s", label="%d"]`, n, shape, n.index)
}

func (rg *RegexGraph) connect(a *node, b *node) {
	rg.writeln(`"%p" -> "%p" [label="%s"]`, a, b, rg.makeState(b))
}

func (rg *RegexGraph) formatSubgraph(n *node, header string) {
	seq := n.state.seq

	rg.writeln(`subgraph cluster_%p {`, n)
	rg.writeln(`%s`, header)
	rg.define(n)
	rg.connect(n, seq)

	membs := seq.Membs()
	max := membs.Max()
	for _, member := range membs {
		rg.format(member)
	}

	rg.writeln(`}`)
	for _, edge := range n.edges {
		rg.connect(max, edge)
	}
}

func (rg *RegexGraph) format(n *node) {
	s := n.state

	switch s.Tag {
	case not:
		rg.formatSubgraph(n, `style=filled;bgcolor="#FBF3F3"`)
	case dash:
		rg.formatSubgraph(n, `style=filled;bgcolor="#F4FDFF"`)
	default:
		rg.define(n)
		for _, edge := range n.edges {
			rg.connect(n, edge)
		}
	}
}

func (rg *RegexGraph) makeState(n *node) string {
	s := n.state

	switch s.Tag {
	case epsilon:
		return "&Sigma;"
	case anything:
		return "&alpha;"
	case none:
		return "&times;"
	case not:
		return "!"
	case dash:
		return "/"

	case text:
		return fmt.Sprintf("'%s'", s.str)

	case set:
		switch len(s.str) {
		case 0:
			return "[]"
		case 1:
			return fmt.Sprintf("[%c]", s.str[0])
		default:
			return fmt.Sprintf("[%c..%c]", s.str[0], s.str[len(s.str)-1])
		}

	case scope:
		return fmt.Sprintf("[%c-%c]", s.a, s.b)
	}

	return "?"
}
