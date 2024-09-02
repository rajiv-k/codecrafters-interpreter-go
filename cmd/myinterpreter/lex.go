package main

import (
	"fmt"
	"os"
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
	TokenIllegal
)

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	switch t.Type {
	case TokenNumber:
		return fmt.Sprintf("NUMBER %v %v", t.Literal, t.Literal)
	case TokenString:
		return fmt.Sprintf("STRING %q %v", t.Literal, t.Literal)
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
		// return fmt.Sprintf("COMMENT %q %q", t.Literal, t.Literal)
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

func (l *Lexer) Peek() byte {
	if l.readPosition >= len(l.input) {
		return byte(TokenEOF)
	}
	return l.input[l.readPosition]
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
		if l.Peek() == '/' {
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
		if l.Peek() == '=' {
			l.readChar()
			tok = Token{Type: TokenEqualEqual, Literal: string("==")}
		} else {
			tok = Token{Type: TokenEqual, Literal: string(l.ch)}
		}
	case '<':
		if l.Peek() == '=' {
			l.readChar()
			tok = Token{Type: TokenLessEqual, Literal: string("<=")}
		} else {
			tok = Token{Type: TokenLess, Literal: string(l.ch)}
		}
	case '>':
		if l.Peek() == '=' {
			l.readChar()
			tok = Token{Type: TokenGreaterEqual, Literal: string(">=")}
		} else {
			tok = Token{Type: TokenGreater, Literal: string(l.ch)}
		}
	case '!':
		if l.Peek() == '=' {
			l.readChar()
			tok = Token{Type: TokenBangEqual, Literal: string("!=")}
		} else {
			tok = Token{Type: TokenBang, Literal: string(l.ch)}
		}
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case 0:
		tok.Type = TokenEOF
	default:
		if isDigit(l.ch) {
			tok.Type = TokenNumber
			tok.Literal = l.readNumber()
		} else if (l.ch > 0x41 && l.ch <= 0x5a) || (l.ch >= 0x61 && l.ch <= 0x7a) || l.ch == '_' {
			// TODO(rajiv): lex identifier
		} else {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", l.lineNum, l.ch)
			tok = Token{Type: TokenIllegal, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
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

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' {
		if l.ch == '\n' {
			l.lineNum++
		}
		l.readChar()
	}
}
