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
				text: "2",
			},
		}
	}

	testCases := []struct {
		name string
		expr expression

		wantErr error
		wantVal box
	}{
		{
			name: "binaryPlus",
			expr: binaryExpr(binaryPlus),

			wantErr: nil,
			wantVal: box{
				kind: boxInt,
				val:  13,
			},
		},
		{
			name: "binaryMinus",
			expr: binaryExpr(binaryMinus),

			wantErr: nil,
			wantVal: box{
				kind: boxInt,
				val:  9,
			},
		},
		{
			name: "binaryMultiply",
			expr: binaryExpr(binaryMultiply),

			wantErr: nil,
			wantVal: box{
				kind: boxInt,
				val:  22,
			},
		},
		{
			name: "binaryDivide",
			expr: binaryExpr(binaryDivide),

			wantErr: nil,
			wantVal: box{
				kind: boxInt,
				val:  5,
			},
		},
		{
			name: "binaryLeftShift",
			expr: binaryExpr(binaryLeftShift),

			wantErr: nil,
			wantVal: box{
				kind: boxInt,
				val:  44,
			},
		},
		{
			name: "binaryRightShift",
			expr: binaryExpr(binaryRightShift),

			wantErr: nil,
			wantVal: box{
				kind: boxInt,
				val:  2,
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
