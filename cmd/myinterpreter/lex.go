package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenString
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenColon
	TokenSemiColon
	TokenPlus
	TokenMinus
	TokenDot
	TokenSlash
	TokenStar
	TokenEqualEqual
	TokenBangEqual
	TokenLessEqual
	TokenGreaterEqual
	TokenEqual
	TokenBang
	TokenLess
	TokenComment
	TokenGreater
	TokenIdentifier
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFor
	TokenFun
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile
	TokenIllegal
)

var Keywords = map[string]TokenType{
	"and":    TokenAnd,
	"class":  TokenClass,
	"else":   TokenElse,
	"false":  TokenFalse,
	"for":    TokenFor,
	"fun":    TokenFun,
	"if":     TokenIf,
	"nil":    TokenNil,
	"or":     TokenOr,
	"print":  TokenPrint,
	"return": TokenReturn,
	"super":  TokenSuper,
	"this":   TokenThis,
	"true":   TokenTrue,
	"var":    TokenVar,
	"while":  TokenWhile,
}

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	switch t.Type {
	case TokenNumber:
		if floatVal, err := strconv.ParseFloat(t.Literal, 64); err == nil {
			if _, err2 := strconv.Atoi(t.Literal); err2 == nil {
				return fmt.Sprintf("NUMBER %s %.1f", t.Literal, floatVal)
			} else {
				intVal := int64(floatVal)
				if floatVal-float64(intVal) == 0 {
					return fmt.Sprintf("NUMBER %s %v.0", t.Literal, floatVal)
				}
				return fmt.Sprintf("NUMBER %s %v", t.Literal, floatVal)
			}
		} else {
			panic("invalid format for number")
		}
	case TokenString:
		return fmt.Sprintf("STRING \"%v\" %v", t.Literal, t.Literal)
	case TokenLeftBrace:
		return fmt.Sprintf("LEFT_BRACE %v null", t.Literal)
	case TokenRightBrace:
		return fmt.Sprintf("RIGHT_BRACE %v null", t.Literal)
	case TokenLeftParen:
		return fmt.Sprintf("LEFT_PAREN %v null", t.Literal)
	case TokenRightParen:
		return fmt.Sprintf("RIGHT_PAREN %v null", t.Literal)
	case TokenComma:
		return fmt.Sprintf("COMMA %v null", t.Literal)
	case TokenColon:
		return fmt.Sprintf("COLON %v null", t.Literal)
	case TokenSemiColon:
		return fmt.Sprintf("SEMICOLON %v null", t.Literal)
	case TokenPlus:
		return fmt.Sprintf("PLUS %v null", t.Literal)
	case TokenMinus:
		return fmt.Sprintf("MINUS %v null", t.Literal)
	case TokenStar:
		return fmt.Sprintf("STAR %v null", t.Literal)
	case TokenSlash:
		return fmt.Sprintf("SLASH %v null", t.Literal)
	case TokenDot:
		return fmt.Sprintf("DOT %v null", t.Literal)
	case TokenEqual:
		return fmt.Sprintf("EQUAL %v null", t.Literal)
	case TokenEqualEqual:
		return fmt.Sprintf("EQUAL_EQUAL %v null", t.Literal)
	case TokenBangEqual:
		return fmt.Sprintf("BANG_EQUAL %v null", t.Literal)
	case TokenLessEqual:
		return fmt.Sprintf("LESS_EQUAL %v null", t.Literal)
	case TokenGreaterEqual:
		return fmt.Sprintf("GREATER_EQUAL %v null", t.Literal)
	case TokenBang:
		return fmt.Sprintf("BANG %v null", t.Literal)
	case TokenLess:
		return fmt.Sprintf("LESS %v null", t.Literal)
	case TokenGreater:
		return fmt.Sprintf("GREATER %v null", t.Literal)
	case TokenComment:
		// Comment tokens don't need to be printed
	case TokenIdentifier:
		return fmt.Sprintf("IDENTIFIER %v null", t.Literal)
	case TokenAnd:
		return fmt.Sprintf("AND %v null", t.Literal)
	case TokenClass:
		return fmt.Sprintf("CLASS %v null", t.Literal)
	case TokenElse:
		return fmt.Sprintf("ELSE %v null", t.Literal)
	case TokenFalse:
		return fmt.Sprintf("FALSE %v null", t.Literal)
	case TokenFor:
		return fmt.Sprintf("FOR %v null", t.Literal)
	case TokenFun:
		return fmt.Sprintf("FUN %v null", t.Literal)
	case TokenIf:
		return fmt.Sprintf("IF %v null", t.Literal)
	case TokenNil:
		return fmt.Sprintf("NIL %v null", t.Literal)
	case TokenOr:
		return fmt.Sprintf("OR %v null", t.Literal)
	case TokenPrint:
		return fmt.Sprintf("PRINT %v null", t.Literal)
	case TokenReturn:
		return fmt.Sprintf("RETURN %v null", t.Literal)
	case TokenSuper:
		return fmt.Sprintf("SUPER %v null", t.Literal)
	case TokenThis:
		return fmt.Sprintf("THIS %v null", t.Literal)
	case TokenTrue:
		return fmt.Sprintf("TRUE %v null", t.Literal)
	case TokenVar:
		return fmt.Sprintf("VAR %v null", t.Literal)
	case TokenWhile:
		return fmt.Sprintf("WHILE %v null", t.Literal)
	case TokenEOF:
		return fmt.Sprintf("EOF  null")
	case TokenIllegal:
		return fmt.Sprintf("ILLEGAL %v null", t.Literal)
	default:
		return fmt.Sprintf("unknown token <%v>", t.Literal)
	}
	return "unreachable"
}

type Lexer struct {
	input        string
	position     int // current position in input (points to current char)
	readPosition int // current reading position in input (after current char)
	lineNum      int
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, lineNum: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// func (l *Lexer) next() {
//     if l.readPosition
// }

func (l *Lexer) isAtEnd() bool {
	return l.position >= len(l.input)
}

func (l *Lexer) Peek() byte {
	if l.position >= len(l.input) {
		return byte(TokenEOF)
	}
	return l.input[l.position]
}

func (l *Lexer) PeekNext() byte {
	if l.readPosition >= len(l.input) {
		return byte(TokenEOF)
	}
	return l.input[l.readPosition]
}

func (l *Lexer) backup() {
	l.readPosition = l.position
	if l.position < len(l.input) {
		l.ch = l.input[l.position]
	}
}

func (l *Lexer) Next() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = Token{Type: TokenLeftParen, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: TokenRightParen, Literal: string(l.ch)}
	case '{':
		tok = Token{Type: TokenLeftBrace, Literal: string(l.ch)}
	case '}':
		tok = Token{Type: TokenRightBrace, Literal: string(l.ch)}
	case '[':
		tok = Token{Type: TokenLeftParen, Literal: string(l.ch)}
	case ']':
		tok = Token{Type: TokenRightParen, Literal: string(l.ch)}
	case ',':
		tok = Token{Type: TokenComma, Literal: string(l.ch)}
	case ':':
		tok = Token{Type: TokenColon, Literal: string(l.ch)}
	case ';':
		tok = Token{Type: TokenSemiColon, Literal: string(l.ch)}
	case '+':
		tok = Token{Type: TokenPlus, Literal: string(l.ch)}
	case '-':
		tok = Token{Type: TokenMinus, Literal: string(l.ch)}
	case '*':
		tok = Token{Type: TokenStar, Literal: string(l.ch)}
	case '/':
		if l.PeekNext() == '/' {
			l.readChar() // consume the second slash
			l.skipWhitespace()
			tok.Type = TokenComment
			tok.Literal = l.readComment()
		} else {
			tok = Token{Type: TokenSlash, Literal: string(l.ch)}
		}
	case '.':
		tok = Token{Type: TokenDot, Literal: string(l.ch)}
	case '=':
		if l.PeekNext() == '=' {
			l.readChar()
			tok = Token{Type: TokenEqualEqual, Literal: string("==")}
		} else {
			tok = Token{Type: TokenEqual, Literal: string(l.ch)}
		}
	case '<':
		if l.PeekNext() == '=' {
			l.readChar()
			tok = Token{Type: TokenLessEqual, Literal: string("<=")}
		} else {
			tok = Token{Type: TokenLess, Literal: string(l.ch)}
		}
	case '>':
		if l.PeekNext() == '=' {
			l.readChar()
			tok = Token{Type: TokenGreaterEqual, Literal: string(">=")}
		} else {
			tok = Token{Type: TokenGreater, Literal: string(l.ch)}
		}
	case '!':
		if l.PeekNext() == '=' {
			l.readChar()
			tok = Token{Type: TokenBangEqual, Literal: string("!=")}
		} else {
			tok = Token{Type: TokenBang, Literal: string(l.ch)}
		}
	case '"':
		stringVal, err := l.readString()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		} else {
			tok = Token{Type: TokenString, Literal: stringVal}
		}
	case 0:
		tok.Type = TokenEOF
	default:
		if isDigit(l.ch) {
			valString, err := l.readNumber()
			if err != nil {
				fmt.Printf("[BOO] %v\n", err)
				tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
			} else {
				tok = Token{Type: TokenNumber, Literal: valString}
			}
		} else if isAlpha(l.ch) {
			tok.Literal = l.readIdentifier()
			if v, ok := Keywords[tok.Literal]; ok {
				tok.Type = v
			} else {
				tok.Type = TokenIdentifier
			}
		} else {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", l.lineNum, l.ch)
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for !l.isAtEnd() && isAlphaNumeric(l.ch) {
		l.readChar()
	}
	l.backup()
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() (string, error) {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	containsDecimal := false
	if l.ch == '.' {
		containsDecimal = true
		l.readChar()
	}
	foundAtleastOneDigitAfterDecimal := false
	for isDigit(l.ch) {
		foundAtleastOneDigitAfterDecimal = true
		l.readChar()
	}

	if containsDecimal && !foundAtleastOneDigitAfterDecimal {
		return "", errors.New("invalid number")
	}

	// we have advanced one more character ahead of the number
	l.backup()

	return l.input[start:l.position], nil
}

func (l *Lexer) readString() (string, error) {
	start := l.position + 1
	for l.PeekNext() != '"' && !l.isAtEnd() {
		l.readChar()
		if l.ch == '\n' {
			l.lineNum++
		}
	}

	if l.isAtEnd() {
		return "", fmt.Errorf("[line %d] Error: Unterminated string.\n", l.lineNum)
	}

	l.readChar()

	return l.input[start:l.position], nil
}

func (l *Lexer) readComment() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == 0 || l.ch == '\n' {
			l.lineNum++
			break
		}
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlpha(ch byte) bool {
	return (ch >= 0x41 && ch <= 0x5a) || (ch >= 0x61 && ch <= 0x7a) || ch == '_'
}

func isAlphaNumeric(ch byte) bool {
	return isDigit(ch) || isAlpha(ch)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' {
		if l.ch == '\n' {
			l.lineNum++
		}
		l.readChar()
	}
}
