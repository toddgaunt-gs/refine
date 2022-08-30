package refine

import (
	"errors"
	"fmt"
)

// visitor is an interface for a visitor that is meant to traverse an Abstract
// Syntax Tree (AST) for the refine package.
type visitor interface {
	VisitIntegerExpression(i *integerExpression)
	VisitStringExpression(s *stringExpression)
	VisitSymbolExpression(s *symbolExpression)
	VisitSelectorExpression(s *selectorExpression)
	VisitUnaryExpression(u *unaryExpression)
	VisitBinaryExpression(b *binaryExpression)
}

type expression interface {
	// Accept calls the appropriate method of a visitor for a given
	// implementation of expression.
	Accept(v visitor)
}

// Accepts calls a visitor on an integer expression.
func (ie *integerExpression) Accept(v visitor) {
	v.VisitIntegerExpression(ie)
}

// Accepts calls a visitor on a string expression.
func (se *stringExpression) Accept(v visitor) {
	v.VisitStringExpression(se)
}

// Accepts calls a visitor on a symbol expression.
func (se *symbolExpression) Accept(v visitor) {
	v.VisitSymbolExpression(se)
}

// Accepts calls a visitor on a selector expression.
func (se *selectorExpression) Accept(v visitor) {
	v.VisitSelectorExpression(se)
}

// Accepts calls a visitor on a unary expression.
func (ue *unaryExpression) Accept(v visitor) {
	v.VisitUnaryExpression(ue)
}

// Accepts calls a visitor on a binary expression.
func (be *binaryExpression) Accept(v visitor) {
	v.VisitBinaryExpression(be)
}

type integerExpression struct {
	text string
}

type stringExpression struct {
	text string
}

type symbolExpression struct {
	text string
}

type selectorExpression struct {
	sym       *symbolExpression
	selection *symbolExpression
}

type unaryOperator int

const (
	unaryMinus unaryOperator = iota
	unaryPlus
	unaryNot
	unaryDereference
)

type unaryExpression struct {
	op   unaryOperator
	expr expression
}

type binaryOperator int

const (
	binaryMultiply binaryOperator = iota
	binaryDivide

	binaryMinus
	binaryPlus

	binaryEqual
	binaryNotEqual
	binaryLessThan
	binaryLessThanOrEqual
	binaryGreaterThan
	binaryGreaterThanOrEqual

	binaryLogicalOr
	binaryLogicalAnd
)

type binaryExpression struct {
	op    binaryOperator
	left  expression
	right expression
}

type parser struct {
	// last is the token that was last accepted token within a parsing function.
	last token
	// tok is the current token looking to be accepted from a parsing function.
	tok token
	// tokens streams tokens from the lexer to the parser.
	tokens chan token
}

// accept looks for a kind of token waiting in channel. If the token in the
// channel matches the kind provided, it is consumed from the channel.
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
		sym := &symbolExpression{
			text: p.last.text,
		}
		if p.accept(tokenPeriod) {
			if p.accept(tokenSymbol) {
				selection := &symbolExpression{
					text: p.last.text,
				}
				return &selectorExpression{
					sym:       sym,
					selection: selection,
				}, nil
			}
			return nil, errors.New("expected an identifier following selector")
		}
		return sym, nil
	}

	if p.accept(tokenString) {
		return &stringExpression{
			text: p.last.text[1 : len(p.last.text)-1],
		}, nil
	}

	if p.accept(tokenInteger) {
		return &integerExpression{
			text: p.last.text,
		}, nil
	}

	return nil, errors.New("couldn't parse expression!")
}

func parseUnary(p *parser) (expression, error) {
	var accepted = map[tokenKind]unaryOperator{
		tokenPlus:       unaryPlus,
		tokenMinus:      unaryMinus,
		tokenLogicalNot: unaryNot,
	}

	for kind, op := range accepted {
		if p.accept(kind) {
			expr, err := parseUnary(p)
			return &unaryExpression{
				op:   op,
				expr: expr,
			}, err
		}
	}

	return parseAtom(p)
}

func parseBinaryMultiplicative(p *parser) (expression, error) {
	left, err := parseUnary(p)
	if err != nil {
		return nil, err
	}

	var accepted = map[tokenKind]binaryOperator{
		tokenAsterisk: binaryMultiply,
		tokenDivide:   binaryDivide,
	}

	for kind, op := range accepted {
		if p.accept(kind) {
			right, err := parseBinaryMultiplicative(p)
			return &binaryExpression{
				op:    op,
				left:  left,
				right: right,
			}, err
		}
	}

	return left, nil
}

func parseBinaryAdditive(p *parser) (expression, error) {
	left, err := parseBinaryMultiplicative(p)
	if err != nil {
		return nil, err
	}

	var accepted = map[tokenKind]binaryOperator{
		tokenPlus:  binaryPlus,
		tokenMinus: binaryMinus,
	}

	for kind, op := range accepted {
		if p.accept(kind) {
			right, err := parseBinaryAdditive(p)
			return &binaryExpression{
				op:    op,
				left:  left,
				right: right,
			}, err
		}
	}

	return left, nil
}

func parseBinaryComparative(p *parser) (expression, error) {
	left, err := parseBinaryAdditive(p)
	if err != nil {
		return nil, err
	}

	var accepted = map[tokenKind]binaryOperator{
		tokenEqual:              binaryEqual,
		tokenNotEqual:           binaryNotEqual,
		tokenLessThan:           binaryLessThan,
		tokenLessThanOrEqual:    binaryLessThanOrEqual,
		tokenGreaterThan:        binaryGreaterThan,
		tokenGreaterThanOrEqual: binaryGreaterThanOrEqual,
	}

	for kind, op := range accepted {
		if p.accept(kind) {
			right, err := parseBinaryComparative(p)
			return &binaryExpression{
				op:    op,
				left:  left,
				right: right,
			}, err
		}
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
// precedence level, working its way up the chain of functions according to
// precedence of operations.
func parseExpression(p *parser) (expression, error) {
	/*
		if p.accept(tokenQuestionMark) {
			return &aliasExpression{
				text.p.last.text,
			}, nil
		}
	*/
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
		return nil, fmt.Errorf("Unparsed tokens: %v", unparsedTokens)
	}

	return expr, err
}
