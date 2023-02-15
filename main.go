package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println(`No sources specified in the command line arguments`)
		os.Exit(1)
	}

	tok := Token{}
	src, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sn := NewScanner(string(src), NewBeeSyntax())

	for tok.Trait != End {
		tok = sn.Tokenize()
		fmt.Printf("%20s: %s\n", strconv.Quote(tok.Expr), BeeTraitName(tok.Trait))
	}
}

// func main() {
// 	fmt.Println("Bee - programming language")
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
// 			fmt.Printf(`"%s" -> "%s" : no-match`, src, expr)
// 		} else {
// 			fmt.Printf(`"%s" -> "%s" : %s`, src, expr, expr[:match])
// 		}
// 	}
// }
