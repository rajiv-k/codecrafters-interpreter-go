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
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	eof := Token{Type: TokenEOF}
	// fmt.Printf("source: %v\n", string(fileContents))

	switch command {
	case "tokenize":
		var foundIllegalToken bool
		if len(fileContents) > 0 {
			lexer := NewLexer(string(fileContents))
			for tok := lexer.Next(); tok.Type != TokenEOF; tok = lexer.Next() {
				if tok.Type == TokenIllegal {
					foundIllegalToken = true
				} else if tok.Type != TokenComment {
					fmt.Printf("%+v\n", tok)
				}
			}
		}
		fmt.Printf("%+v\n", eof)
		if foundIllegalToken {
			os.Exit(65)
		}
	case "parse":
		tokens := make([]Token, 0)
		if len(fileContents) > 0 {
			lexer := NewLexer(string(fileContents))
			for tok := lexer.Next(); tok.Type != TokenEOF; tok = lexer.Next() {
				if tok.Type == TokenIllegal {
					os.Exit(65)
				}
				tokens = append(tokens, tok)
			}
		}

		tokens = append(tokens, eof)
		parser := NewParser(tokens)
		block := parser.Parse(tokens)
		fmt.Printf("%v\n", block)
	default:
		fmt.Printf("invalid command: %v\n", command)
		os.Exit(1)
	}
}
