package lexer

import (
	"strconv"
	"strings"
	"testing"

	"github.com/0xmukesh/coco/internal/tokens"
)

type validateFailureFn func(t *testing.T, input string)

type lexerTestItem struct {
	name                  string
	input                 string
	expectedTokenType     tokens.TokenType
	expectedTokenTypes    []tokens.TokenType
	expectedTokenLiterals []string
	shouldFail            bool
	validateFailure       validateFailureFn
}

func newLexerTest(name, input string, expectedTokenType tokens.TokenType) lexerTestItem {
	return lexerTestItem{
		name:              name,
		input:             input,
		expectedTokenType: expectedTokenType,
		shouldFail:        false,
	}
}

func newLexerTestVerbose(name, input string, expectedTokenTypes []tokens.TokenType, expectedTokenLiterals []string) lexerTestItem {
	return lexerTestItem{
		name:                  name,
		input:                 input,
		expectedTokenTypes:    expectedTokenTypes,
		expectedTokenLiterals: expectedTokenLiterals,
	}
}

func newLexerTestFail(name, input string, validateFailure validateFailureFn) lexerTestItem {
	return lexerTestItem{
		name:            name,
		input:           input,
		shouldFail:      true,
		validateFailure: validateFailure,
	}
}

func runLexerTest(t *testing.T, tt lexerTestItem) {
	if tt.shouldFail {
		if tt.validateFailure == nil {
			t.Fatal("shouldFail is true but validateFailure function is not provided")
		}

		tt.validateFailure(t, tt.input)
	} else {
		l := New(tt.input)
		tks := l.Lex()

		tokenTypesToCheck := []tokens.TokenType{}

		if tt.expectedTokenTypes != nil {
			if len(tks) != len(tt.expectedTokenTypes) {
				t.Fatalf("tt.expectedTokenTypes (%d) and tokens returned by lexer (%d) length mismatch", len(tt.expectedTokenTypes), len(tks))
			}

			if len(tks) != len(tt.expectedTokenLiterals) {
				t.Fatalf("tt.expectedTokenLiterals (%d) and tokens returned by lexer (%d) length mismatch", len(tt.expectedTokenLiterals), len(tks))
			}

			tokenTypesToCheck = append(tokenTypesToCheck, tt.expectedTokenTypes...)
		} else {
			if len(tks) != 1 {
				t.Fatalf("expected one token, got %d tokens", len(tks))
			}

			tokenTypesToCheck = append(tokenTypesToCheck, tt.expectedTokenType)
		}

		for i, tk := range tks {
			if tk.Type != tokenTypesToCheck[i] {
				t.Fatalf("expected token type - %q, got - %q", tt.expectedTokenType, tk.Type)
			}

			input := ""

			if tt.expectedTokenLiterals != nil {
				input = tt.expectedTokenLiterals[i]
			} else {
				input = tt.input
			}

			// if the token is string then first normalize the input and the wrap it with double quotes at the end
			if tk.Type == tokens.STRING {
				normalizedInput, err := strconv.Unquote(input)
				if err == nil {
					input = "\"" + normalizedInput + "\""
				}
			}

			if tk.Literal != input {
				t.Fatalf("expected token literal - %s, got - %s", input, tk.Literal)
			}
		}
	}
}

func expectWrongTokenLiteral() validateFailureFn {
	return func(t *testing.T, input string) {
		l := New(input)
		tks := l.Lex()
		gotToken := tks[0]

		if gotToken.Literal == input {
			t.Fatalf("expected wrong token literal, but got expected token literal: %q", gotToken.Literal)
		}

		t.Logf("literal correctly wrong: expected=%q, got=%q", input, gotToken.Literal)
	}
}

func expectMalformedFloatLiteral() validateFailureFn {
	return func(t *testing.T, input string) {
		splits := strings.Split(input, ".")
		if len(splits) < 2 {
			t.Fatalf("expected length of splits to be >= 2, got %d", len(splits))
		}

		l := New(input)
		tks := l.Lex()

		if len(tks) < 2 {
			t.Fatalf("expected atleast 2 tokens, got %d tokens", len(tks))
		}

		if tks[0].Type != tokens.FLOAT {
			t.Fatalf("expected first token to be a float, got %q", tks[0].Type)
		}

		floatLiteral := splits[0] + "." + splits[1]

		if tks[0].Literal != floatLiteral {
			t.Fatalf("expected float literal to be %s, got %s", floatLiteral, tks[0].Literal)
		}

		if tks[1].Type != tokens.ILLEGAL {
			t.Fatalf("expected next token after float literal to be an illegal token, got %q", tks[1].Type)
		}

		if tks[1].Literal != "." {
			t.Fatalf("expected next token literal after float literal to be \".\", got %s", tks[1].Literal)
		}
	}
}

func expectIllegalToken(expectedErrMsg string) validateFailureFn {
	return func(t *testing.T, input string) {
		l := New(input)
		tks := l.Lex()

		foundMatch := false

		for _, tk := range tks {
			foundMatch = tk.Type == tokens.ILLEGAL && tk.Literal == expectedErrMsg
			if foundMatch {
				break
			}
		}

		if !foundMatch {
			t.Fatalf("couldn't find an illegal token with %s as err msg among the following tokens: \n%+v", expectedErrMsg, tks)
		}
	}
}
