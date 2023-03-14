package main

import (
	"fmt"
	"os"
	// "strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println(`No sources specified in the command line arguments`)
		os.Exit(1)
	}

	src, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("Scanner: ", err)
		os.Exit(1)
	}

	sn := NewScanner(string(src), NewBeeSyntax())
	ps := NewParser(sn)
	ast, err := ps.Parse()
	if err != nil {
		fmt.Println("Parser: ", err)
		os.Exit(1)
	}

	for _, node := range ast.Body {
			fmt.Printf("%v\n", node)
	}
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	} else {
// 		return b
// 	}
// }

// func main() {
// 	args := os.Args[1:]

// 	if len(args) < 1 {
// 		fmt.Println(`No sources specified in the command line arguments`)
// 		os.Exit(1)
// 	}

// 	src, err := os.ReadFile(args[0])
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	sn := NewScanner(string(src), NewBeeSyntax())
// 	toks := []Token{}
// 	maxLen := 0

// 	for !sn.Finished() {
// 		toks = append(toks, sn.Tokenize())
// 	}
// 	for _, tok := range toks {
// 		maxLen = max(maxLen, len(tok.Expr))
// 	}
// 	for _, tok := range toks {
// 		fmt.Printf("%*s : %s\n", maxLen, strconv.Quote(tok.Expr), BeeTraitName(tok.Trait))
// 	}
// }

// func main() {
// 	args := os.Args[1:]

// 	switch len(args) {
// 	case 0:
// 		fmt.Printf(`No regex source given to the program`)
// 		os.Exit(1)
// 	case 1:
// 		rx, err := NewRegex(args[0])
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		fmt.Printf(`%s`, rx.Graph(args[0]))

// 	default: // 2 or more arguments (Rest is just ignored)
// 		src, expr := args[0], args[1]

// 		rx, err := NewRegex(src)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		match := rx.Match(expr)
// 		if match < 0 {
// 			fmt.Printf(`"%s" -> "%s" :: no-match :(`, src, expr)
// 		} else {
// 			fmt.Printf(`"%s" -> "%s" :: %s`, src, expr, expr[:match])
// 		}
// 	}
// }

// func main() {
// 	fmt.Println("Bee")
// }
