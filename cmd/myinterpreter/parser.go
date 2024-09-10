package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	createTokenLookup()
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) Parse(tokens []Token) BlockStmt {
	body := make([]Statement, 0)
	for p.hasNext() {
		body = append(body, parseStatement(p))
	}
	return BlockStmt{Body: body}
}

func (p *Parser) advance() Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{
			Type: TokenEOF,
		}
	}
	return p.tokens[p.pos]
}

func (p *Parser) nextTokenBindingPower() BindingPower {
	if p.pos >= len(p.tokens) {
		return Lowest
	}
	return bindingPowerLookup[p.tokens[p.pos].Type]
}

func (p *Parser) hasNext() bool {
	return p.pos < len(p.tokens) && p.current().Type != TokenEOF
}

type BindingPower int

const (
	Lowest BindingPower = iota
	Comma
	Assignment
	Logical
	Relational
	Additive
	Multiplicative
	Unary
	Call
	Member
	Group
	Primary
)

var bindingPowerStr = map[BindingPower]string{
	Lowest:         "Lowest",
	Comma:          "Comma",
	Assignment:     "Assignment",
	Logical:        "Logical",
	Relational:     "Relational",
	Additive:       "Additive",
	Multiplicative: "Multiplicative",
	Unary:          "Unary",
	Call:           "Call",
	Member:         "Member",
	Group:          "Group",
	Primary:        "Primary",
}

func (bp BindingPower) String() string {
	return fmt.Sprintf("%s(%d)", bindingPowerStr[bp], bp)
}

type StatementHandler func(p *Parser) Statement
type NudHandler func(p *Parser) Expression
type LedHandler func(p *Parser, left Expression, bp BindingPower) Expression

type StatementLookup map[TokenType]StatementHandler
type NudLookup map[TokenType]NudHandler
type LedLookup map[TokenType]LedHandler
type BindingPowerLookup map[TokenType]BindingPower

var (
	statementLookup    = StatementLookup{}
	nudLookup          = NudLookup{}
	ledLookup          = LedLookup{}
	bindingPowerLookup = BindingPowerLookup{
		TokenEOF:          Lowest,
		TokenRightParen:   Lowest,
		TokenEqual:        Assignment,
		TokenTrue:         Logical,
		TokenFalse:        Logical,
		TokenLess:         Relational,
		TokenLessEqual:    Relational,
		TokenGreater:      Relational,
		TokenGreaterEqual: Relational,
		TokenPlus:         Additive,
		TokenMinus:        Additive,
		TokenStar:         Multiplicative,
		TokenSlash:        Multiplicative,
		TokenBang:         Unary,
		TokenLeftParen:    Group,
		TokenNumber:       Primary,
		TokenString:       Primary,
		TokenIdentifier:   Primary,
	}
)

func nud(tokenType TokenType, nudFn NudHandler) {
	nudLookup[tokenType] = nudFn
}

func led(tokenType TokenType, ledFn LedHandler) {
	ledLookup[tokenType] = ledFn
}

func statement(tokenType TokenType, statementFn StatementHandler) {
	bindingPowerLookup[tokenType] = Lowest
	statementLookup[tokenType] = statementFn
}

func createTokenLookup() {
	led(TokenAnd, parseBinaryExpr)
	led(TokenOr, parseBinaryExpr)

	led(TokenLess, parseBinaryExpr)
	led(TokenLessEqual, parseBinaryExpr)
	led(TokenGreater, parseBinaryExpr)
	led(TokenGreaterEqual, parseBinaryExpr)

	led(TokenPlus, parseBinaryExpr)
	led(TokenMinus, parseBinaryExpr)
	led(TokenStar, parseBinaryExpr)
	led(TokenSlash, parseBinaryExpr)

	nud(TokenNumber, parsePrimaryExpr)
	nud(TokenString, parsePrimaryExpr)
	nud(TokenIdentifier, parsePrimaryExpr)
	nud(TokenTrue, parsePrimaryExpr)
	nud(TokenFalse, parsePrimaryExpr)
	nud(TokenLeftParen, parseGroupExpr)
	nud(TokenNil, parsePrimaryExpr)
	nud(TokenMinus, parseUnaryExpr)
	nud(TokenBang, parseUnaryExpr)
}

func parseExpression(p *Parser, bp BindingPower) Expression {
	token := p.current()
	tokenType := token.Type
	nudFn, ok := nudLookup[tokenType]
	if !ok {
		panic(fmt.Sprintf("no nud handler for %v", tokenType))
	}

	left := nudFn(p)

	// Parse prefix
	nextBindingPower := p.nextTokenBindingPower()

	for p.hasNext() && nextBindingPower > bp {
		nextTokenType := p.current().Type
		if nextTokenType == TokenRightParen {
			// End of a group expression
			return left
		}
		ledFn, ok := ledLookup[nextTokenType]
		if !ok {
			panic(fmt.Sprintf("no led handler for %v", nextTokenType))
		}
		left = ledFn(p, left, nextBindingPower)
	}
	return left
}

func parseGroupExpr(p *Parser) Expression {
	// skip past the open paren
	p.advance()

	// parse the contained expression
	expr := GroupExpr{
		Expression: parseExpression(p, Lowest),
	}

	// consume the closing paren
	p.advance()
	return expr
}

func parseUnaryExpr(p *Parser) Expression {
	op := p.advance()
	return UnaryExpr{
		Op:      op,
		Operand: parseExpression(p, Unary),
	}
}

func parseBinaryExpr(p *Parser, left Expression, bp BindingPower) Expression {
	op := p.advance()
	right := parseExpression(p, bindingPowerLookup[op.Type])
	return BinaryExpr{
		Left:  left,
		Op:    op,
		Right: right,
	}
}

func parsePrimaryExpr(p *Parser) Expression {
	currentTokenType := p.current().Type
	switch currentTokenType {
	case TokenNumber:
		num, _ := strconv.ParseFloat(p.advance().Literal, 64)
		return NumberExpr{
			Value: num,
		}
	case TokenString:
		return StringExpr{
			Value: p.advance().Literal,
		}
	case TokenIdentifier:
		return IdentifierExpr{
			Value: p.advance().Literal,
		}
	case TokenTrue, TokenFalse:
		return BoolExpr{
			Value: p.advance().Literal == "true",
		}
	case TokenNil:
		_ = p.advance()
		return NilExpr{}
	default:
		panic(fmt.Sprintf("coul not create primary expr from unexpected token: %v", currentTokenType))
	}
}

func parseStatement(p *Parser) Statement {
	tokenType := p.current().Type
	stmtFn, ok := statementLookup[tokenType]
	if ok {
		return stmtFn(p)
	}

	expr := parseExpression(p, Lowest)
	return ExpressionStmt{
		Expression: expr,
	}
}
