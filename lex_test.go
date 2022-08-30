package refine

import (
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
