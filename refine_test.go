package refine

import (
	"errors"
	"reflect"
	"testing"
)

//const text = "2 < 30 && 1 > 5"
//const text = "(2 < 5) == (2 < 1)"
//const text = "2 + 5 + -1 * -2"
const text = "foo > bar"

func TestLexer(t *testing.T) {
	testCases := []struct {
		predicate string
		want      []token
	}{
		// Positive cases
		{"0", []token{{tokenInteger, "0"}}},
		{"-1", []token{{tokenMinus, "-"}, {tokenInteger, "1"}}},
		{"`string`", []token{{tokenString, "`string`"}}},
		{"symbol", []token{{tokenSymbol, "symbol"}}},
		{"foo > bar", []token{{tokenSymbol, "foo"}, {tokenGreaterThan, ">"}, {tokenSymbol, "bar"}}},
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

func TestParser(t *testing.T) {
	testCases := []struct {
		name string
		expr string
		want error
	}{
		{"greater than", "foo > bar", nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.expr, func(t *testing.T) {
			tokens := lex(tc.name, tc.expr)
			_, err := parse(tokens)
			if !errors.Is(err, tc.want) {
				t.Fatalf("failed to parse %q: %v", tc.expr, err)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	type checkDependent struct {
		A int `refine:"A > B"`
		B int `refine:"B >= 0"`
	}

	type checkNil struct {
		A *int           `refine:"A != nil"`
		B *int           `refine:"B == nil"`
		C *int           `refine:"C == B"`
		D []string       `refine:"D != nil"`
		E map[int]string `refine:"E != nil"`
	}

	type checkString struct {
		S string "refine:\"S == `foo`\""
	}

	type checkStruct struct {
		C checkNil `refine:"C.A != nil"`
	}

	testCases := []struct {
		name  string
		value any

		want error
	}{
		{
			name:  "NotStructErr",
			value: 2,

			want: ErrNotStruct,
		},
		{
			name: "DependentFields",
			value: checkDependent{
				A: 2,
				B: 1,
			},

			want: nil,
		},
		{
			name: "DependentFieldsErr",
			value: checkDependent{
				A: 2,
				B: 2,
			},

			want: ErrNotMet,
		},
		{
			name: "NilFields",
			value: checkNil{
				A: func(x int) *int { return &x }(2),
				B: nil,
				C: nil,
				D: []string{},
				E: map[int]string{},
			},

			want: nil,
		},
		{
			name: "StringMet",
			value: checkString{
				S: "foo",
			},

			want: nil,
		},
		{
			name: "StringNotMet",
			value: checkString{
				S: "not foo",
			},

			want: ErrNotMet,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := Check(tc.value)
			if !errors.Is(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestCheckPredicate(t *testing.T) {
	testCases := []struct {
		name      string
		value     any
		predicate string

		want error
	}{
		{
			name:      "not negative",
			value:     2,
			predicate: "? >= 0",

			want: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := CheckPredicate(tc.value, tc.predicate)
			if !errors.Is(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}
