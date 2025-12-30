package lexer

import (
	"testing"

	"github.com/0xmukesh/coco/internal/tokens"
)

func TestTokens(t *testing.T) {
	input := `let animal = "cat";
let six = 6;
let seven = 7;`
	lexer := New(input)
	tests := []struct {
		expectedType    tokens.TokenType
		expectedLiteral string
	}{
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "animal"},
		{tokens.ASSIGN, "="},
		{tokens.QUOTES, "\""},
		{tokens.IDENTIFIER, "cat"},
		{tokens.QUOTES, "\""},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "six"},
		{tokens.ASSIGN, "="},
		{tokens.INTEGER, "6"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "seven"},
		{tokens.ASSIGN, "="},
		{tokens.INTEGER, "7"},
		{tokens.SEMICOLON, ";"},
	}

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("[test %d] wrong token type. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("[test %d] wrong token literal. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
