package lexer

import (
	"testing"

	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/utils"
)

func TestLexer(t *testing.T) {
	source := `+-*/
(){}
  !<>
=`
	tests := []struct {
		wantTokenType tokens.TokenType
		wantLiteral   string
		wantLine      int
		wantColumn    int
	}{
		{tokens.PLUS, "+", 1, 0},
		{tokens.MINUS, "-", 1, 1},
		{tokens.STAR, "*", 1, 2},
		{tokens.SLASH, "/", 1, 3},
		{tokens.LPAREN, "(", 2, 0},
		{tokens.RPAREN, ")", 2, 1},
		{tokens.LBRACE, "{", 2, 2},
		{tokens.RBRACE, "}", 2, 3},
		{tokens.BANG, "!", 3, 2},
		{tokens.LESS_THAN, "<", 3, 3},
		{tokens.GREATER_THAN, ">", 3, 4},
		{tokens.ASSIGN, "=", 4, 0},
		{tokens.EOF, "", 4, 0},
	}
	l := New(source)

	for i, tt := range tests {
		tok := l.NextToken()

		if tt.wantTokenType != tok.Type {
			t.Fatal(utils.TestMismatchErrorBuilder(i, "token type", tt.wantTokenType, tok.Type))
		}

		if tt.wantLiteral != tok.Literal {
			t.Fatal(utils.TestMismatchErrorBuilder(i, "token literal", tt.wantLiteral, tok.Literal))
		}

		if tt.wantLine != tok.Line {
			t.Fatal(utils.TestMismatchErrorBuilder(i, "token line", tt.wantLine, tok.Line))
		}

		if tt.wantColumn != tok.Column {
			t.Fatal(utils.TestMismatchErrorBuilder(i, "token column", tt.wantColumn, tok.Column))
		}
	}
}
