package refine

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokenKind int

const (
	tokenError tokenKind = iota - 1
	tokenEOF
	// Syntax
	tokenPeriod
	tokenComma
	tokenLeftParen
	tokenRightParen
	// Operators
	tokenLogicalOr
	tokenLogicalAnd
	tokenEqual
	tokenLessThan
	tokenLessThanOrEqual
	tokenGreaterThan
	tokenGreaterThanOrEqual
	tokenLogicalNot
	tokenNotEqual
	tokenMinus
	tokenPlus
	tokenAsterisk
	tokenDivide
	tokenBitwiseOr
	tokenBitwiseAnd
	tokenLeftShift
	tokenRightShift
	// Atoms
	tokenInteger
	tokenString
	tokenSymbol
)

const eof = 0
const whitespace = " \t\r\v\n"
const delimiters = string(rune(eof)) + whitespace + "()=<>+-*/"

type token struct {
	kind tokenKind
	text string
}

type lexer struct {
	name  string // Name of the lexer to identify it when errors occur.
	input string // String being lexed.

	start int // Start position of the current token.
	index int // Current index in the input.
	width int // Width of the last read rune.

	tokens chan token // Channel scanned tokens are output to.
}

type stateFunc func(l *lexer) stateFunc

// next checks if the next rune is part of the valid set provided, but does not
// consume it.
func (l *lexer) next(valid string) bool {
	if strings.ContainsRune(valid, l.get()) {
		l.unget()
		return true
	}
	l.unget()
	return false
}

// accept consumes a rune from the lexer if it is part of the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.get()) {
		return true
	}
	l.unget()
	return false
}

// acceptRun consumes runes from the valid set of runes provided until none in
// the set are encountered.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.get()) {
	}
	l.unget()
}

// ignore skips over runes read so far so they are not emitted.
func (l *lexer) ignore() {
	l.start = l.index
}

// unget moves the lexer back by one rune in the input. This can only be
// performed once for every call to l.get()
func (l *lexer) unget() {
	l.index -= l.width
	l.width = 0
}

// peek returns the next rune and consumes it from the input.
func (l *lexer) get() rune {
	if l.index >= len(l.input) {
		l.width = 0
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.index:])
	l.index += width
	l.width = width
	return r
}

// peek returns the next rune without consuming it from the input.
func (l *lexer) peek() rune {
	r := l.get()
	l.unget()
	return r
}

// errorf sends an error token with a formatted error message to the tokens
// channel, and returns nil to end the state machine.
func (l *lexer) errorf(format string, args ...any) stateFunc {
	return func(l *lexer) stateFunc {
		l.tokens <- token{
			kind: tokenError,
			text: fmt.Sprintf(format, args...),
		}
		return nil
	}
}

func lexExclamation(l *lexer) stateFunc {
	if l.accept("!") {
		if l.accept("=") {
			l.emit(tokenNotEqual)
		} else {
			l.emit(tokenLogicalNot)
		}
		return lexStart
	} else {
		return l.errorf("expected '!'")
	}
}

func lexEqual(l *lexer) stateFunc {
	if l.accept("=") {
		if l.accept("=") {
			l.emit(tokenEqual)
		} else {
			return l.errorf("expected '=='")
		}
		return lexStart
	} else {
		return l.errorf("expected '='")
	}
}

func lexLessThan(l *lexer) stateFunc {
	if l.accept("<") {
		if l.accept("=") {
			l.emit(tokenLessThanOrEqual)
		} else if l.accept("<") {
			l.emit(tokenLeftShift)
		} else {
			l.emit(tokenLessThan)
		}
		return lexStart
	} else {
		return l.errorf("expected '<'")
	}
}

func lexGreaterThan(l *lexer) stateFunc {
	if l.accept(">") {
		if l.accept("=") {
			l.emit(tokenGreaterThanOrEqual)
		} else if l.accept(">") {
			l.emit(tokenRightShift)
		} else {
			l.emit(tokenGreaterThan)
		}
		return lexStart
	} else {
		return l.errorf("expected '>'")
	}
}

func lexNumber(l *lexer) stateFunc {
	if l.accept("123456789") {
		for {
			l.accept("_")
			if !l.accept("0123456789") {
				break
			}
		}
		l.unget()
		l.emit(tokenInteger)
		return lexStart
	}

	// Numbers prefixed with 0 might be octal, binary, hex, or just plain 0.
	if l.accept("0") {
		if l.accept("bB") {
			for {
				l.accept("_")
				if !l.accept("01") {
					break
				}
			}
			l.unget()
			l.emit(tokenInteger)
			return lexStart
		}
		if l.accept("hH") {
			for {
				l.accept("_")
				if !l.accept("0123456789abcdefABCDEF") {
					break
				}
			}
			l.unget()
			l.emit(tokenInteger)
			return lexStart
		}
		// For octal numbers, the o rune is optional.
		l.accept("oO")
		for {
			l.accept("_")
			if !l.accept("01234567") {
				break
			}
		}
		l.unget()
		l.emit(tokenInteger)
		return lexStart
	}

	return l.errorf("expected an ASCII digit")
}

func lexSymbol(l *lexer) stateFunc {
	for r := l.get(); unicode.IsLetter(r) || unicode.IsDigit(r); r = l.get() {
		// intentionally empty.
	}
	l.unget()
	if l.index != l.start {
		l.emit(tokenSymbol)
		return lexStart
	} else {
		return l.errorf("invalid symbol")
	}
}

func lexString(l *lexer) stateFunc {
	if l.accept("`") {
		for r := l.get(); r != '`'; r = l.get() {
			if r == eof {
				return l.errorf("reached EOF when reading string")
			}
		}
		l.emit(tokenString)
		return lexStart
	} else {
		return l.errorf("expected '`'")
	}
}

func lexAmpersand(l *lexer) stateFunc {
	if l.accept("&") {
		if l.accept("&") {
			l.emit(tokenLogicalAnd)
		} else {
			l.emit(tokenBitwiseAnd)
		}
		return lexStart
	} else {
		return l.errorf("expected '&'")
	}
}

func lexPipe(l *lexer) stateFunc {
	if l.accept("|") {
		if l.accept("|") {
			l.emit(tokenLogicalOr)
		} else {
			l.emit(tokenBitwiseOr)
		}
		return lexStart
	} else {
		return l.errorf("expected '|'")
	}
}

func lexStart(l *lexer) stateFunc {
	for {
		// Skip leading whitespace.
		l.acceptRun(whitespace)
		l.ignore()

		switch {
		case l.next("."):
			l.accept(".")
			l.emit(tokenPeriod)
			return lexStart
		case l.next(","):
			l.accept(",")
			l.emit(tokenComma)
			return lexStart
		case l.next("("):
			l.accept("(")
			l.emit(tokenLeftParen)
			return lexStart
		case l.next(")"):
			l.accept(")")
			l.emit(tokenRightParen)
			return lexStart
		case l.next("!"):
			return lexExclamation
		case l.next("="):
			return lexEqual
		case l.next("<"):
			return lexLessThan
		case l.next(">"):
			return lexGreaterThan
		case l.next("&"):
			return lexAmpersand
		case l.next("|"):
			return lexPipe
		case l.next("*"):
			l.accept("*")
			l.emit(tokenAsterisk)
			return lexStart
		case l.next("/"):
			l.accept("/")
			l.emit(tokenDivide)
			return lexStart
		case l.next("+"):
			l.accept("+")
			l.emit(tokenPlus)
			return lexStart
		case l.next("-"):
			l.accept("-")
			l.emit(tokenMinus)
			return lexStart
		case l.next("`"):
			return lexString
		case l.next("0123456789"):
			return lexNumber
		case unicode.IsLetter(l.peek()):
			return lexSymbol
		}

		if r := l.get(); r == eof {
			l.emit(tokenEOF)
			return nil
		} else {
			return l.errorf("unexpected rune '%c'", r)
		}
	}
}

func (l *lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) text() string {
	return l.input[l.start:l.index]
}

func (l *lexer) emit(k tokenKind) {
	tok := token{
		kind: k,
		text: l.text(),
	}
	l.tokens <- tok
	l.start = l.index
}

// lex splits the input string into tokens and outputs them to the channel it returns.
func lex(name string, expr string) chan token {
	var l = lexer{
		name:   name,
		input:  expr,
		start:  0,
		index:  0,
		tokens: make(chan token),
	}
	go l.run()
	return l.tokens
}
