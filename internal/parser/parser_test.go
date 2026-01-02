package parser

import (
	"reflect"
	"testing"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/tokens"
)

func TestLetStatement(t *testing.T) {
	input := `let x = 1;
let y = 2;`
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

func TestIdentifierExpression(t *testing.T) {
	input := `abc;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 ast.Statements. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("failed to parse statement to *ast.ExpressionStatement")
	}

	identifier, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("failed to parse expression to *ast.Identifier")
	}

	if identifier.TokenLiteral() != "abc" {
		t.Fatalf("expected abc as identifier.TokenLiteral but got %s", identifier.TokenLiteral())
	}

	if identifier.Value != "abc" {
		t.Fatalf("expected abc as identifier.Value but got %s", identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `100;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 ast.Statements. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("failed to parse statement to *ast.ExpressionStatement")
	}

	expr, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("failed to parse expression to *ast.IntegerLiteral")
	}

	if expr.TokenLiteral() != "100" {
		t.Fatalf("expected 100 as identifier.IntegerLiteral but got %s", expr.TokenLiteral())
	}

	if expr.Value != 100 {
		t.Fatalf("expected 100 as identifier.Value but got %s", expr.TokenLiteral())
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	input := `1.54678;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 ast.Statements. got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("failed to parse statement to *ast.ExpressionStatement")
	}

	expr, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("failed to parse expression to *ast.FloatLiteral")
	}

	if expr.TokenLiteral() != "1.54678" {
		t.Fatalf("expected 1.54678 as identifier.IntegerLiteral but got %s", expr.TokenLiteral())
	}

	if expr.Value != 1.54678 {
		t.Fatalf("expected 1.54678 as identifier.Value but got %s", expr.TokenLiteral())
	}
}

func TestUnaryExpression(t *testing.T) {
	tests := []struct {
		input                  string
		expectedTokenType      tokens.TokenType
		expectedExpressionType reflect.Type
		expectedOutput         string
	}{
		{"!5;", tokens.BANG, reflect.TypeFor[*ast.IntegerLiteral](), "(! 5)"},
		{"-xyz;", tokens.MINUS, reflect.TypeFor[*ast.Identifier](), "(- xyz)"},
		{"-2.345;", tokens.MINUS, reflect.TypeFor[*ast.FloatLiteral](), "(- 2.345)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()

		for _, stmt := range program.Statements {
			exprStmt, ok := stmt.(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected *ast.ExpressionStatement. got %T", exprStmt)
			}

			unaryExpr, ok := exprStmt.Expression.(*ast.UnaryExpression)
			if !ok {
				t.Fatalf("expected *ast.UnaryExpression. got %T", unaryExpr)
			}

			if unaryExpr.Token.Type != tt.expectedTokenType {
				t.Fatalf("expected %s token type. got %s", tt.expectedTokenType, unaryExpr.Token.Type)
			}

			if reflect.TypeOf(unaryExpr.Expression) != tt.expectedExpressionType {
				t.Fatalf("expected %v expression type. got %T", tt.expectedExpressionType, unaryExpr.Expression)
			}

			if unaryExpr.String() != tt.expectedOutput {
				t.Fatalf("expected %s expression output. got %s", tt.expectedOutput, unaryExpr.String())
			}
		}
	}
}
