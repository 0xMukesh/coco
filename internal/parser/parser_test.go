package parser

import (
	"fmt"
	"testing"

	"github.com/0xmukesh/coco/internal/lexer"
)

func TestParser(t *testing.T) {
	input := `let x = 1 + 1;
return x + x;`
	l := lexer.New(input)
	tks := l.Tokenize()

	p := New(tks)
	program := p.Parse()

	for _, v := range program.Statements {
		fmt.Printf("%+v\n", v)
	}
}
