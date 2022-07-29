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

type evaluator struct {
	symbols map[string]Value
}

func eval(e *evaluator, expr expression) (Value, error) {
	switch t := expr.(type) {
	case *stringExpression:
		return evalStringExpression(e, t)
	case *symbolExpression:
		return evalSymbolExpression(e, t)
	case *integerExpression:
		return evalIntegerExpression(e, t)
	case *unaryExpression:
		return evalUnaryExpression(e, t)
	case *binaryExpression:
		return evalBinaryExpression(e, t)
	}

	panic("unknown expression type!")
}

func evalStringExpression(e *evaluator, se *stringExpression) (Value, error) {
	return Value{kind: boxString, val: se.text}, nil
}

func evalIntegerExpression(e *evaluator, ie *integerExpression) (Value, error) {
	i, err := strconv.Atoi(ie.text)
	if err != nil {
		panic("couldn't convert integer token to integer value!")
	}

	return Value{kind: boxInt, val: i}, nil
}

func evalSymbolExpression(e *evaluator, se *symbolExpression) (Value, error) {
	var val, ok = e.symbols[se.text]
	if !ok {
		return Value{}, fmt.Errorf("Couldn't find value for symbol %s", se.text)
	}

	return val, nil
}

func evalUnaryExpression(e *evaluator, u *unaryExpression) (Value, error) {
	val, err := eval(e, u.expr)
	if err != nil {
		return Value{}, err
	}

	switch u.op {
	case unaryMinus:
		return evalUnaryMinus(val)
	case unaryPlus:
		return evalUnaryPlus(val)
	case unaryNot:
		return evalUnaryNot(val)
	}

	panic(fmt.Sprintf("unknown unary operator: %d!", u.op))
}

func evalBinaryExpression(e *evaluator, b *binaryExpression) (Value, error) {
	left, err := eval(e, b.left)
	if err != nil {
		return Value{}, err
	}

	right, err := eval(e, b.right)
	if err != nil {
		return Value{}, err
	}

	// Coerce nil constants.
	if left.kind == boxUntypedNilConstant && (right.kind == boxSlice || right.kind == boxMap || right.kind == boxPointer) {
		left.kind = right.kind
	}

	// Coerce nil constants.
	if right.kind == boxUntypedNilConstant && (left.kind == boxSlice || left.kind == boxMap || left.kind == boxPointer) {
		right.kind = left.kind
	}

	switch b.op {
	case binaryMultiply:
		return evalMultiply(left, right)
	case binaryDivide:
		return evalDivide(left, right)
	case binaryPlus:
		return evalAdd(left, right)
	case binaryMinus:
		return evalSubtract(left, right)
	case binaryEqual:
		return evalEqual(left, right)
	case binaryNotEqual:
		return evalNotEqual(left, right)
	case binaryLessThan:
		return evalLessThan(left, right)
	case binaryLessThanOrEqual:
		return evalLessThanOrEqual(left, right)
	case binaryGreaterThan:
		return evalGreaterThan(left, right)
	case binaryGreaterThanOrEqual:
		return evalGreaterThanOrEqual(left, right)
	}

	panic(fmt.Sprintf("unknown binary operator: %d!", b.op))
}

type Visitor interface {
	Visit(e expression)
}

func (b *binaryExpression) Accept(v Visitor) {
	v.Visit(b)
}

func (u *unaryExpression) Accept(v Visitor) {
	v.Visit(u)
}

func (s *stringExpression) Accept(v Visitor) {
	v.Visit(s)
}

func (s *symbolExpression) Accept(v Visitor) {
	v.Visit(s)
}

func (s *integerExpression) Accept(v Visitor) {
	v.Visit(s)
}

func (e *evaluator) Visit(expr expression) {
	switch t := expr.(type) {
	case *stringExpression:
		evalStringExpression(e, t)
	case *symbolExpression:
		evalSymbolExpression(e, t)
	case *unaryExpression:
		evalUnaryExpression(e, t)
	case *binaryExpression:
		evalBinaryExpression(e, t)
	}

	panic("unknown expression type!")
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
