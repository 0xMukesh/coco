package codegen

import (
	"fmt"
	"os"
	"testing"

	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/parser"
	"github.com/0xmukesh/coco/internal/typechecker"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

func TestCodegen(t *testing.T) {
	source := `exit 3 + 2;`
	l := lexer.New(source)
	tks := l.Lex()

	p := parser.New(tks)
	program := p.ParseProgram()

	if p.HasErrors() {
		for _, e := range p.Errors() {
			fmt.Println(e)
		}

		t.FailNow()
	}

	fmt.Println(program.String())

	tenv := cotypes.NewTypeEnvironment()
	tc := typechecker.New(tenv)
	tc.Transform(program)

	if tc.HasErrors() {
		for _, e := range tc.Errors() {
			fmt.Println(e)
		}

		t.FailNow()
	}

	cg := New()
	cg.Generate(program)
	ir := cg.EmitIR()

	outFile := "../../build/test.ll"

	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		if _, err := os.Create(outFile); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.WriteFile(outFile, []byte(ir), 0644); err != nil {
		t.Fatal(err)
	}
}
