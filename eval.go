package refine

import (
	"errors"
	"fmt"
	"strconv"
)

type Kind int

const (
	boxUntypedNilConstant Kind = iota
	boxBool
	boxInt
	boxUint
	boxFloat32
	boxFloat64
	boxString
	boxSlice
	boxMap
	boxPointer
)

type Value struct {
	kind Kind
	val  any
}

type evaluator struct {
	symbols map[string]Value
	Result  Value
	Err     error
}

func newEvaluator() *evaluator {
	ev := &evaluator{
		symbols: map[string]Value{},
	}

	// Populate constant values

	ev.symbols["nil"] = Value{
		kind: boxUntypedNilConstant,
		val:  nil,
	}

	ev.symbols["true"] = Value{
		kind: boxBool,
		val:  true,
	}

	ev.symbols["false"] = Value{
		kind: boxBool,
		val:  false,
	}

	return ev
}

func evalMultiply(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: left.kind, val: v}, nil
}

func evalDivide(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: left.kind, val: v}, nil
}

func evalAdd(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: left.kind, val: v}, nil
}

func evalSubtract(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: left.kind, val: v}, nil
}

func evalEqual(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{
		kind: boxBool,
		val:  v,
	}, nil
}

func evalNotEqual(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{
		kind: boxBool,
		val:  v,
	}, nil
}

func evalLessThan(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: boxBool, val: v}, nil
}

func evalLessThanOrEqual(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: boxBool, val: v}, nil
}

func evalGreaterThan(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: boxBool, val: v}, nil
}

func evalGreaterThanOrEqual(left, right Value) (Value, error) {
	if left.kind != right.kind {
		return Value{}, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: boxBool, val: v}, nil
}

func evalUnaryNot(val Value) (Value, error) {
	var v any
	switch val.kind {
	case boxBool:
		v = !val.val.(bool)
	default:
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: val.kind, val: v}, nil
}

func evalUnaryMinus(val Value) (Value, error) {
	var v any
	switch val.kind {
	case boxInt:
		v = -val.val.(int)
	case boxFloat32:
		v = -val.val.(float32)
	case boxFloat64:
		v = -val.val.(float64)
	default:
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: val.kind, val: v}, nil
}

func evalUnaryPlus(val Value) (Value, error) {
	var v any
	switch val.kind {
	case boxInt:
		v = +val.val.(int)
	case boxFloat32:
		v = +val.val.(float32)
	case boxFloat64:
		v = +val.val.(float64)
	default:
		return Value{}, errors.New("invalid type!")
	}

	return Value{kind: val.kind, val: v}, nil
}

func evalUnaryDereference(expr Value) (Value, error) {
	switch expr.kind {
	case boxPointer:
		return expr.val.(Value), nil
	default:
		return Value{}, errors.New("invalid type!")
	}
}

func (e *evaluator) VisitIntegerExpression(ie *integerExpression) {
	i, err := strconv.Atoi(ie.text)
	if err != nil {
		panic("couldn't convert integer token to integer value!")
	}

	e.Result, e.Err = Value{kind: boxInt, val: i}, nil
}

func (e *evaluator) VisitSymbolExpression(se *symbolExpression) {
	var val, ok = e.symbols[se.text]
	if !ok {
		e.Result, e.Err = Value{}, fmt.Errorf("Couldn't find value for symbol %s", se.text)
	} else {
		e.Result, e.Err = val, nil
	}
}

func (e *evaluator) VisitStringExpression(se *stringExpression) {
	e.Result, e.Err = Value{kind: boxString, val: se.text}, nil
}

func (e *evaluator) VisitSelectorExpression(se *selectorExpression) {
	e.Result, e.Err = Value{}, errors.New("selectors are unimplemented!")
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
		panic(fmt.Sprintf("unknown unary operator: %d!", ue.op))
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
		panic(fmt.Sprintf("unknown binary operator: %d!", be.op))
	}
}

// eval is a wrapper around passing the evaluator as a visitor to expr.
func eval(e *evaluator, expr expression) (Value, error) {
	expr.Accept(e)
	return e.Result, e.Err
}
