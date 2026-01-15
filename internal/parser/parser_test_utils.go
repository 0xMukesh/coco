package parser

import (
	"testing"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
)

type parserTestItem struct {
	name        string
	input       string
	expectedAst *ast.Program
}

func newParserTestItem(name, input string, expectedAst *ast.Program) parserTestItem {
	return parserTestItem{
		name:        name,
		input:       input,
		expectedAst: expectedAst,
	}
}

func assertType[T any](t *testing.T, idx int, actual any) T {
	t.Helper()
	result, ok := actual.(T)
	if !ok {
		var zero T
		t.Fatalf("statement #%d: expected %T, got %T", idx, zero, actual)
	}

	return result
}

type astBuilder struct {
	program    *ast.Program
	currentCol int
}

func newAstBuilder() astBuilder {
	return astBuilder{
		program: &ast.Program{},
	}
}

func (b astBuilder) addIntegerLiteralExpression(value int64) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: ast.NewIntegerExpr(value),
	})

	return b
}

func (b astBuilder) addFloatLiteralExpression(value float64) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: ast.NewFloatExpr(value),
	})

	return b
}

func (b astBuilder) addBooleanLiteralExpression(value bool) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: ast.NewBooleanExpr(value),
	})

	return b

}

func (b astBuilder) addStringLiteralExpression(value string) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: ast.NewStringExpr(value),
	})

	return b
}

func (b astBuilder) addUnaryExpression(right ast.Expression) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: &ast.UnaryExpression{
			Expr: right,
		},
	})

	return b
}

func (b astBuilder) toProgram() *ast.Program {
	return b.program
}

func compareAst(t *testing.T, expected, actual *ast.Program) {
	t.Helper()

	if len(expected.Statements) != len(actual.Statements) {
		t.Fatalf("statement count mismatch. expected %d, got %d", len(expected.Statements), len(actual.Statements))
	}

	for i := range actual.Statements {
		compareStatement(t, i, expected.Statements[i], actual.Statements[i])
	}
}

func compareStatement(t *testing.T, idx int, expected, actual ast.Statement) {
	t.Helper()

	switch exp := expected.(type) {
	case *ast.ExpressionStatement:
		act := assertType[*ast.ExpressionStatement](t, idx, actual)
		compareExpression(t, idx, exp.Expr, act.Expr)
	default:
		t.Fatalf("unknown statement type %T", expected)
	}
}

func compareExpression(t *testing.T, idx int, expected, actual ast.Expression) {
	t.Helper()

	switch exp := expected.(type) {
	case *ast.IntegerExpression:
		act := assertType[*ast.IntegerExpression](t, idx, actual)
		if exp.Value != act.Value {
			t.Errorf("statement #%d: integer value mismatch: expected %d, got %d", idx, exp.Value, act.Value)
		}
	case *ast.FloatExpression:
		act := assertType[*ast.FloatExpression](t, idx, actual)
		if exp.Value != act.Value {
			t.Errorf("statement #%d: integer value mismatch: expected %f, got %f", idx, exp.Value, act.Value)
		}
	case *ast.BooleanExpression:
		act := assertType[*ast.BooleanExpression](t, idx, actual)
		if exp.Value != act.Value {
			t.Errorf("statement #%d: integer value mismatch: expected %t, got %t", idx, exp.Value, act.Value)
		}
	case *ast.StringExpression:
		act := assertType[*ast.StringExpression](t, idx, actual)
		if utils.NormalizeQuotedString(exp.Value) != act.Value {
			t.Errorf("statement #%d: integer value mismatch: expected %s, got %s", idx, exp.Value, act.Value)
		}
	case *ast.UnaryExpression:
		act := assertType[*ast.UnaryExpression](t, idx, actual)
		compareExpression(t, idx, exp.Expr, act.Expr)
	default:
		t.Fatalf("unknown expression type %T", expected)
	}
}

func runParserTest(t *testing.T, tt parserTestItem) {
	l := lexer.New(tt.input)
	tks := l.Lex()

	toParse := true

	for _, tok := range tks {
		if tok.Type == tokens.ILLEGAL {
			t.Logf("found illegal token while lexing - %+v", tok)
			toParse = false
		}
	}

	if !toParse {
		t.Fatalf("fail to lex program without errors")
	}

	p := New(tks)
	program := p.ParseProgram()

	if p.HasErrors() {
		for _, e := range p.Errors() {
			t.Log(e)
		}

		t.Fatalf("failed to parse program without errors. got %d errors", len(p.Errors()))
	} else {
		compareAst(t, tt.expectedAst, program)
	}
}
