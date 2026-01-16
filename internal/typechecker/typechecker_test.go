package typechecker

import (
	"fmt"
	"testing"

	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/parser"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

func TestTypeChecker(t *testing.T) {
	source := `let a = 1;
let b = 2;
exit a + b;`
	l := lexer.New(source)
	tks := l.Lex()
	p := parser.New(tks)
	program := p.ParseProgram()
	tenv := cotypes.NewTypeEnvironment()
	tc := New(tenv)

	tc.Transform(program)

	if tc.HasErrors() {
		for _, e := range tc.Errors() {
			fmt.Println(e)
		}
	}
}
