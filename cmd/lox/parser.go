package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
	errors []error
	log    *log.Logger
}

func NewParser(tokens []Token, log *log.Logger) *Parser {
	initLookups()
	return &Parser{
		tokens: tokens,
		pos:    0,
		log:    log,
	}
}

func (p *Parser) Parse(tokens []Token) BlockStmt {
	p.log.Println("BEGIN Parse")
	body := make([]Statement, 0)
	for p.hasNext() {
		body = append(body, parseStatement(p))
	}
	p.log.Printf("END Parse: parsed %v Statements. body: %v", len(body), body)
	return BlockStmt{Body: body}
}

func (p *Parser) advance() Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) expect(tokenType TokenType) {
	if p.current().Type != tokenType {
		p.log.Printf("expected: %v, got: '%v', exiting...\n", tokenType, p.current().Literal)
		os.Exit(65)
	}
	p.advance()
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
	nudLookup          = NudLookup{}
	ledLookup          = LedLookup{}
	statementLookup    = StatementLookup{}
	bindingPowerLookup = BindingPowerLookup{
		TokenEOF:          Lowest,
		TokenRightParen:   Lowest,
		TokenEqual:        Assignment,
		TokenBangEqual:    Logical,
		TokenEqualEqual:   Logical,
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

func initLookups() {
	led(TokenAnd, parseBinaryExpr)
	led(TokenOr, parseBinaryExpr)
	led(TokenBangEqual, parseBinaryExpr)
	led(TokenEqualEqual, parseBinaryExpr)
	led(TokenEqual, parseAssignmentExpr)

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
	statementLookup = StatementLookup{
		TokenPrint:     parsePrintStmt,
		TokenVar:       parseVarDeclStmt,
		TokenLeftBrace: parseBlockStmt,
	}
}

func parseExpression(p *Parser, bp BindingPower) Expression {
	p.log.Printf("BEGIN parseExpression")
	token := p.current()
	tokenType := token.Type
	if tokenType == TokenSemiColon {
		return nil
	}
	nudFn, ok := nudLookup[tokenType]
	if !ok {
		p.errors = append(p.errors, fmt.Errorf("Expected 'operand', got: '%v'\n", token.Literal))
		p.log.Printf("Expected 'operand', got: '%v'\n", token.Literal)
		os.Exit(65)
	}

	left := nudFn(p)

	// Parse prefix
	nextBindingPower := p.nextTokenBindingPower()

	for p.hasNext() && nextBindingPower > bp {
		nextTokenType := p.current().Type
		if nextTokenType == TokenRightParen || nextTokenType == TokenSemiColon {
			// End of a group expression
			p.log.Printf("{END} parseExpression at level: %v", bp)
			return left
		}
		ledFn, ok := ledLookup[nextTokenType]
		if !ok {
			panic(fmt.Sprintf("no led handler for %v", nextTokenType))
		}
		left = ledFn(p, left, nextBindingPower)
	}

	p.log.Println("END parseExpression")
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
	p.expect(TokenRightParen)
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

func parseAssignmentExpr(p *Parser, left Expression, bp BindingPower) Expression {
	// Assignment operator
	p.advance()
	identifierExpr, ok := left.(IdentifierExpr)
	if !ok {
		os.Exit(65)
	}
	value := parseExpression(p, Lowest)
	return AssignmentExpr{
		Identifier: Token{Literal: identifierExpr.Value, Type: TokenIdentifier},
		Value:      value,
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
		p.log.Printf("parsePrimaryExpr: tokenType: TokenIdentifier")
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
	p.log.Printf("BEGIN parseStatement")
	tokenType := p.current().Type
	stmtFn, ok := statementLookup[tokenType]
	if ok {
		stmt := stmtFn(p)
		p.log.Printf("END parseStatement: parsed stmt: %v", stmt)
		return stmt
	}

	expr := parseExpression(p, Lowest)
	p.log.Printf("parseStatment (found ExpressionStmt): parsed expr: %v", expr)
	p.log.Print("END parseStatement")
	if p.current().Type == TokenSemiColon {
		p.advance()
	}
	return ExpressionStmt{
		Expression: expr,
	}
}

func parsePrintStmt(p *Parser) Statement {
	p.log.Printf("BEGIN parsePrintStmt")
	// print keyword
	p.expect(TokenPrint)

	expr := parseExpression(p, Lowest)
	p.log.Printf("parsePrintStmt: parsed expr: %v\n", expr)
	if expr == nil {
		os.Exit(65)
	}
	p.expect(TokenSemiColon)
	p.log.Println("END parsePrintStmt")
	return PrintStmt{Expression: expr}
}

func parseVarDeclStmt(p *Parser) Statement {
	p.log.Printf("BEGIN parseVarDeclStmt")
	// var keyword
	p.expect(TokenVar)

	varName := p.advance()
	var expr Expression = NilExpr{}
	if p.current().Type == TokenEqual {
		p.advance()
		expr = parseExpression(p, Lowest)
	}
	p.log.Println("END parseVarDeclStmt")
	p.expect(TokenSemiColon)
	return VarDeclStmt{
		Name:       varName.Literal,
		Expression: expr,
	}

}

func parseExpressionStmt(p *Parser) Statement {
	p.log.Printf("BEGIN parseExpressionStmt")
	expr := parseExpression(p, Lowest)
	p.log.Println("END parseExpressionStmt")
	p.expect(TokenSemiColon)
	return ExpressionStmt{Expression: expr}
}

func parseBlockStmt(p *Parser) Statement {
	p.log.Println("BEGIN parseBlockStmt")
	body := make([]Statement, 0)
	p.expect(TokenLeftBrace)
	for p.current().Type != TokenRightBrace {
		body = append(body, parseStatement(p))
	}
	p.log.Println("END parseBlockStmt")
	p.expect(TokenRightBrace)
	return BlockStmt{Body: body}
}
