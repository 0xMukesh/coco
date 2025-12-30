package parser

import (
	"testing"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `let x = 1;
let y = 2;
let = 1;`
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"}, {"y"},
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	errors := p.Errors()

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}

		t.FailNow()
	}

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 ast.Statements. got %d", len(program.Statements))
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if stmt.TokenLiteral() != "let" {
			t.Fatalf("expected let statement. got %s", stmt.TokenLiteral())
		}

		letStmt, ok := stmt.(*ast.LetStatement)
		if !ok {
			t.Fatalf("failed to parse #%d statement to *ast.LetStatement", i)
		}

		if letStmt.Identifier.Value != tt.expectedIdentifier {
			t.Fatalf("expected identifier with %s value. got %s", tt.expectedIdentifier, letStmt.Identifier.Value)
		}

		if letStmt.Identifier.TokenLiteral() != tt.expectedIdentifier {
			t.Fatalf("expected identifier with %s value. got %s", tt.expectedIdentifier, letStmt.Identifier.Value)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `return 5;
return 10 + 11 + 43;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 ast.Statements. got %d", len(program.Statements))
	}

	for i, stmt := range program.Statements {
		returntStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("failed to parse #%d statement to *ast.ReturnStatement", i)
		}

		if returntStmt.TokenLiteral() != "return" {
			t.Fatalf("expected return statement. got %s", returntStmt.TokenLiteral())
		}
	}
}
