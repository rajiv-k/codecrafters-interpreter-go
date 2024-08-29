package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	eof := Token{Type: TokenEOF}
	var foundIllegalToken bool
	if len(fileContents) > 0 {
		lexer := NewLexer(string(fileContents))
		for tok := lexer.Next(); tok.Type != TokenEOF; tok = lexer.Next() {
			if tok.Type == TokenIllegal {
				foundIllegalToken = true
			} else {
				fmt.Printf("%+v\n", tok)
			}
		}
	}
	fmt.Printf("%+v\n", eof)
	if foundIllegalToken {
		os.Exit(65)
	}
}
