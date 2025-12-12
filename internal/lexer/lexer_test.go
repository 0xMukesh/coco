package lexer

import (
	"testing"

	"github.com/0xmukesh/coco/internal/tokens"
)

func TestNextToken(t *testing.T) {
	input := `let x = 2;
let y = 3;
let z = x + y;
let u = 7;

fn sum(a, b) {
  a + b;
}

let isEquals = sum(x, y) == z;
let notEquals = sum(x, y) != u;`

	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "x"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "2"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "y"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "3"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "z"},
		{tokens.ASSIGN, "="},
		{tokens.IDENTIFIER, "x"},
		{tokens.PLUS, "+"},
		{tokens.IDENTIFIER, "y"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "u"},
		{tokens.ASSIGN, "="},
		{tokens.INT, "7"},
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
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "isEquals"},
		{tokens.ASSIGN, "="},
		{tokens.IDENTIFIER, "sum"},
		{tokens.LPAREN, "("},
		{tokens.IDENTIFIER, "x"},
		{tokens.COMMA, ","},
		{tokens.IDENTIFIER, "y"},
		{tokens.RPAREN, ")"},
		{tokens.EQUALS, "=="},
		{tokens.IDENTIFIER, "z"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "notEquals"},
		{tokens.ASSIGN, "="},
		{tokens.IDENTIFIER, "sum"},
		{tokens.LPAREN, "("},
		{tokens.IDENTIFIER, "x"},
		{tokens.COMMA, ","},
		{tokens.IDENTIFIER, "y"},
		{tokens.RPAREN, ")"},
		{tokens.NOT_EQUALS, "!="},
		{tokens.IDENTIFIER, "u"},
		{tokens.SEMICOLON, ";"},
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
