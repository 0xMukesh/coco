package parser

import (
	"regexp"
	"testing"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/lexer"
	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
)

type validateFailureFn func(t *testing.T, input string)

type parserTestItem struct {
	name            string
	input           string
	shouldFail      bool
	expectedAst     *ast.Program
	validateFailure validateFailureFn
}

func newParserTest(name, input string, expectedAst *ast.Program) parserTestItem {
	return parserTestItem{
		name:        name,
		input:       input,
		expectedAst: expectedAst,
		shouldFail:  false,
	}
}

func newParserTestFail(name, input string, validateFailureFn validateFailureFn) parserTestItem {
	return parserTestItem{
		name:            name,
		input:           input,
		shouldFail:      true,
		validateFailure: validateFailureFn,
	}
}

func lexAndCheckTokens(t *testing.T, input string) []tokens.Token {
	l := lexer.New(input)
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

	return tks
}

func expectParseFailure(expectedError string) validateFailureFn {
	return func(t *testing.T, input string) {
		tks := lexAndCheckTokens(t, input)
		p := New(tks)
		p.ParseProgram()

		if !p.HasErrors() {
			t.Fatalf("expected program to have parser errors")
		}

		err := p.Errors()[0]

		parserErrorPrefix := regexp.MustCompile(`^\[line\s+\d+,\s+column\s+\d+:\d+\]\s*`)
		err = parserErrorPrefix.ReplaceAllString(err, "")

		if err != expectedError {
			t.Fatalf("parser error mismatched. expected - %s, got - %s", expectedError, err)
		}
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

func (b astBuilder) addIdentifierExpression(literal string) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: ast.NewIdentifierExpr(literal),
	})

	return b
}

func (b astBuilder) addUnaryExpression(operator tokens.Token, right ast.Expression) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: ast.NewUnaryExpr(operator, right),
	})

	return b
}

func (b astBuilder) addBinaryExpression(operator tokens.Token, left, right ast.Expression) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: ast.NewBinaryExpr(operator, left, right),
	})

	return b
}

func (b astBuilder) addGroupedExpression(expr ast.Expression) astBuilder {
	b.program.Statements = append(b.program.Statements, &ast.ExpressionStatement{
		Expr: ast.NewGroupedExpr(expr),
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
			t.Errorf("statement #%d: float value mismatch: expected %f, got %f", idx, exp.Value, act.Value)
		}
	case *ast.BooleanExpression:
		act := assertType[*ast.BooleanExpression](t, idx, actual)
		if exp.Value != act.Value {
			t.Errorf("statement #%d: boolean value mismatch: expected %t, got %t", idx, exp.Value, act.Value)
		}
	case *ast.StringExpression:
		act := assertType[*ast.StringExpression](t, idx, actual)
		if utils.NormalizeQuotedString(exp.Value) != act.Value {
			t.Errorf("statement #%d: string value mismatch: expected %s, got %s", idx, exp.Value, act.Value)
		}
	case *ast.IdentifierExpression:
		act := assertType[*ast.IdentifierExpression](t, idx, actual)
		if exp.Literal != act.Literal {
			t.Errorf("statement #%d: identifier literal mismatch: expected %s, got %s", idx, exp.Literal, act.Literal)
		}
	case *ast.UnaryExpression:
		act := assertType[*ast.UnaryExpression](t, idx, actual)

		if exp.Token.Literal != act.Token.Literal {
			t.Errorf("statament #%d: unary expression operator mismatch: expected %s, got %s", idx, exp.Token.Literal, act.Token.Literal)
		}

		compareExpression(t, idx, exp.Expr, act.Expr)
	case *ast.BinaryExpression:
		act := assertType[*ast.BinaryExpression](t, idx, actual)

		if exp.Operator.Literal != act.Operator.Literal {
			t.Errorf("statament #%d: binary expression operator mismatch: expected %s, got %s", idx, exp.Operator.Literal, act.Operator.Literal)
		}

		compareExpression(t, idx, exp.Left, act.Left)
		compareExpression(t, idx, exp.Right, act.Right)
	case *ast.GroupedExpression:
		act := assertType[*ast.GroupedExpression](t, idx, actual)

		compareExpression(t, idx, exp.Expr, act.Expr)
	default:
		t.Fatalf("unknown expression type %T", expected)
	}
}

func runParserTest(t *testing.T, tt parserTestItem) {
	tks := lexAndCheckTokens(t, tt.input)
	p := New(tks)
	program := p.ParseProgram()

	if tt.shouldFail {
		if tt.validateFailure == nil {
			t.Fatal("shouldFail is true but validateFailure function is not provided")
		}

		tt.validateFailure(t, tt.input)
	} else {
		if p.HasErrors() {
			for _, e := range p.Errors() {
				t.Log(e)
			}

			t.Fatalf("failed to parse program without errors. got %d errors", len(p.Errors()))
		} else {
			compareAst(t, tt.expectedAst, program)
		}
	}
}
