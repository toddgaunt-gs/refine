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
	unaryNot
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
	boxBool
	boxInteger
	boxString
)

type box struct {
	kind boxKind
	val  any
}

func evalMultiply(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	switch left.kind {
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  left.val.(int) * right.val.(int),
		}, nil
	default:
		return nil, errors.New("invalid type!")
	}
}

func evalAdd(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
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
	default:
		return nil, errors.New("invalid type!")
	}
}

func evalSubtract(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	switch left.kind {
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  left.val.(int) - right.val.(int),
		}, nil
	default:
		return nil, errors.New("invalid type!")
	}
}

func evalEqual(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	switch left.kind {
	case boxString:
		return &box{
			kind: boxBool,
			val:  left.val.(string) == right.val.(string),
		}, nil
	case boxInteger:
		return &box{
			kind: boxBool,
			val:  left.val.(int) == right.val.(int),
		}, nil
	case boxBool:
		return &box{
			kind: boxBool,
			val:  left.val.(bool) == right.val.(bool),
		}, nil
	default:
		return nil, errors.New("invalid type!")
	}
}

func evalLessThan(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	switch left.kind {
	case boxString:
		return &box{
			kind: boxBool,
			val:  left.val.(string) < right.val.(string),
		}, nil
	case boxInteger:
		return &box{
			kind: boxBool,
			val:  left.val.(int) < right.val.(int),
		}, nil
	default:
		return nil, errors.New("invalid type!")
	}
}

func evalGreaterThan(left, right *box) (*box, error) {
	if left.kind != right.kind {
		return nil, fmt.Errorf("type mismatch: %d != %d", left.kind, right.kind)
	}

	switch left.kind {
	case boxString:
		return &box{
			kind: boxBool,
			val:  left.val.(string) > right.val.(string),
		}, nil
	case boxInteger:
		return &box{
			kind: boxBool,
			val:  left.val.(int) > right.val.(int),
		}, nil
	default:
		return nil, errors.New("invalid type!")
	}
}

func evalUnaryNot(expr *box) (*box, error) {
	switch expr.kind {
	case boxBool:
		return &box{
			kind: boxBool,
			val:  !expr.val.(bool),
		}, nil
	default:
		//TODO
		return nil, errors.New("invalid type!")
	}
}

func evalUnaryMinus(expr *box) (*box, error) {
	switch expr.kind {
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  -expr.val.(int),
		}, nil
	default:
		//TODO
		return nil, errors.New("invalid type!")
	}
}

func evalUnaryPlus(expr *box) (*box, error) {
	switch expr.kind {
	case boxInteger:
		return &box{
			kind: boxInteger,
			val:  +expr.val.(int),
		}, nil
	default:
		//TODO
		return nil, errors.New("invalid type!")
	}
}

type expression interface {
	Eval() (*box, error)
}

type unaryExpression struct {
	op   unaryOperator
	expr expression
}

func (u *unaryExpression) Eval() (*box, error) {
	expr, err := u.expr.Eval()
	if err != nil {
		return nil, err
	}

	switch u.op {
	case unaryMinus:
		return evalUnaryMinus(expr)
	case unaryPlus:
		return evalUnaryPlus(expr)
	case unaryNot:
		return evalUnaryNot(expr)
	}

	panic(fmt.Sprintf("unknown unary operator: %d!", u.op))
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
	case binaryEqual:
		return evalEqual(left, right)
	case binaryLessThan:
		return evalLessThan(left, right)
	case binaryGreaterThan:
		return evalGreaterThan(left, right)
	}

	panic(fmt.Sprintf("unknown binary operator: %d!", b.op))
}

type atomExpression struct {
	kind atom
	text string
}

func (a *atomExpression) Eval() (*box, error) {
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
	if p.accept(tokenLeftParen) {
		expr, err := parseExpression(p)
		if err != nil {
			return nil, err
		}
		if !p.accept(tokenRightParen) {
			return nil, errors.New("expected ')'")
		}
		return expr, nil
	}

	if p.accept(tokenSymbol) {
		return &atomExpression{
			kind: atomSymbol,
			text: p.last.text,
		}, nil
	}

	if p.accept(tokenInteger) {
		return &atomExpression{
			kind: atomInteger,
			text: p.last.text,
		}, nil
	}

	return nil, errors.New("couldn't parse expression!")
}

func parseUnary(p *parser) (expression, error) {
	if p.accept(tokenPlus) {
		expr, err := parseUnary(p)
		return &unaryExpression{
			op:   unaryPlus,
			expr: expr,
		}, err
	}

	if p.accept(tokenMinus) {
		expr, err := parseUnary(p)
		return &unaryExpression{
			op:   unaryMinus,
			expr: expr,
		}, err
	}

	if p.accept(tokenNot) {
		expr, err := parseUnary(p)
		return &unaryExpression{
			op:   unaryNot,
			expr: expr,
		}, err
	}

	return parseAtom(p)
}

func parseBinaryMultiplicative(p *parser) (expression, error) {
	left, err := parseUnary(p)
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

	if p.accept(tokenDivide) {
		right, err := parseBinaryMultiplicative(p)
		return &binaryExpression{
			op:    binaryDivide,
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

// parseExpression is the top-level parsing function starting at the lowest
// precedence level, working its way up the chain of precedence functions.
func parseExpression(p *parser) (expression, error) {
	return parseBinaryLogicalOr(p)
}

func parse(tokens chan token) (expression, error) {
	p := &parser{
		tok:    <-tokens,
		tokens: tokens,
	}
	expr, err := parseExpression(p)
	if err != nil {
		return nil, err
	}

	// Create a list of any unparsed tokens as an error
	var unparsedTokens []token
	for tok := p.tok; tok.kind != tokenEOF; tok = <-tokens {
		if tok.kind == tokenEOF {
			break
		}
		unparsedTokens = append(unparsedTokens, tok)
	}
	if len(unparsedTokens) > 0 {
		return nil, fmt.Errorf("Remaining unparsed tokens: %v", unparsedTokens)
	}
	return expr, err
}
