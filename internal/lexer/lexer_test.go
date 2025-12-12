package lexer

import (
	"testing"

	"github.com/0xmukesh/coco/internal/tokens"
)

func TestNextToken(t *testing.T) {
	input := `let x = 5;
let y = 10;
let z = x + y;

fn sum(a, b) {
  a + b;
}`

	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "x"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "5"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "y"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "10"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "z"},
		{tokens.ASSIGN, "="},
		{tokens.IDENTIFIER, "x"},
		{tokens.PLUS, "+"},
		{tokens.IDENTIFIER, "y"},
		{tokens.SEMICOLON, ";"},
		{tokens.FUNCTION, "fn"},
		{tokens.IDENTIFIER, "sum"},
		{tokens.LPAREN, "("},
		{tokens.IDENTIFIER, "a"},
		{tokens.COMMA, ","},
		{tokens.IDENTIFIER, "b"},
		{tokens.RPAREN, ")"},
		{tokens.LBRACE, "{"},
		{tokens.IDENTIFIER, "a"},
		{tokens.PLUS, "+"},
		{tokens.IDENTIFIER, "b"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%q, got =%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
