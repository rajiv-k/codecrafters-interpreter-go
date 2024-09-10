package main

import "fmt"

type Visitor interface {
	VisitNumberExpr(NumberExpr) (float64, error)
	VisitStringExpr(StringExpr) (string, error)
	VisitUnaryExpr(UnaryExpr) (any, error)
	VisitBinaryExpr(BinaryExpr) (any, error)
	VisitBoolExpr(BoolExpr) (any, error)
	VisitIdentifierExpr(IdentifierExpr) (any, error)
	VisitGroupExpr(GroupExpr) (any, error)
	VisitNilExpr(NilExpr) (any, error)
}

type RuntimeError struct {
	wrapped error
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("%v", r.wrapped)
}

type Evaluator struct{}

func (e *Evaluator) VisitNumberExpr(n NumberExpr) (float64, error) {
	return n.Value, nil
}

func (e *Evaluator) VisitStringExpr(n StringExpr) (string, error) {
	return n.Value, nil
}

func (e *Evaluator) VisitUnaryExpr(u UnaryExpr) (any, error) {
	val, err := e.Eval(u.Operand)
	if err != nil {
		return nil, err
	}
	switch u.Op.Type {
	case TokenMinus:
		floatVal := val.(float64)
		return -floatVal, nil
	case TokenBang:
		if boolVal, ok := val.(bool); ok {
			return !boolVal, nil
		}
		if strVal, ok := val.(string); ok && strVal == "nil" {
			return true, nil
		}
		if _, ok := val.(float64); ok {
			return false, nil
		}
		return false, nil
	default:
		return 0, RuntimeError{fmt.Errorf("unary expression: unsupported operator '%v'", u.Op.Literal)}
	}
}

func (e *Evaluator) VisitBinaryExpr(b BinaryExpr) (any, error) {
	leftOpaque, err := e.Eval(b.Left)
	if err != nil {
		return nil, err
	}
	rightOpaque, err := e.Eval(b.Right)
	if err != nil {
		return nil, err
	}
	switch b.Op.Type {
	case TokenPlus:
		if isNumber(leftOpaque) && isNumber(rightOpaque) {
			left, _ := leftOpaque.(float64)
			right, _ := rightOpaque.(float64)
			return left + right, nil
		} else if isString(leftOpaque) && isString(rightOpaque) {
			return fmt.Sprintf("%v%v", leftOpaque, rightOpaque), nil
		}
		return nil, RuntimeError{fmt.Errorf("")}
	case TokenMinus:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left - right, nil
	case TokenStar:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left * right, nil
	case TokenSlash:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left / right, nil
	case TokenLess:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left < right, nil
	case TokenLessEqual:
		if isNumber(leftOpaque) && isNumber(rightOpaque) {
			left, _ := leftOpaque.(float64)
			right, _ := rightOpaque.(float64)
			return left <= right, nil
		}
		return nil, RuntimeError{fmt.Errorf("")}
	case TokenGreater:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left > right, nil
	case TokenGreaterEqual:
		left, _ := leftOpaque.(float64)
		right, _ := rightOpaque.(float64)
		return left >= right, nil
	case TokenEqualEqual:
		if isNumber(leftOpaque) && isNumber(rightOpaque) {
			left, _ := leftOpaque.(float64)
			right, _ := rightOpaque.(float64)
			return left == right, nil
		} else if isString(leftOpaque) && isString(rightOpaque) {
			return leftOpaque == rightOpaque, nil
		}
		return false, nil
	case TokenBangEqual:
		if isNumber(leftOpaque) && isNumber(rightOpaque) {
			left, _ := leftOpaque.(float64)
			right, _ := rightOpaque.(float64)
			return left != right, nil
		} else if isString(leftOpaque) && isString(rightOpaque) {
			return leftOpaque != rightOpaque, nil
		}
		return true, nil
	default:
		return nil, RuntimeError{fmt.Errorf("binary expression: unsupported operand '%v'", b.Op.Type)}
	}

}

func (e *Evaluator) VisitIdentifierExpr(b IdentifierExpr) (any, error) {
	return b.Value, nil
}

func (e *Evaluator) VisitBoolExpr(b BoolExpr) (any, error) {
	return b.Value, nil
}

func (e *Evaluator) VisitGroupExpr(b GroupExpr) (any, error) {
	return e.Eval(b.Expression)
}

func (e *Evaluator) VisitNilExpr(b NilExpr) (any, error) {
	return b.accept(e)
}
func (e *Evaluator) Eval(expr Expression) (any, error) {
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
