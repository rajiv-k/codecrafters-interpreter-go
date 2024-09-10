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
		if boolVal, ok := val.(bool); ok {
			return !boolVal
		}
		if strVal, ok := val.(string); ok && strVal == "nil" {
			return true
		}
		if _, ok := val.(float64); ok {
			return false
		}
		return false
	default:
		panic(fmt.Sprintf("unary expression: unsupported operand '%v'", u.Op.Type))
	}
}

func (e *Evaluator) VisitBinaryExpr(b BinaryExpr) any {
	leftOpaque := e.Eval(b.Left)
	rightOpaque := e.Eval(b.Right)
	switch b.Op.Type {
	case TokenPlus:
		if isNumber(leftOpaque) && isNumber(rightOpaque) {
			left, _ := leftOpaque.(float64)
			right, _ := rightOpaque.(float64)
			return left + right
		} else if isString(leftOpaque) && isString(rightOpaque) {
			return fmt.Sprintf("%v%v", leftOpaque, rightOpaque)
		}
		return nil
	case TokenMinus:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left - right
	case TokenStar:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left * right
	case TokenSlash:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left / right
	case TokenLess:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left < right
	case TokenLessEqual:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left < right
	case TokenGreater:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left > right
	case TokenGreaterEqual:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left >= right

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

func isString(a any) bool {
	_, ok := a.(string)
	return ok
}

func isNumber(a any) bool {
	_, ok := a.(float64)
	return ok
}
