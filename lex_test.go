package refine

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		predicate string
		want      []token
	}{
		// Basic cases
		{"0", []token{{tokenInteger, "0"}}},
		{"1_000_000", []token{{tokenInteger, "1_000_000"}}},
		{"`string`", []token{{tokenString, "`string`"}}},
		{"symbol", []token{{tokenSymbol, "symbol"}}},
		{". , ()", []token{{tokenPeriod, "."}, {tokenComma, ","}, {tokenLeftParen, "("}, {tokenRightParen, ")"}}},
		{"== != <= >= < >", []token{{tokenEqual, "=="}, {tokenNotEqual, "!="}, {tokenLessThanOrEqual, "<="}, {tokenGreaterThanOrEqual, ">="}, {tokenLessThan, "<"}, {tokenGreaterThan, ">"}}},
		{"! | & || &&", []token{{tokenLogicalNot, "!"}, {tokenBitwiseOr, "|"}, {tokenBitwiseAnd, "&"}, {tokenLogicalOr, "||"}, {tokenLogicalAnd, "&&"}}},
		{"* / + - << >>", []token{{tokenAsterisk, "*"}, {tokenDivide, "/"}, {tokenPlus, "+"}, {tokenMinus, "-"}, {tokenLeftShift, "<<"}, {tokenRightShift, ">>"}}},

		// Complex cases
		{"-1", []token{{tokenMinus, "-"}, {tokenInteger, "1"}}},
		{"foo > bar", []token{{tokenSymbol, "foo"}, {tokenGreaterThan, ">"}, {tokenSymbol, "bar"}}},
		{"0`hello`symbol", []token{{tokenInteger, "0"}, {tokenString, "`hello`"}, {tokenSymbol, "symbol"}}},
		{"-+-+400_000", []token{{tokenMinus, "-"}, {tokenPlus, "+"}, {tokenMinus, "-"}, {tokenPlus, "+"}, {tokenInteger, "400_000"}}},
		{"5 * (-1>>2)", []token{{tokenInteger, "5"}, {tokenAsterisk, "*"}, {tokenLeftParen, "("}, {tokenMinus, "-"}, {tokenInteger, "1"}, {tokenRightShift, ">>"}, {tokenInteger, "2"}}},

		// Negative cases
		{`"string"`, []token{{tokenError, `unexpected rune '"'`}}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.predicate, func(t *testing.T) {
			tokens := lex(tc.predicate, tc.predicate)
			for _, want := range tc.want {
				got := <-tokens
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("failed to lex %q: got token %v, want token %v", tc.predicate, got, want)
				}
			}
		})
	}
}

func FuzzLexerIntegers(f *testing.F) {
	testCases := []int{
		-1,
		0,
		1,
	}

	for _, tc := range testCases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, i int) {
		predicate := fmt.Sprintf("%d", i)
		tokens := lex(predicate, predicate)
		if i >= 0 {
			// Verify that a number is produced
			tok := <-tokens
			if tok.kind != tokenInteger {
				t.Fatalf("got token %+#v, wanted an integer token", tok)
			}
			if got, want := tok.text, fmt.Sprintf("%d", i); got != want {
				t.Fatalf("got token text %s, want %s", got, want)
			}
		} else {
			// Verify that a minus sign is produced for negative numbers
			tok := <-tokens
			if tok.kind != tokenMinus {
				t.Fatalf("got token %+#v, wanted an minus token", tok)
			}
			if got, want := tok.text, "-"; got != want {
				t.Fatalf("got token text %s, want %s", got, want)
			}
			tok = <-tokens
			if tok.kind != tokenInteger {
				t.Fatalf("got token %+#v, wanted an integer token", tok)
			}
			if got, want := tok.text, fmt.Sprintf("%d", -i); got != want {
				t.Fatalf("got token text %s, want %s", got, want)
			}
		}
		tok := <-tokens
		if tok.kind != tokenEOF {
			t.Fatalf("got token %+#v, wanted EOF token", tok)
		}
	})
}
