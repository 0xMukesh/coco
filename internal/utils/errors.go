package utils

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/tokens"
)

func TestMismatchErrorBuilder(testIdx int, v string, expected, got any) string {
	return fmt.Sprintf("[test #%d] %s mismatch. expected=%v, got=%v", testIdx, v, expected, got)
}

func ParserErrorBuilder(token tokens.Token, message string) string {
	return fmt.Sprintf("[line %d, column %d:%d] %s", token.Line, token.StartColumn, token.EndColumn, message)
}

func ParserExpectedNextTokenToBeErrorBuilder(token tokens.Token, expectedTokenType tokens.TokenType) string {
	msg := fmt.Sprintf("expected type of next token to be %s, got %s instead", expectedTokenType, token.Type)
	return ParserErrorBuilder(token, msg)
}

func ParserExpectedCurrentTokenToBeErrorBuilder(token tokens.Token, expectedTokenType tokens.TokenType) string {
	msg := fmt.Sprintf("expected type of current token to be %s, got %s instead", expectedTokenType, token.Type)
	return ParserErrorBuilder(token, msg)
}

func ParserNoPrefixFnErrorBuilder(token tokens.Token) string {
	msg := fmt.Sprintf("no prefix function found for %s token", token.Type)
	return ParserErrorBuilder(token, msg)
}

func ParserExpressionExpectedErrorBuilder(token tokens.Token) string {
	msg := fmt.Sprintf("expression expected after %s token", token.Type)
	return ParserErrorBuilder(token, msg)
}

func ParserFailedToParseExpressionErrorBuilder(token tokens.Token, err string) string {
	msg := fmt.Sprintf("failed to parse expression: %s", err)
	return ParserErrorBuilder(token, msg)
}
