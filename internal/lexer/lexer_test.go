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
=
== >= <=
!=

// this is a single line comment

/*
this
is
a
multi line
comment
*/

==

// test test

/*
test test
*/

abc
let +=
const xyz`
	tests := []struct {
		wantTokenType   tokens.TokenType
		wantLiteral     string
		wantLine        int
		wantStartColumn int
		wantEndColumn   int
	}{
		{tokens.PLUS, "+", 1, 0, 1},
		{tokens.MINUS, "-", 1, 1, 2},
		{tokens.STAR, "*", 1, 2, 3},
		{tokens.SLASH, "/", 1, 3, 4},
		{tokens.LPAREN, "(", 2, 0, 1},
		{tokens.RPAREN, ")", 2, 1, 2},
		{tokens.LBRACE, "{", 2, 2, 3},
		{tokens.RBRACE, "}", 2, 3, 4},
		{tokens.BANG, "!", 3, 2, 3},
		{tokens.LESS_THAN, "<", 3, 3, 4},
		{tokens.GREATER_THAN, ">", 3, 4, 5},
		{tokens.ASSIGN, "=", 4, 0, 1},
		{tokens.EQUALS, "==", 5, 0, 2},
		{tokens.GREATER_THAN_EQUALS, ">=", 5, 3, 5},
		{tokens.LESS_THAN_EQUALS, "<=", 5, 6, 8},
		{tokens.NOT_EQUALS, "!=", 6, 0, 2},
		{tokens.EQUALS, "==", 18, 0, 2},
		{tokens.IDENTIFIER, "abc", 26, 0, 3},
		{tokens.LET, "let", 27, 0, 3},
		{tokens.PLUS, "+", 27, 4, 5},
		{tokens.ASSIGN, "=", 27, 5, 6},
		{tokens.CONST, "const", 28, 0, 5},
		{tokens.IDENTIFIER, "xyz", 28, 6, 9},
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

		if tt.wantStartColumn != tok.StartColumn {
			t.Fatal(utils.TestMismatchErrorBuilder(i, "token start column", tt.wantStartColumn, tok.StartColumn))
		}

		if tt.wantEndColumn != tok.EndColumn {
			t.Fatal(utils.TestMismatchErrorBuilder(i, "token end column", tt.wantEndColumn, tok.EndColumn))
		}
	}
}
