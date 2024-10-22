package main

import (
	"fmt"
	"log"
)

type Visitor interface {
	VisitNumberExpr(NumberExpr) (float64, error)
	VisitStringExpr(StringExpr) (string, error)
	VisitUnaryExpr(UnaryExpr) (any, error)
	VisitBinaryExpr(BinaryExpr) (any, error)
	VisitBoolExpr(BoolExpr) (any, error)
	VisitIdentifierExpr(IdentifierExpr) (any, error)
	VisitGroupExpr(GroupExpr) (any, error)
	VisitNilExpr(NilExpr) (any, error)
	VisitAssignmentExpr(AssignmentExpr) (any, error)
	VisitPrintStmt(PrintStmt) error
	VisitExpressionStmt(ExpressionStmt) error
	VisitVarDeclStmt(VarDeclStmt) error
	VisitBlockStmt(BlockStmt) error
}

type RuntimeError struct {
	wrapped error
}

type Environment struct {
	values map[string]any
	outer  *Environment
	log    *log.Logger
}

func NewEnvironment(log *log.Logger, outer *Environment) *Environment {
	return &Environment{
		values: make(map[string]any),
		log:    log,
		outer:  outer,
	}
}

func (e Environment) DefineVar(name string, value any) error {
	e.values[name] = value
	return nil
}

func (e Environment) GetVar(name string) (any, error) {
	e.log.Printf("GetVar for '%v' from %v", name, e.values)
	expr, ok := e.values[name]
	if ok {
		return expr, nil
	}
	if e.outer != nil {
		return e.outer.GetVar(name)
	}
	e.log.Printf("unknown variable '%v'", name)
	return nil, RuntimeError{fmt.Errorf("unknown variable '%v'", name)}
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("%v", r.wrapped)
}

type Evaluator struct {
	env *Environment
	log *log.Logger
}

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
	e.log.Printf("VisitBinaryExpr: left:%#v, right:%#v", b.Left, b.Right)
	leftExpr, err := e.EvalExpr(b.Left)
	if err != nil {
		return nil, err
	}
	leftNum, ok := leftExpr.(NumberExpr)
	if ok {
		leftExpr = leftNum.Value
	}
	e.log.Printf("leftExpr: %#v", leftExpr)
	rightExpr, err := e.EvalExpr(b.Right)
	if err != nil {
		return nil, err
	}

	rightNum, ok := rightExpr.(NumberExpr)
	if ok {
		rightExpr = rightNum.Value
	}
	e.log.Printf("rightExpr: %#v", rightExpr)
	switch b.Op.Type {
	case TokenPlus:
		e.log.Printf("isNumber(%v): %v", leftExpr, isNumber(leftExpr))
		e.log.Printf("isNumber(%v): %v", rightExpr, isNumber(rightExpr))
		if isNumber(leftExpr) && isNumber(rightExpr) {
			left, _ := leftExpr.(float64)
			right, _ := rightExpr.(float64)
			return left + right, nil
		} else if isString(leftExpr) && isString(rightExpr) {
			return fmt.Sprintf("%v%v", leftExpr, rightExpr), nil
		}
		e.log.Printf("here")
		return nil, RuntimeError{fmt.Errorf("Both operands must be numbers or strings")}
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
		if isNumber(leftExpr) && isNumber(rightExpr) {
			left, _ := leftExpr.(float64)
			right, _ := rightExpr.(float64)
			return left < right, nil
		}
		return nil, RuntimeError{fmt.Errorf("Operands must be numbers")}
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
		e.log.Printf("unsupported operand: %v", b.Op.Type)
		return nil, RuntimeError{fmt.Errorf("binary expression: unsupported operand '%v'", b.Op.Type)}
	}

}

func (e *Evaluator) VisitIdentifierExpr(b IdentifierExpr) (any, error) {
	e.log.Printf("VisitIdentifierExpr: identifier name: '%v'", b.Value)
	return e.env.GetVar(b.Value)
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
	e.log.Printf("VisitPrintStmt: %v", p)
	v, err := e.EvalExpr(p.Expression)
	if err != nil {
		return err
	}
	if v == nil {
		return RuntimeError{fmt.Errorf("empty expression")}
	}
	fmt.Println(v)
	return nil
}

func (e *Evaluator) VisitBlockStmt(b BlockStmt) error {
	return e.evalBlock(b.Body, NewEnvironment(e.log, e.env))
}

func (e *Evaluator) VisitExpressionStmt(s ExpressionStmt) error {
	_, err := e.EvalExpr(s.Expression)
	if err != nil {
		return err
	}
	return nil
}

func (e *Evaluator) VisitVarDeclStmt(p VarDeclStmt) error {
	value, err := e.EvalExpr(p.Expression)
	if err != nil {
		return err
	}
	return e.env.DefineVar(p.Name, value)
}

func (e *Evaluator) VisitAssignmentExpr(p AssignmentExpr) (any, error) {
	_, err := e.env.GetVar(p.Identifier.Literal)
	if err != nil {
		return nil, err
	}
	value, err := e.EvalExpr(p.Value)
	if err != nil {
		return nil, err
	}
	return value, e.env.DefineVar(p.Identifier.Literal, value)
}

func (e *Evaluator) EvalExpr(expr Expression) (any, error) {
	return expr.accept(e)
}

func (e *Evaluator) Eval(block BlockStmt) error {
	for i := 0; i < len(block.Body); i++ {
		s := block.Body[i]
		e.log.Printf("Evaluating statement: %v", s)
		if err := s.accept(e); err != nil {
			return err
		}
		e.log.Println("-------------")
	}
	return nil
}

func (e *Evaluator) evalBlock(statments []Statement, env *Environment) error {
	parentEnv := e.env
	defer func() {
		e.env = parentEnv
	}()
	e.env = env
	for _, s := range statments {
		if err := s.accept(e); err != nil {
			return err
		}
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
