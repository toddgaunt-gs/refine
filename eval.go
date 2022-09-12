package refine

import (
	"errors"
	"fmt"
)

type kind int

const (
	boxUntypedNilConstant kind = iota
	boxBool
	boxInt
	boxUint
	boxFloat32
	boxFloat64
	boxString
	boxSlice
	boxMap
	boxPointer
	boxStruct
)

type box struct {
	kind kind
	val  any
}

type evaluator struct {
	symbols map[string]box
	Result  box
	Err     error
}

func newEvaluator() *evaluator {
	ev := &evaluator{
		symbols: map[string]box{},
	}

	// Populate constant values

	ev.symbols["nil"] = box{
		kind: boxUntypedNilConstant,
		val:  nil,
	}

	return ev
}

func evalMultiply(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}
	var v any
	switch left.kind {
	case boxInt:
		v = left.val.(int) * right.val.(int)
	case boxFloat32:
		v = left.val.(float32) * right.val.(float32)
	case boxFloat64:
		v = left.val.(float64) * right.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: left.kind, val: v}, nil
}

func evalDivide(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v any
	switch left.kind {
	case boxInt:
		v = left.val.(int) / right.val.(int)
	case boxFloat32:
		v = left.val.(float32) / right.val.(float32)
	case boxFloat64:
		v = left.val.(float64) / right.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: left.kind, val: v}, nil
}

func evalAdd(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v any
	switch left.kind {
	case boxString:
		v = left.val.(string) + right.val.(string)
	case boxInt:
		v = left.val.(int) + right.val.(int)
	case boxFloat32:
		v = left.val.(float32) + right.val.(float32)
	case boxFloat64:
		v = left.val.(float64) + right.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: left.kind, val: v}, nil
}

func evalSubtract(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v any
	switch left.kind {
	case boxInt:
		v = left.val.(int) - right.val.(int)
	case boxFloat32:
		v = left.val.(float32) - right.val.(float32)
	case boxFloat64:
		v = left.val.(float64) - right.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: left.kind, val: v}, nil
}

func evalLeftShift(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v any
	switch left.kind {
	case boxInt:
		v = left.val.(int) << right.val.(int)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: left.kind, val: v}, nil
}

func evalRightShift(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v any
	switch left.kind {
	case boxInt:
		v = left.val.(int) >> right.val.(int)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: left.kind, val: v}, nil
}

func evalEqual(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v bool
	switch left.kind {
	case boxString:
		v = left.val.(string) == right.val.(string)
	case boxInt:
		v = left.val.(int) == right.val.(int)
	case boxBool:
		v = left.val.(bool) == right.val.(bool)
	case boxSlice:
		v = left.val == right.val
	case boxMap:
		v = left.val == right.val
	case boxPointer:
		v = left.val == right.val
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{
		kind: boxBool,
		val:  v,
	}, nil
}

func evalNotEqual(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v bool
	switch left.kind {
	case boxString:
		v = left.val.(string) != right.val.(string)
	case boxInt:
		v = left.val.(int) != right.val.(int)
	case boxBool:
		v = left.val.(bool) != right.val.(bool)
	case boxSlice:
		v = left.val != right.val
	case boxMap:
		v = left.val != right.val
	case boxPointer:
		v = left.val != right.val
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{
		kind: boxBool,
		val:  v,
	}, nil
}

func evalLessThan(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v bool
	switch left.kind {
	case boxString:
		v = left.val.(string) < right.val.(string)
	case boxInt:
		v = left.val.(int) < right.val.(int)
	case boxFloat32:
		v = left.val.(float32) < right.val.(float32)
	case boxFloat64:
		v = left.val.(float64) < right.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: boxBool, val: v}, nil
}

func evalLessThanOrEqual(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v bool
	switch left.kind {
	case boxString:
		v = left.val.(string) <= right.val.(string)
	case boxInt:
		v = left.val.(int) <= right.val.(int)
	case boxFloat32:
		v = left.val.(float32) <= right.val.(float32)
	case boxFloat64:
		v = left.val.(float64) <= right.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: boxBool, val: v}, nil
}

func evalGreaterThan(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v bool
	switch left.kind {
	case boxString:
		v = left.val.(string) > right.val.(string)
	case boxInt:
		v = left.val.(int) > right.val.(int)
	case boxFloat32:
		v = left.val.(float32) > right.val.(float32)
	case boxFloat64:
		v = left.val.(float64) > right.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: boxBool, val: v}, nil
}

func evalGreaterThanOrEqual(left, right box) (box, error) {
	if left.kind != right.kind {
		return box{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	var v bool
	switch left.kind {
	case boxString:
		v = left.val.(string) >= right.val.(string)
	case boxInt:
		v = left.val.(int) >= right.val.(int)
	case boxFloat32:
		v = left.val.(float32) >= right.val.(float32)
	case boxFloat64:
		v = left.val.(float64) >= right.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: boxBool, val: v}, nil
}

func evalUnaryNot(val box) (box, error) {
	var v any
	switch val.kind {
	case boxBool:
		v = !val.val.(bool)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: val.kind, val: v}, nil
}

func evalUnaryMinus(val box) (box, error) {
	var v any
	switch val.kind {
	case boxInt:
		v = -val.val.(int)
	case boxFloat32:
		v = -val.val.(float32)
	case boxFloat64:
		v = -val.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: val.kind, val: v}, nil
}

func evalUnaryPlus(val box) (box, error) {
	var v any
	switch val.kind {
	case boxInt:
		v = +val.val.(int)
	case boxFloat32:
		v = +val.val.(float32)
	case boxFloat64:
		v = +val.val.(float64)
	default:
		return box{}, errors.New("invalid type!")
	}

	return box{kind: val.kind, val: v}, nil
}

func evalUnaryDereference(expr box) (box, error) {
	switch expr.kind {
	case boxPointer:
		return expr.val.(box), nil
	default:
		return box{}, errors.New("invalid type!")
	}
}

func (e *evaluator) VisitBooleanExpression(be *booleanExpression) {
	e.Result, e.Err = box{kind: boxBool, val: be.value}, nil
}

func (e *evaluator) VisitIntegerExpression(ie *integerExpression) {
	e.Result, e.Err = box{kind: boxInt, val: ie.value}, nil
}

func (e *evaluator) VisitSymbolExpression(se *symbolExpression) {
	var val, ok = e.symbols[se.text]
	if !ok {
		e.Result, e.Err = box{}, fmt.Errorf("refine.eval: couldn't find value for symbol %s", se.text)
	} else {
		e.Result, e.Err = val, nil
	}
}

func (e *evaluator) VisitStringExpression(se *stringExpression) {
	e.Result, e.Err = box{kind: boxString, val: se.text}, nil
}

func (e *evaluator) VisitSelectorExpression(se *selectorExpression) {
	e.Result, e.Err = box{}, errors.New("selectors are unimplemented!")
}

func (e *evaluator) VisitUnaryExpression(ue *unaryExpression) {
	ue.expr.Accept(e)
	if e.Err != nil {
		return
	}

	switch ue.op {
	case unaryMinus:
		e.Result, e.Err = evalUnaryMinus(e.Result)
	case unaryPlus:
		e.Result, e.Err = evalUnaryPlus(e.Result)
	case unaryNot:
		e.Result, e.Err = evalUnaryNot(e.Result)
	default:
		e.Result, e.Err = box{}, fmt.Errorf("refine.eval: unknown unary operator: %d!", ue.op)
	}
}

func (e *evaluator) VisitBinaryExpression(be *binaryExpression) {
	be.left.Accept(e)
	if e.Err != nil {
		return
	}
	left := e.Result

	be.right.Accept(e)
	if e.Err != nil {
		return
	}
	right := e.Result

	// Coerce nil constants from the left.
	if left.kind == boxUntypedNilConstant && (right.kind == boxSlice || right.kind == boxMap || right.kind == boxPointer) {
		left.kind = right.kind
	}

	// Coerce nil constants from the right.
	if right.kind == boxUntypedNilConstant && (left.kind == boxSlice || left.kind == boxMap || left.kind == boxPointer) {
		right.kind = left.kind
	}

	switch be.op {
	case binaryMultiply:
		e.Result, e.Err = evalMultiply(left, right)
	case binaryDivide:
		e.Result, e.Err = evalDivide(left, right)
	case binaryPlus:
		e.Result, e.Err = evalAdd(left, right)
	case binaryMinus:
		e.Result, e.Err = evalSubtract(left, right)
	case binaryLeftShift:
		e.Result, e.Err = evalLeftShift(left, right)
	case binaryRightShift:
		e.Result, e.Err = evalRightShift(left, right)
	case binaryEqual:
		e.Result, e.Err = evalEqual(left, right)
	case binaryNotEqual:
		e.Result, e.Err = evalNotEqual(left, right)
	case binaryLessThan:
		e.Result, e.Err = evalLessThan(left, right)
	case binaryLessThanOrEqual:
		e.Result, e.Err = evalLessThanOrEqual(left, right)
	case binaryGreaterThan:
		e.Result, e.Err = evalGreaterThan(left, right)
	case binaryGreaterThanOrEqual:
		e.Result, e.Err = evalGreaterThanOrEqual(left, right)
	default:
		e.Result, e.Err = box{}, fmt.Errorf("refine.eval: unknown binary operator: %d!", be.op)
	}
}

// eval is a wrapper around passing the evaluator as a visitor to expr.
func eval(e *evaluator, expr expression) (box, error) {
	expr.Accept(e)
	return e.Result, e.Err
}
