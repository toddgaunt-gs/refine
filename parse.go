package refine

import (
	"errors"
	"fmt"
	"strconv"
)

type unaryOperator int

const (
	unaryMinus unaryOperator = iota
	unaryPlus
)

type binaryOperator int

const (
	binaryMinus binaryOperator = iota
	binaryPlus
	binaryMultiply
	binaryDivide
	binaryEqual
	binaryLessThan
	binaryLessThanOrEqual
	binaryGreaterThan
	binaryGreaterThanOrEqual
	binaryLogicalAnd
	binaryLogicalOr
)

type atom int

const (
	atomSymbol atom = iota
	atomInteger
	atomString
)

type boxKind int

const (
	boxNil boxKind = iota
	boxString
	boxInteger
)

type box struct {
	kind boxKind
	val  any
}

func evalMultiply(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, errors.New("left and right operands don't share a type!")
	}

	switch left.kind {
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  left.val.(int) * right.val.(int),
		}, nil
	case boxNil:
		//TODO
		return nil, errors.New("invalid operator for type!")
	}
	panic("foo!")
}

func evalAdd(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, errors.New("left and right operands don't share a type!")
	}

	switch left.kind {
	case boxString:
		return &box{
			kind: boxString,
			val:  left.val.(string) + right.val.(string),
		}, nil
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  left.val.(int) + right.val.(int),
		}, nil
	case boxNil:
		//TODO
		return nil, errors.New("invalid operator for type!")
	}
	panic("foo!")
}

func evalSubtract(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, errors.New("left and right operands don't share a type!")
	}

	switch left.kind {
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  left.val.(int) - right.val.(int),
		}, nil
	case boxNil:
		//TODO
		return nil, errors.New("invalid operator for type!")
	}
	panic("foo!")
}

func evalLessThan(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, errors.New("left and right operands don't share a type!")
	}

	switch left.kind {
	case boxString:
		return &box{
			kind: boxString,
			val:  left.val.(string) < right.val.(string),
		}, nil
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  left.val.(int) < right.val.(int),
		}, nil
	case boxNil:
		//TODO
		return nil, errors.New("invalid type!")
	}
	panic("foo!")
}

func evalGreaterThan(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, errors.New("left and right operands don't share a type!")
	}

	switch left.kind {
	case boxString:
		return &box{
			kind: boxString,
			val:  left.val.(string) > right.val.(string),
		}, nil
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  left.val.(int) > right.val.(int),
		}, nil
	case boxNil:
		//TODO
		return nil, errors.New("invalid type!")
	}
	panic("foo!")
}

type expression interface {
	Eval() (*box, error)
}

type unaryExpression struct {
	op   unaryOperator
	expr expression
}

func (u *unaryExpression) Eval() (*box, error) {
	switch u.op {
	case unaryMinus:
	case unaryPlus:
	default:
	}
	panic("unknown unary operator")
}

type binaryExpression struct {
	op    binaryOperator
	left  expression
	right expression
}

func (b *binaryExpression) Eval() (*box, error) {
	left, err := b.left.Eval()
	if err != nil {
		return nil, err
	}

	right, err := b.right.Eval()
	if err != nil {
		return nil, err
	}

	switch b.op {
	case binaryMultiply:
		return evalMultiply(left, right)
	case binaryPlus:
		return evalAdd(left, right)
	case binaryMinus:
		return evalSubtract(left, right)
	case binaryLessThan:
		return evalLessThan(left, right)
	case binaryGreaterThan:
		return evalGreaterThan(left, right)
	}
	panic("unknown binary operator!")
}

type AtomExpression struct {
	kind atom
	text string
}

func (a *AtomExpression) Eval() (*box, error) {
	switch a.kind {
	case atomSymbol:
		panic("not implemented!")
	case atomInteger:
		i, err := strconv.Atoi(a.text)
		if err != nil {
			panic("couldn't convert integer token to integer value!")
		}
		return &box{
			kind: boxInteger,
			val:  i,
		}, nil
	case atomString:
		return &box{
			kind: boxString,
			val:  a.text,
		}, nil
	}
	panic("unkown atom!")
}

type parser struct {
	last   token
	tok    token
	tokens chan token
}

func (p *parser) accept(kind tokenKind) bool {
	if p.tok.kind == kind {
		p.last = p.tok
		p.tok = <-p.tokens
		return true
	} else {
		return false
	}
}

func parseAtom(p *parser) (expression, error) {
	if p.accept(tokenSymbol) {
		return &AtomExpression{
			kind: atomSymbol,
			text: p.last.text,
		}, nil
	}

	if p.accept(tokenInteger) {
		return &AtomExpression{
			kind: atomInteger,
			text: p.last.text,
		}, nil
	}

	return nil, errors.New("couldn't parse as atom")
}

func parseUnary(p *parser) (expression, error) {
	if p.accept(tokenMinus) {
		expr, err := parseExpression(p)
		return &unaryExpression{
			op:   unaryMinus,
			expr: expr,
		}, err
	}

	if p.accept(tokenPlus) {
		expr, err := parseExpression(p)
		return &unaryExpression{
			op:   unaryPlus,
			expr: expr,
		}, err
	}

	return nil, errors.New("couldn't parse as unary expression")
}

func parseBinaryMultiplicative(p *parser) (expression, error) {
	left, err := parseAtom(p)
	if err != nil {
		return nil, err
	}

	if p.accept(tokenAsterisk) {
		right, err := parseBinaryMultiplicative(p)
		return &binaryExpression{
			op:    binaryMultiply,
			left:  left,
			right: right,
		}, err
	}

	return left, nil
}

func parseBinaryAdditive(p *parser) (expression, error) {
	left, err := parseBinaryMultiplicative(p)
	if err != nil {
		return nil, err
	}

	if p.accept(tokenPlus) {
		fmt.Printf("plus\n")
		right, err := parseBinaryAdditive(p)
		return &binaryExpression{
			op:    binaryPlus,
			left:  left,
			right: right,
		}, err
	}

	if p.accept(tokenMinus) {
		fmt.Printf("minus\n")
		right, err := parseBinaryAdditive(p)
		return &binaryExpression{
			op:    binaryMinus,
			left:  left,
			right: right,
		}, err
	}

	return left, nil
}

func parseBinaryComparative(p *parser) (expression, error) {
	left, err := parseBinaryAdditive(p)
	if err != nil {
		return nil, err
	}

	if p.accept(tokenEqual) {
		right, err := parseBinaryComparative(p)
		return &binaryExpression{
			op:    binaryEqual,
			left:  left,
			right: right,
		}, err
	}

	if p.accept(tokenLessThan) {
		right, err := parseBinaryComparative(p)
		return &binaryExpression{
			op:    binaryLessThan,
			left:  left,
			right: right,
		}, err
	}

	if p.accept(tokenGreaterThan) {
		right, err := parseBinaryComparative(p)
		return &binaryExpression{
			op:    binaryGreaterThan,
			left:  left,
			right: right,
		}, err
	}

	return left, nil
}

func parseBinaryLogicalAnd(p *parser) (expression, error) {
	left, err := parseBinaryComparative(p)
	if err != nil {
		return nil, err
	}

	if p.accept(tokenLogicalAnd) {
		right, err := parseBinaryLogicalAnd(p)
		return &binaryExpression{
			op:    binaryLogicalAnd,
			left:  left,
			right: right,
		}, err
	}

	return left, nil
}

func parseBinaryLogicalOr(p *parser) (expression, error) {
	left, err := parseBinaryLogicalAnd(p)
	if err != nil {
		return nil, err
	}

	if p.accept(tokenLogicalOr) {
		right, err := parseBinaryLogicalOr(p)
		return &binaryExpression{
			op:    binaryLogicalOr,
			left:  left,
			right: right,
		}, err
	}

	return left, nil
}

func parseExpression(p *parser) (expression, error) {
	expr, err := parseBinaryLogicalOr(p)
	if err == nil {
		return expr, nil
	}
	return parseAtom(p)
}

func parse(tokens chan token) (expression, error) {
	p := &parser{
		tok:    <-tokens,
		tokens: tokens,
	}
	return parseExpression(p)
}
