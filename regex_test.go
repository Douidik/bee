package main

import (
	"fmt"
	"testing"
)

const LoremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
labore et dolore magna aliqua. Id porta nibh venenatis cras sed felis eget velit. Viverra tellus
in hac habitasse. Sed risus pretium quam vulputate dignissim suspendisse in est. In eu mi
bibendum neque egestas congue quisque egestas. Mi proin sed libero enim sed faucibus turpis in.
Aliquam vestibulum morbi blandit cursus. Tellus in hac habitasse platea dictumst vestibulum.
Massa ultricies mi quis hendrerit. Molestie a iaculis at erat pellentesque adipiscing commodo.
Vulputate eu scelerisque felis imperdiet proin fermentum. Vitae congue eu consequat ac felis. Nec
ultrices dui sapien eget mi proin sed. Nunc mattis enim ut tellus elementum sagittis vitae et.
Mauris ultrices eros in cursus turpis massa tincidunt dui ut. Nisi porta lorem mollis aliquam ut
porttitor leo a diam. Diam phasellus vestibulum lorem sed risus ultricies. Arcu vitae elementum
curabitur vitae nunc sed velit dignissim. Ut eu sem integer vitae justo eget magna fermentum
iaculis.In eu mi bibendum neque.`

func expectRegex(ts *testing.T, src string, expr string, crashExpected bool, matchExpected bool, eq int) {
	regex, err := NewRegex(src)
	if err != nil {
		if !crashExpected {
			ts.Log(err)
			ts.Fail()
		}
		return
	}
	regex.Graph("test")

	match := regex.Match(expr)
	if matchExpected && match == -1 {
		ts.Logf(`"%s" doesn't matched "%s"`, src, expr)
		ts.Fail()
	}
	if !matchExpected && match != -1 {
		ts.Logf(`"%s" matched "%s" of "%s"`, src, expr[:match], expr)
		ts.Fail()
	}
	if match != -1 && eq != -1 && expr[:match] != expr[:eq] {
		ts.Logf(`"%s" matched "%s" of "%s" instead of "%s"`, src, expr[:match], expr, expr[:eq])
		ts.Fail()
	}
}

func expectMatch(ts *testing.T, src string, expr string) {
	expectRegex(ts, src, expr, false, true, -1)
}

func expectMatchEq(ts *testing.T, src string, expr string, eq int) {
	expectRegex(ts, src, expr, false, true, eq)
}

func expectNoMatch(ts *testing.T, src string, expr string) {
	expectRegex(ts, src, expr, false, false, -1)
}

func expectError(ts *testing.T, src string) {
	expectRegex(ts, src, "", true, false, -1)
}

func TestRegexUnknown(ts *testing.T) {
	expectError(ts, "N")
	expectError(ts, ")")
	expectError(ts, "\"")
}

func TestRegexText(ts *testing.T) {
	expectMatch(ts, "'abc'", "abc")
	expectMatch(ts, "'abc'", "abcccccccccc")
	expectMatch(ts, "'hello' ' ' 'world'", "hello world")
	expectMatch(ts, "'hello\nworld'", "hello\nworld")
	expectMatch(ts, fmt.Sprintf(`'%s'`, LoremIpsum), LoremIpsum)

	expectMatch(ts, "`abc`", "abc")
	expectMatch(ts, "`abc`", "abcccccccccc")
	expectMatch(ts, "`hello` ` ` `world`", "hello world")
	expectMatch(ts, "`hello\nworld`", "hello\nworld")

	expectError(ts, "`hello'")
	expectError(ts, "'hello`")
	expectError(ts, "'hello")
	expectError(ts, "hello'")
	expectError(ts, "hello`")
	expectError(ts, "`hello")
	expectError(ts, "hello")

	expectNoMatch(ts, "'cba'", "abc")
	expectNoMatch(ts, "'cbaa'", "abcc")
	expectNoMatch(ts, fmt.Sprintf("`%s`", LoremIpsum), LoremIpsum[1:])
	expectNoMatch(ts, fmt.Sprintf("`%s`", LoremIpsum), LoremIpsum[2:len(LoremIpsum)-2])
}

func TestRegexRange(ts *testing.T) {
	expectMatchEq(ts, "[0-9]+", "0123456789yeet", 10)
	expectMatchEq(ts, "[a-f]+", "abcdefghijklmnopqrstuvwxyz", 6)
	expectMatchEq(ts, "[a-a]+", "aaaaaaah", 7)
	expectMatchEq(ts, "[[-]]+", `[\]`, 3)
	expectMatchEq(ts, "[---]+", "--", 2)

	expectNoMatch(ts, "[a-z]", "`")
	expectNoMatch(ts, "[a-z]", "{")

	expectError(ts, "[")
	expectError(ts, "[0")
	expectError(ts, "[0-")
	expectError(ts, "[0-9")
	expectError(ts, "]")
	expectError(ts, "9]")
	expectError(ts, "-9]")
	expectError(ts, "0-9]")
}

func TestRegexSet(ts *testing.T) {
	expectMatch(ts, "_", "\n")
	expectMatch(ts, "a", "a")
	expectMatch(ts, "o", "+")
	expectMatch(ts, "n", "7")
	expectMatch(ts, "Q", "\"")
	expectMatch(ts, "q", "'")

	expectNoMatch(ts, "_", "b")
	expectNoMatch(ts, "a", "4")
	expectNoMatch(ts, "o", "\t")
	expectNoMatch(ts, "n", "|")
	expectNoMatch(ts, "Q", "^")
	expectNoMatch(ts, "q", "&")
}

func TestRegexSeq(ts *testing.T) {
	expectMatch(ts, "{'abc'}", "abc")
	expectMatch(ts, "{'ab'} {'c'}", "abc")
	expectMatch(ts, "{{{{{{'ab'} {'c'}}}}}}", "abc")

	expectError(ts, "{'abc'")
	expectError(ts, "{")
	expectError(ts, "}")
	expectError(ts, "{{{'abc'")
	expectError(ts, "'abc'}}}")
}

func TestRegexPlus(ts *testing.T) {
	expectMatch(ts, "{'abc'}+", "abcabcabc")
	expectMatch(ts, "{'ab'n}+", "ab1ab2ab3")
	expectMatch(ts, "n+n+", "12")

	expectError(ts, "+")
	expectError(ts, "++")
	expectError(ts, "+a")
	expectError(ts, "{}+")
}

func TestRegexStar(ts *testing.T) {
	expectMatch(ts, "{'abc'}*", "abc")
	expectMatch(ts, "{'abc'}*", "")
	expectMatch(ts, "{'ab'n}*", "ab1ab2ab3")
	expectMatch(ts, "{{{'hello'}}}*", "")
	expectMatch(ts, "{{{'hello'}}}*", "hellohellohello")

	expectError(ts, "*")
	expectError(ts, "***")
	expectError(ts, "*a")
	expectError(ts, "{}*")
}

func TestRegexQuest(ts *testing.T) {
	expectMatch(ts, "{'abc'}?", "abc")
	expectMatch(ts, "{'abc'}?", "")
	expectMatch(ts, "{'ab'n}?", "ab1")
	expectMatch(ts, "{{{'hello'}}}?", "")
	expectMatch(ts, "{{{'hello'}}}?", "hello")

	expectError(ts, "?")
	expectError(ts, "???")
	expectError(ts, "?a")
	expectError(ts, "{}?")
}

func TestRegexOr(ts *testing.T) {
	expectMatch(ts, "{'a'|'b'}", "a")
	expectMatch(ts, "{'a'|'b'}", "a")
	expectMatch(ts, "{'a' | 'b'}", "a")
	expectMatch(ts, "{'a' | 'b'}", "b")
	expectMatch(ts, "a{a|'_'|n}*", "snake_case_variable123")

	expectError(ts, "|")
	expectError(ts, "||")
	expectError(ts, "|||")
	expectError(ts, "'a'|{}")
	expectError(ts, "{}|'b'")
	expectError(ts, "'a'|")
	expectError(ts, "|'b'")
}

func TestRegexWave(ts *testing.T) {
	expectMatch(ts, "^~'c'", "abc")
	expectMatch(ts, "a~'z'", "ahjklz")
	expectMatchEq(ts, "'//' {a|' '} ~ '//'", "// The program starts here // int main() {", 29)
	expectMatch(ts, "n ~ {'z'|'9'}", "0123456789")
	expectMatch(ts, "n ~ {'z'|'9'}", "012345678z")
	expectMatch(ts, "{' '} ~ 'sus'", "                           sus               ")
	expectNoMatch(ts, "{' '} ~ 'sus'", "            |             sus               ")

	expectError(ts, "~")
	expectError(ts, "a~")
	expectError(ts, "~{}")
	expectError(ts, "{}~")
}
