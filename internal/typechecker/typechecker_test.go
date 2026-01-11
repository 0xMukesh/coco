package typechecker

import (
	"fmt"
	"testing"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/parser"
	"github.com/0xmukesh/coco/internal/types"
)

func TestTypeChecker(t *testing.T) {
	source := `true == true;`
	l := lexer.New(source)
	tks := l.Lex()
	p := parser.New(tks)
	program := p.ParseProgram()
	tenv := types.NewTypeEnvironment()
	tc := New(tenv)

	tc.Transform(program)

	if tc.HasErrors() {
		for _, e := range tc.Errors() {
			fmt.Println(e)
		}
	} else {
		fmt.Printf("%+v\n", (program.Statements[0].(*ast.ExpressionStatement)).Expr.GetType())
	}
}
