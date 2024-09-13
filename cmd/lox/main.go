package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sanity-io/litter"
)

var (
	verbose bool
)

func main() {
	tokenizeCmd := flag.NewFlagSet("tokenize", flag.ExitOnError)
	parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)
	evaluateCmd := flag.NewFlagSet("evaluate", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	for _, fs := range []*flag.FlagSet{tokenizeCmd, parseCmd, evaluateCmd, runCmd} {
		fs.BoolVar(&verbose, "verbose", false, "enable verbose mode")
	}
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh COMMAND <filename>")
		os.Exit(1)
	}

	logger := log.Default()

	command := os.Args[1]
	eof := Token{Type: TokenEOF}

	switch command {
	case "tokenize":
		tokenizeCmd.Parse(os.Args[2:])
		if !verbose {
			logger.SetOutput(io.Discard)
		}
		if len(tokenizeCmd.Args()) != 1 {
			logger.Fatal("Usage: ./your_program.sh tokenize <filename>")
		}
		source := fileContents(tokenizeCmd.Args()[0])
		var foundIllegalToken bool
		if len(source) > 0 {
			lexer := NewLexer(string(source))
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
		parseCmd.Parse(os.Args[2:])
		if !verbose {
			logger.SetOutput(io.Discard)
		}
		if len(parseCmd.Args()) != 1 {
			logger.Fatal("Usage: ./your_program.sh parse <filename>")
		}
		source := fileContents(parseCmd.Args()[0])
		tokens := make([]Token, 0)
		if len(source) > 0 {
			lexer := NewLexer(string(source))
			for tok := lexer.Next(); tok.Type != TokenEOF; tok = lexer.Next() {
				if tok.Type == TokenIllegal {
					os.Exit(65)
				}
				tokens = append(tokens, tok)
			}
		}

		tokens = append(tokens, eof)
		parser := NewParser(tokens, logger)
		block := parser.Parse(tokens)
		fmt.Printf("%v\n", block)
		logger.Println(litter.Sdump(block))
	case "evaluate":
		evaluateCmd.Parse(os.Args[2:])
		if !verbose {
			logger.SetOutput(io.Discard)
		}
		if len(evaluateCmd.Args()) != 1 {
			logger.Fatal("Usage: ./your_program.sh evaluate <filename>")
		}
		source := fileContents(evaluateCmd.Args()[0])
		tokens := make([]Token, 0)
		if len(source) > 0 {
			lexer := NewLexer(string(source))
			for tok := lexer.Next(); tok.Type != TokenEOF; tok = lexer.Next() {
				if tok.Type == TokenIllegal {
					os.Exit(65)
				}
				tokens = append(tokens, tok)
			}
		}

		tokens = append(tokens, eof)
		parser := NewParser(tokens, logger)
		expr := parseExpression(parser, Lowest)
		evaluator := Evaluator{log: logger, env: NewEnvironment(logger, nil)}
		result, err := evaluator.EvalExpr(expr)
		if err != nil {
			os.Exit(70)
		}
		fmt.Println(result)
	case "run":
		runCmd.Parse(os.Args[2:])
		if !verbose {
			logger.SetOutput(io.Discard)
		}
		if len(runCmd.Args()) != 1 {
			logger.Fatal("Usage: ./your_program.sh run <filename>")
		}
		source := fileContents(runCmd.Args()[0])
		tokens := make([]Token, 0)
		if len(source) > 0 {
			lexer := NewLexer(string(source))
			for tok := lexer.Next(); tok.Type != TokenEOF; tok = lexer.Next() {
				if tok.Type == TokenIllegal {
					os.Exit(65)
				}
				tokens = append(tokens, tok)
			}
		}

		logger.Printf("--- END of lexing ---")
		tokens = append(tokens, eof)
		parser := NewParser(tokens, logger)
		block := parser.Parse(tokens)
		logger.Printf("--- END of parsing ---")
		evaluator := Evaluator{env: NewEnvironment(logger, nil), log: logger}
		err := evaluator.Eval(block)
		if err != nil {
			os.Exit(70)
		}
	default:
		fmt.Printf("invalid command: %v\n", command)
	}
}

func fileContents(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)

	}
	return string(b)
}
