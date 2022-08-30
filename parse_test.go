package refine

import (
	"errors"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		name    string
		tokens  []token
		want    expression
		wantErr error
	}{
		{
			name:   "basic math",
			tokens: []token{{tokenInteger, "1"}, {tokenPlus, "+"}, {tokenInteger, "2"}},
			want: &binaryExpression{
				op:    binaryPlus,
				left:  &integerExpression{text: "1"},
				right: &integerExpression{text: "2"},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tokens := make(chan token)
			go func() {
				for _, tok := range tc.tokens {
					tokens <- tok
				}
				close(tokens)
			}()
			expr, err := parse(tokens)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got err %v, want err %v", err, tc.wantErr)
			}
			if !reflect.DeepEqual(expr, tc.want) {
				t.Fatalf("got expr %+#v, want expr %+#v", expr, tc.want)
			}
		})
	}
}
