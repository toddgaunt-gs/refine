package refine

import (
	"fmt"
	"testing"
)

//const text = "2 < 30 && 1 > 5"
//const text = "2 + 5 + -1 * -2"
const text = "(2 < 5) (==) (2 < 1)"

func TestLexer(t *testing.T) {
	tokens := lex("test", text)
	for {
		select {
		case tok := <-tokens:
			if tok.kind == tokenEOF {
				return
			}
			fmt.Printf("Token: %d '%s'\n", tok.kind, tok.text)
		}
	}
}

func TestParser(t *testing.T) {
	tokens := lex("test", text)
	expr, err := parse(tokens)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	box, err := expr.Eval()
	if err != nil {
		t.Fatalf("failed to eval: %v", err)
	}
	fmt.Printf("(%s) evaluates to: %v\n", text, box.val)
}
