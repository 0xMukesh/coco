package parser

import (
	"testing"

	"github.com/0xmukesh/coco/internal/ast"
)

func TestParser_Literals(t *testing.T) {
	tests := []parserTestItem{
		newParserTestItem("integer literal simple five", "5", newAstBuilder().addIntegerLiteralExpression(5).toProgram()),
		newParserTestItem("integer literal simple zero", "0", newAstBuilder().addIntegerLiteralExpression(0).toProgram()),
		newParserTestItem("integer literal large", "99999", newAstBuilder().addIntegerLiteralExpression(99999).toProgram()),
		newParserTestItem("integer literal negative", "-42", newAstBuilder().addUnaryExpression(ast.NewIntegerExpr(42)).toProgram()),
		newParserTestItem("float literal simple", "3.14", newAstBuilder().addFloatLiteralExpression(3.14).toProgram()),
		newParserTestItem("float literal without leading digit", ".314", newAstBuilder().addFloatLiteralExpression(.314).toProgram()),
		newParserTestItem("float literal without trailing digit", "314.", newAstBuilder().addFloatLiteralExpression(314.).toProgram()),
		newParserTestItem("bool literal simple true", "true", newAstBuilder().addBooleanLiteralExpression(true).toProgram()),
		newParserTestItem("bool literal simple false", "false", newAstBuilder().addBooleanLiteralExpression(false).toProgram()),
		newParserTestItem("string literal simple", `"hello"`, newAstBuilder().addStringLiteralExpression(`"hello"`).toProgram()),
		newParserTestItem("string literal newline", `"hello\nworld"`, newAstBuilder().addStringLiteralExpression(`"hello\nworld"`).toProgram()),
		newParserTestItem("string literal tab", `"hello\tworld"`, newAstBuilder().addStringLiteralExpression(`"hello\tworld"`).toProgram()),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParserTest(t, tt)
		})
	}
}
