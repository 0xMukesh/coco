package typechecker

import (
	"fmt"
	"testing"

	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/parser"
)

func TestTypeChecker(t *testing.T) {
	source := `let a = 1;
{
  let a = 1;
  let a = 3;
}`
	l := lexer.New(source)
	tks := l.Lex()
	p := parser.New(tks)
	program := p.ParseProgram()
	tc := New()

	tc.Transform(program)

	if tc.HasErrors() {
		for _, e := range tc.Errors() {
			fmt.Println(e)
		}
	}
}
