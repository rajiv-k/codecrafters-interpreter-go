package main

import "fmt"

type Visitor interface {
	VisitNumberExpr(NumberExpr) float64
	VisitStringExpr(StringExpr) string
	VisitUnaryExpr(UnaryExpr) any
	VisitBinaryExpr(BinaryExpr) any
	VisitBoolExpr(BoolExpr) bool
	VisitIdentifierExpr(IdentifierExpr) any
	VisitGroupExpr(GroupExpr) any
	VisitNilExpr(NilExpr) any
}

type Evaluator struct{}

func (e *Evaluator) VisitNumberExpr(n NumberExpr) float64 {
	return n.Value
}

func (e *Evaluator) VisitStringExpr(n StringExpr) string {
	return n.Value
}

func (e *Evaluator) VisitUnaryExpr(u UnaryExpr) any {
	val := e.Eval(u.Operand)
	switch u.Op.Type {
	case TokenMinus:
		floatVal := val.(float64)
		return -floatVal
	case TokenBang:
		boolVal := val.(bool)
		return !boolVal
	default:
		panic(fmt.Sprintf("unary expression: unsupported operand '%v'", u.Op.Type))
	}
}

func (e *Evaluator) VisitBinaryExpr(b BinaryExpr) any {
	leftOpaque := e.Eval(b.Left)
	rightOpaque := e.Eval(b.Right)
	left, _ := leftOpaque.(float64)
	right, _ := rightOpaque.(float64)
	switch b.Op.Type {
	case TokenPlus:
		return left + right
	case TokenMinus:
		return left - right
	case TokenStar:
		return left * right
	case TokenSlash:
		return left / right
	default:
		panic(fmt.Sprintf("binary expression: unsupported operand '%v'", b.Op.Type))
	}

}

func (e *Evaluator) VisitIdentifierExpr(b IdentifierExpr) any {
	return b.Value
}

func (e *Evaluator) VisitBoolExpr(b BoolExpr) bool {
	return b.Value
}

func (e *Evaluator) VisitGroupExpr(b GroupExpr) any {
	return e.Eval(b.Expression)
}

func (e *Evaluator) VisitNilExpr(b NilExpr) any {
	return b.accept(e)
}
func (e *Evaluator) Eval(expr Expression) any {
	return expr.accept(e)
}
