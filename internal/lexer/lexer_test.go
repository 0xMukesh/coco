package lexer

import (
	"testing"

	"github.com/0xmukesh/coco/internal/tokens"
)

func TestLexer_SingleCharacterTokens(t *testing.T) {
	tests := []lexerTestItem{
		newLexerTest("plus", "+", tokens.PLUS),
		newLexerTest("minus", "-", tokens.MINUS),
		newLexerTest("star", "*", tokens.STAR),
		newLexerTest("slash", "/", tokens.SLASH),
		newLexerTest("modulo", "%", tokens.MODULO),
		newLexerTest("assign", "=", tokens.ASSIGN),
		newLexerTest("less than", "<", tokens.LESS_THAN),
		newLexerTest("greater than", ">", tokens.GREATER_THAN),
		newLexerTest("bang", "!", tokens.BANG),
		newLexerTest("left paren", "(", tokens.LPAREN),
		newLexerTest("right paren", ")", tokens.RPAREN),
		newLexerTest("left brace", "{", tokens.LBRACE),
		newLexerTest("right brace", "}", tokens.RBRACE),
		newLexerTest("left square", "[", tokens.LSQUARE),
		newLexerTest("right square", "]", tokens.RSQUARE),
		newLexerTest("comman", ",", tokens.COMMA),
		newLexerTest("semicolon", ";", tokens.SEMICOLON),
		newLexerTest("colon", ":", tokens.COLON),
		newLexerTest("illegal", "#", tokens.ILLEGAL),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLexerTest(t, tt)
		})
	}
}

func TestLexer_MultiCharacterTokens(t *testing.T) {
	tests := []lexerTestItem{
		newLexerTest("equals", "==", tokens.EQUALS),
		newLexerTest("not equals", "!=", tokens.NOT_EQUALS),
		newLexerTest("less than equals", "<=", tokens.LESS_THAN_EQUALS),
		newLexerTest("greater than equals", ">=", tokens.GREATER_THAN_EQUALS),
		newLexerTest("and", "&&", tokens.AND),
		newLexerTest("or", "||", tokens.OR),
		newLexerTest("increment", "++", tokens.INCREMENT),
		newLexerTest("decrement", "--", tokens.DECREMENT),
		newLexerTest("double star", "**", tokens.DOUBLE_STAR),
		newLexerTest("plus equal", "+=", tokens.PLUS_EQUAL),
		newLexerTest("minus equal", "-=", tokens.MINUS_EQUAL),
		newLexerTest("star equal", "*=", tokens.STAR_EQUAL),
		newLexerTest("slash equal", "/=", tokens.SLASH_EQUAL),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLexerTest(t, tt)
		})
	}
}

func TestLexer_KeywordsAndIdentifiers(t *testing.T) {
	tests := []lexerTestItem{}

	for k, v := range tokens.PROGRAM_KEYWORDS {
		tests = append(tests, newLexerTest(k, k, v))
	}

	tests = append(tests, newLexerTest("letter", "letter", tokens.IDENTIFIER), newLexerTest("let_x", "let_x", tokens.IDENTIFIER), newLexerTest("let2", "let2", tokens.IDENTIFIER), newLexerTestFail("2let", "2let", expectWrongTokenLiteral()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLexerTest(t, tt)
		})
	}
}

func TestLexer_IntegerLiterals(t *testing.T) {
	tests := []lexerTestItem{
		newLexerTest("zero", "0", tokens.INTEGER),
		newLexerTest("simple", "123", tokens.INTEGER),
		newLexerTest("large", "999999", tokens.INTEGER),
		newLexerTest("leading zero", "007", tokens.INTEGER),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLexerTest(t, tt)
		})
	}
}

func TestLexer_FloatLiterals(t *testing.T) {
	tests := []lexerTestItem{
		newLexerTest("simple", "123.34", tokens.FLOAT),
		newLexerTest("no trailing digit", "123.", tokens.FLOAT),
		newLexerTest("no leading digit", ".12", tokens.FLOAT),
		newLexerTestFail("malformed with leading digit", "123.34.45", expectMalformedFloatLiteral()),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLexerTest(t, tt)
		})
	}
}

func TestLexer_StringLiterals(t *testing.T) {
	tests := []lexerTestItem{
		newLexerTest("simple", `"hello world"`, tokens.STRING),
		newLexerTest("empty", `""`, tokens.STRING),
		newLexerTest("newline escape", `"hello\nworld"`, tokens.STRING),
		newLexerTest("newline", `"\n"`, tokens.STRING),
		newLexerTest("tab escape", `"hello\tworld"`, tokens.STRING),
		newLexerTest("tab", `"\t"`, tokens.STRING),
		newLexerTest("double quotes", `"hello \"world\""`, tokens.STRING),
		newLexerTest("with single quotes", `"hello 'world'"`, tokens.STRING),
		newLexerTestFail("unterminated", `"hello`, expectIllegalToken("unterminated string")),
		newLexerTestFail("invalid escape character", `"hello\z"`, expectIllegalToken("invalid escape character")),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLexerTest(t, tt)
		})
	}
}

func TestLexer_Whitespace(t *testing.T) {
	tokenTypes := []tokens.TokenType{tokens.INTEGER, tokens.PLUS, tokens.INTEGER}
	tokenLiterals := []string{"5", "+", "2"}

	tests := []lexerTestItem{
		newLexerTestVerbose("no spaces", "5+2", tokenTypes, tokenLiterals),
		newLexerTestVerbose("simple spaces", "5 + 2", tokenTypes, tokenLiterals),
		newLexerTestVerbose("many spaces", "5    +    2", tokenTypes, tokenLiterals),
		newLexerTestVerbose("tabs", "5\t\t+\t\t2", tokenTypes, tokenLiterals),
		newLexerTestVerbose("newline", "5\n\n+\n\n2", tokenTypes, tokenLiterals),
		newLexerTestVerbose("mixed", "5   \n\n  \t\t+\t\t\n  2", tokenTypes, tokenLiterals),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLexerTest(t, tt)
		})
	}
}

func TestLexer_Comments(t *testing.T) {
	tokensTypes := []tokens.TokenType{tokens.INTEGER}
	tokenLiterals := []string{"5"}

	tests := []lexerTestItem{
		newLexerTestVerbose("simple", `// hello world
5`, tokensTypes, tokenLiterals),
		newLexerTestVerbose("multiple comments", `// hello world
5 // test`, tokensTypes, tokenLiterals),
		newLexerTestVerbose("multiline", `/*
multi
line
comment */ 5`, tokensTypes, tokenLiterals),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLexerTest(t, tt)
		})
	}
}
