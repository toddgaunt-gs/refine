package refine

import (
	"errors"
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	var binaryExpr = func(op binaryOperator) *binaryExpression {
		return &binaryExpression{
			op: op,
			left: &integerExpression{
				text: "11",
			},
			right: &integerExpression{
				text: "7",
			},
		}
	}

	testCases := []struct {
		name string
		expr expression

		wantErr error
		wantVal Value
	}{
		{
			name: "binaryPlus",
			expr: binaryExpr(binaryPlus),

			wantErr: nil,
			wantVal: Value{
				kind: boxInt,
				val:  18,
			},
		},
		{
			name: "binaryMinus",
			expr: binaryExpr(binaryMinus),

			wantErr: nil,
			wantVal: Value{
				kind: boxInt,
				val:  4,
			},
		},
		{
			name: "binaryMultiply",
			expr: binaryExpr(binaryMultiply),

			wantErr: nil,
			wantVal: Value{
				kind: boxInt,
				val:  77,
			},
		},
		{
			name: "binaryDivide",
			expr: binaryExpr(binaryDivide),

			wantErr: nil,
			wantVal: Value{
				kind: boxInt,
				val:  1,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ev := newEvaluator()
			val, err := eval(ev, tc.expr)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got err %v, want err %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(val, tc.wantVal) {
				t.Fatalf("got value %v, want value %v", val, tc.wantVal)
			}
		})
	}
}
