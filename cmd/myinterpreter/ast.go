package main

import (
	"fmt"
	"strings"
)

type Statement interface {
	fmt.Stringer
	stmt()
}

type BlockStmt struct {
	Body []Statement
}

func (b BlockStmt) String() string {
	sb := &strings.Builder{}
	for _, s := range b.Body {
		fmt.Fprintln(sb, s.String())
	}
	return sb.String()
}

type ExpressionStmt struct {
	Expression Expression
}

func (e ExpressionStmt) stmt() {}
func (e ExpressionStmt) String() string {
	return e.Expression.String()
}

type Expression interface {
	fmt.Stringer
	expr()
}

// Number
type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}
func (n NumberExpr) String() string {
	intVal := int64(n.Value)
	if n.Value-float64(intVal) == 0 {
		return fmt.Sprintf("%v.0", n.Value)
	}
	return fmt.Sprintf("%v", n.Value)
}

// String
type StringExpr struct {
	Value string
}

func (s StringExpr) expr() {}
func (s StringExpr) String() string {
	return s.Value
}

// Identifier
type IdentifierExpr struct {
	Value string
}

func (i IdentifierExpr) expr() {}
func (i IdentifierExpr) String() string {
	return i.Value
}

// Binary expression
type BinaryExpr struct {
	Left  Expression
	Op    Token
	Right Expression
}

func (b BinaryExpr) expr() {}
func (b BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", b.Op.Lit(), b.Left, b.Right)
}

type BoolExpr struct {
	Value bool
}

func (b BoolExpr) expr() {}
func (b BoolExpr) String() string {
	return fmt.Sprintf("%v", b.Value)
}

type NilExpr struct{}

func (n NilExpr) expr() {}
func (n NilExpr) String() string {
	return "nil"
}
