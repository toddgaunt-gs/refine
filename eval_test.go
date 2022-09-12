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
				text:  "11",
				value: 11,
			},
			right: &integerExpression{
				text:  "2",
				value: 2,
			},
		}
	}

	var errExpr = func(op binaryOperator) *binaryExpression {
		return &binaryExpression{
			op: op,
			left: &integerExpression{
				text:  "0",
				value: 0,
			},
			right: &booleanExpression{
				text:  "true",
				value: true,
			},
		}
	}

	testCases := []struct {
		name string
		expr expression

		wantErr error
		wantVal box
	}{
		// Successful evals
		{
			name:    "binaryPlus",
			expr:    binaryExpr(binaryPlus),
			wantVal: box{kind: boxInt, val: 13},
		},
		{
			name:    "binaryMinus",
			expr:    binaryExpr(binaryMinus),
			wantVal: box{kind: boxInt, val: 9},
		},
		{
			name:    "binaryMultiply",
			expr:    binaryExpr(binaryMultiply),
			wantVal: box{kind: boxInt, val: 22},
		},
		{
			name:    "binaryDivide",
			expr:    binaryExpr(binaryDivide),
			wantVal: box{kind: boxInt, val: 5},
		},
		{
			name:    "binaryLeftShift",
			expr:    binaryExpr(binaryLeftShift),
			wantVal: box{kind: boxInt, val: 44},
		},
		{
			name:    "binaryRightShift",
			expr:    binaryExpr(binaryRightShift),
			wantVal: box{kind: boxInt, val: 2},
		},
		{
			name:    "binaryEqual",
			expr:    binaryExpr(binaryEqual),
			wantVal: box{kind: boxBool, val: false},
		},
		{
			name:    "binaryNotEqual",
			expr:    binaryExpr(binaryNotEqual),
			wantVal: box{kind: boxBool, val: true},
		},
		{
			name:    "binaryLessThan",
			expr:    binaryExpr(binaryLessThan),
			wantVal: box{kind: boxBool, val: false},
		},
		{
			name:    "binaryGreaterThan",
			expr:    binaryExpr(binaryGreaterThan),
			wantVal: box{kind: boxBool, val: true},
		},
		{
			name:    "binaryLessThanOrEqual",
			expr:    binaryExpr(binaryLessThanOrEqual),
			wantVal: box{kind: boxBool, val: false},
		},
		{
			name:    "binaryGreaterThanOrEqual",
			expr:    binaryExpr(binaryGreaterThanOrEqual),
			wantVal: box{kind: boxBool, val: true},
		},
		// Erroneous evals
		{
			name:    "binaryPlusTypeError",
			expr:    errExpr(binaryPlus),
			wantErr: errors.New("type mismatch: 2 != 1"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ev := newEvaluator()
			tc.expr.Accept(ev)
			val, err := ev.Result, ev.Err

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got err %v, want err %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(val, tc.wantVal) {
				t.Fatalf("got value %v, want value %v", val, tc.wantVal)
			}
		})
	}
}
