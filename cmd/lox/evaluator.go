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
	VisitPrintStmt(PrintStmt) error
	VisitExpressionStmt(ExpressionStmt) error
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
	val, err := e.EvalExpr(u.Operand)
	if err != nil {
		return nil, err
	}
	switch u.Op.Type {
	case TokenMinus:
		if !isNumber(val) {
			return nil, RuntimeError{fmt.Errorf("Operand must be a number")}
		}
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
	leftExpr, err := e.EvalExpr(b.Left)
	if err != nil {
		return nil, err
	}
	rightExpr, err := e.EvalExpr(b.Right)
	if err != nil {
		return nil, err
	}
	switch b.Op.Type {
	case TokenPlus:
		if isNumber(leftExpr) && isNumber(rightExpr) {
			left, _ := leftExpr.(float64)
			right, _ := rightExpr.(float64)
			return left + right, nil
		} else if isString(leftExpr) && isString(rightExpr) {
			return fmt.Sprintf("%v%v", leftExpr, rightExpr), nil
		}
		return nil, RuntimeError{fmt.Errorf("")}
	case TokenMinus:
		if !isNumber(leftExpr) || !isNumber(rightExpr) {
			return nil, RuntimeError{fmt.Errorf("Operands must be numbers")}
		}
		left, _ := leftExpr.(float64)
		right, _ := rightExpr.(float64)
		return left - right, nil
	case TokenStar:
		if !isNumber(leftExpr) || !isNumber(rightExpr) {
			return nil, RuntimeError{fmt.Errorf("Operands must be numbers")}
		}
		left, _ := leftExpr.(float64)
		right, _ := rightExpr.(float64)
		return left * right, nil
	case TokenSlash:
		if !isNumber(leftExpr) || !isNumber(rightExpr) {
			return nil, RuntimeError{fmt.Errorf("Operands must be numbers")}
		}
		left, _ := leftExpr.(float64)
		right, _ := rightExpr.(float64)
		if right == 0.0 {
			return nil, RuntimeError{fmt.Errorf("Division by 0 is not allowed")}
		}
		return left / right, nil
	case TokenLess:
		if !isNumber(leftExpr) || !isNumber(rightExpr) {
			return nil, RuntimeError{fmt.Errorf("Operands must be numbers")}
		}
		left, _ := leftExpr.(float64)
		right, _ := rightExpr.(float64)
		return left < right, nil
	case TokenLessEqual:
		if !isNumber(leftExpr) || !isNumber(rightExpr) {
			return nil, RuntimeError{fmt.Errorf("Operands must be numbers")}
		}
		left, _ := leftExpr.(float64)
		right, _ := rightExpr.(float64)
		return left <= right, nil
	case TokenGreater:
		if !isNumber(leftExpr) || !isNumber(rightExpr) {
			return nil, RuntimeError{fmt.Errorf("Operands must be numbers")}
		}
		left, _ := leftExpr.(float64)
		right, _ := rightExpr.(float64)
		return left > right, nil
	case TokenGreaterEqual:
		if !isNumber(leftExpr) || !isNumber(rightExpr) {
			return nil, RuntimeError{fmt.Errorf("Operands must be numbers")}
		}
		left, _ := leftExpr.(float64)
		right, _ := rightExpr.(float64)
		return left >= right, nil
	case TokenEqualEqual:
		return leftExpr == rightExpr, nil
	case TokenBangEqual:
		if isNumber(leftExpr) && isNumber(rightExpr) {
			left, _ := leftExpr.(float64)
			right, _ := rightExpr.(float64)
			return left != right, nil
		} else if isString(leftExpr) && isString(rightExpr) {
			return leftExpr != rightExpr, nil
		}
		leftBool := asBool(leftExpr)
		rightBool := asBool(rightExpr)
		return leftBool != rightBool, nil
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
	return e.EvalExpr(b.Expression)
}

func (e *Evaluator) VisitNilExpr(b NilExpr) (any, error) {
	return b.accept(e)
}

func (e *Evaluator) VisitPrintStmt(p PrintStmt) error {
	v, err := e.EvalExpr(p.Expression)
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}
func (e *Evaluator) VisitExpressionStmt(s ExpressionStmt) error {
	expr, err := e.EvalExpr(s.Expression)
	if err != nil {
		return err
	}
	fmt.Printf("expr: %v\n", expr)
	return nil
}

func (e *Evaluator) EvalExpr(expr Expression) (any, error) {
	return expr.accept(e)
}

func (e *Evaluator) Eval(block BlockStmt) error {
	for _, s := range block.Body {
		s.accept(e)
	}
	return nil
}

func isString(a any) bool {
	_, ok := a.(string)
	return ok
}

func isNumber(a any) bool {
	_, ok := a.(float64)
	return ok
}

func asBool(a any) bool {
	if a == nil {
		return false
	}

	if s, ok := a.(string); ok {
		return s != "nil"
	}

	if b, ok := a.(bool); ok {
		return b
	}
	return true
}
