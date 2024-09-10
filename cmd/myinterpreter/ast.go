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
	accept(Visitor) any
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
func (n NumberExpr) accept(v Visitor) any {
	return v.VisitNumberExpr(n)
}

// String
type StringExpr struct {
	Value string
}

func (s StringExpr) expr() {}
func (s StringExpr) String() string {
	return s.Value
}
func (s StringExpr) accept(v Visitor) any {
	return v.VisitStringExpr(s)
}

// Identifier
type IdentifierExpr struct {
	Value string
}

func (i IdentifierExpr) expr() {}
func (i IdentifierExpr) String() string {
	return i.Value
}
func (i IdentifierExpr) accept(v Visitor) any {
	return v.VisitIdentifierExpr(i)
}

// Unary expression

type UnaryExpr struct {
	Op      Token
	Operand Expression
}

func (u UnaryExpr) expr() {}
func (u UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", u.Op.Lit(), u.Operand)
}
func (u UnaryExpr) accept(v Visitor) any {
	return v.VisitUnaryExpr(u)
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
func (b BinaryExpr) accept(v Visitor) any {
	return v.VisitBinaryExpr(b)
}

type BoolExpr struct {
	Value bool
}

func (b BoolExpr) expr() {}
func (b BoolExpr) String() string {
	return fmt.Sprintf("%v", b.Value)
}
func (b BoolExpr) accept(v Visitor) any {
	return v.VisitBoolExpr(b)
}

type NilExpr struct{}

func (n NilExpr) expr() {}
func (n NilExpr) String() string {
	return "nil"
}
func (n NilExpr) accept(v Visitor) any {
	return "nil"
}

type GroupExpr struct {
	Expression Expression
}

func (g GroupExpr) expr() {}
func (g GroupExpr) String() string {
	return fmt.Sprintf("(group %v)", g.Expression)
}
func (g GroupExpr) accept(v Visitor) any {
	return v.VisitGroupExpr(g)
}
