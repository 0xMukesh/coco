package lexer

import (
	"testing"

	"github.com/0xmukesh/coco/internal/tokens"
)

func TestTokens(t *testing.T) {
	input := `let animal = "cat";
let isCat = animal == "cat";
let isNotCat = animal != "cat";`
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
		{tokens.IDENTIFIER, "isCat"},
		{tokens.ASSIGN, "="},
		{tokens.IDENTIFIER, "animal"},
		{tokens.EQUALS, "=="},
		{tokens.QUOTES, "\""},
		{tokens.IDENTIFIER, "cat"},
		{tokens.QUOTES, "\""},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENTIFIER, "isNotCat"},
		{tokens.ASSIGN, "="},
		{tokens.IDENTIFIER, "animal"},
		{tokens.NOT_EQUALS, "!="},
		{tokens.QUOTES, "\""},
		{tokens.IDENTIFIER, "cat"},
		{tokens.QUOTES, "\""},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, ""},
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
