package parser

import (
	"testing"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
)

func TestParser_Literals(t *testing.T) {
	tests := []parserTestItem{
		newParserTest("integer literal simple five", "5", newAstBuilder().addIntegerLiteralExpression(5).toProgram()),
		newParserTest("integer literal simple zero", "0", newAstBuilder().addIntegerLiteralExpression(0).toProgram()),
		newParserTest("integer literal large", "99999", newAstBuilder().addIntegerLiteralExpression(99999).toProgram()),
		newParserTest("float literal simple", "3.14", newAstBuilder().addFloatLiteralExpression(3.14).toProgram()),
		newParserTest("float literal without leading digit", ".314", newAstBuilder().addFloatLiteralExpression(.314).toProgram()),
		newParserTest("float literal without trailing digit", "314.", newAstBuilder().addFloatLiteralExpression(314.).toProgram()),
		newParserTest("bool literal simple true", "true", newAstBuilder().addBooleanLiteralExpression(true).toProgram()),
		newParserTest("bool literal simple false", "false", newAstBuilder().addBooleanLiteralExpression(false).toProgram()),
		newParserTest("string literal simple", `"hello"`, newAstBuilder().addStringLiteralExpression(`"hello"`).toProgram()),
		newParserTest("string literal newline", `"hello\nworld"`, newAstBuilder().addStringLiteralExpression(`"hello\nworld"`).toProgram()),
		newParserTest("string literal tab", `"hello\tworld"`, newAstBuilder().addStringLiteralExpression(`"hello\tworld"`).toProgram()),
		newParserTest("string literal escape newline", `"hello\\nworld"`, newAstBuilder().addStringLiteralExpression(`"hello\\nworld"`).toProgram()),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParserTest(t, tt)
		})
	}
}

func TestParser_Identifiers(t *testing.T) {
	tests := []parserTestItem{
		newParserTest("simple", "x", newAstBuilder().addIdentifierExpression("x").toProgram()),
		newParserTest("with underscore", "xy_ab", newAstBuilder().addIdentifierExpression("xy_ab").toProgram()),
		newParserTest("with numbers at the end", "xy123", newAstBuilder().addIdentifierExpression("xy123").toProgram()),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParserTest(t, tt)
		})
	}
}

func TestParser_UnaryExpressions(t *testing.T) {
	tests := []parserTestItem{
		newParserTest("single bang", "!true", newAstBuilder().addUnaryExpression(tokens.NewMinimal(tokens.BANG, "!"), ast.NewBooleanExpr(true)).toProgram()),
		newParserTest(
			"double bang",
			"!!true",
			newAstBuilder().addUnaryExpression(
				tokens.NewMinimal(tokens.BANG, "!"),
				ast.NewUnaryExpr(
					tokens.NewMinimal(tokens.BANG, "!"),
					ast.NewBooleanExpr(true),
				),
			).toProgram(),
		),
		newParserTest(
			"triple bang",
			"!!!true",
			newAstBuilder().addUnaryExpression(
				tokens.NewMinimal(tokens.BANG, "!"),
				ast.NewUnaryExpr(
					tokens.NewMinimal(tokens.BANG, "!"),
					ast.NewUnaryExpr(
						tokens.NewMinimal(tokens.BANG, "!"),
						ast.NewBooleanExpr(true),
					),
				),
			).toProgram(),
		),
		newParserTest("single minus", "-5", newAstBuilder().addUnaryExpression(tokens.NewMinimal(tokens.MINUS, "-"), ast.NewIntegerExpr(5)).toProgram()),
		newParserTestFail("invalid unary bang", "!", expectParseFailure("expression expected after ! token")),
		newParserTestFail("invalid unary minus", "-", expectParseFailure("expression expected after - token")),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParserTest(t, tt)
		})
	}
}

func TestParser_BinaryExpressions(t *testing.T) {
	// generates ast for "1 <operator> 2" style binary expressions
	astGenerator := func(tokenType tokens.TokenType, tokenLiteral string) *ast.Program {
		return newAstBuilder().addBinaryExpression(
			tokens.NewMinimal(tokenType, tokenLiteral),
			ast.NewIntegerExpr(1),
			ast.NewIntegerExpr(2),
		).toProgram()
	}

	tests := []parserTestItem{
		newParserTest("add", "1 + 2", astGenerator(tokens.PLUS, "+")),
		newParserTest("subtract", "1 - 2", astGenerator(tokens.MINUS, "-")),
		newParserTest("multiply", "1 * 2", astGenerator(tokens.STAR, "*")),
		newParserTest("divide", "1 / 2", astGenerator(tokens.SLASH, "/")),
		newParserTest("modulo", "1 % 2", astGenerator(tokens.MODULO, "%")),
		newParserTest("equals", "1 == 2", astGenerator(tokens.EQUALS, "==")),
		newParserTest("not equals", "1 != 2", astGenerator(tokens.NOT_EQUALS, "!=")),
		newParserTest("less than", "1 < 2", astGenerator(tokens.LESS_THAN, "<")),
		newParserTest("less than equals", "1 <= 2", astGenerator(tokens.LESS_THAN_EQUALS, "<=")),
		newParserTest("greater than", "1 > 2", astGenerator(tokens.GREATER_THAN, ">")),
		newParserTest("greater than equals", "1 >= 2", astGenerator(tokens.GREATER_THAN_EQUALS, ">=")),
		newParserTest("and", "1 && 2", astGenerator(tokens.AND, "&&")),
		newParserTest("or", "1 || 2", astGenerator(tokens.OR, "||")),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParserTest(t, tt)
		})
	}
}

func TestParser_Precedence(t *testing.T) {
	tests := []parserTestItem{
		// 5 + 3 * 2 = 5 + (3 * 2)
		newParserTest(
			"simple",
			"5 + 3 * 2",
			newAstBuilder().addBinaryExpression(
				tokens.NewMinimal(tokens.PLUS, "+"),
				ast.NewIntegerExpr(5),
				ast.NewBinaryExpr(
					tokens.NewMinimal(tokens.STAR, "*"),
					ast.NewIntegerExpr(3),
					ast.NewIntegerExpr(2),
				),
			).toProgram(),
		),
		// 2 * 3 + 4 = (2 * 3) + 4
		newParserTest(
			"simple",
			"2 * 3 + 4",
			newAstBuilder().addBinaryExpression(
				tokens.NewMinimal(tokens.PLUS, "+"),
				ast.NewBinaryExpr(
					tokens.NewMinimal(tokens.STAR, "*"),
					ast.NewIntegerExpr(2),
					ast.NewIntegerExpr(3),
				),
				ast.NewIntegerExpr(4),
			).toProgram(),
		),
		// 10 - 2 - 3 = (10 - 2) - 3
		newParserTest(
			"left associative",
			"10 - 2 - 3",
			newAstBuilder().addBinaryExpression(
				tokens.NewMinimal(tokens.MINUS, "-"),
				ast.NewBinaryExpr(
					tokens.NewMinimal(tokens.MINUS, "-"),
					ast.NewIntegerExpr(10),
					ast.NewIntegerExpr(2),
				),
				ast.NewIntegerExpr(3),
			).toProgram(),
		),
		// 5 == 5 && 3 != 2 == (5 == 5) && (3 != 2)
		newParserTest(
			"logical + comparison",
			"5 == 5 && 3 != 2",
			newAstBuilder().addBinaryExpression(
				tokens.NewMinimal(tokens.AND, "&&"),
				ast.NewBinaryExpr(tokens.NewMinimal(tokens.EQUALS, "=="), ast.NewIntegerExpr(5), ast.NewIntegerExpr(5)),
				ast.NewBinaryExpr(tokens.NewMinimal(tokens.NOT_EQUALS, "!="), ast.NewIntegerExpr(3), ast.NewIntegerExpr(2)),
			).toProgram(),
		),
		// -5 + 3 = (-5) + 3
		newParserTest(
			"unary + binary",
			"-5 + 3",
			newAstBuilder().addBinaryExpression(
				tokens.NewMinimal(tokens.PLUS, "+"),
				ast.NewUnaryExpr(tokens.NewMinimal(tokens.MINUS, "-"), ast.NewIntegerExpr(5)),
				ast.NewIntegerExpr(3),
			).toProgram(),
		),
		// 5 + 3 * 2 - 4 / 2 = (5 + (3 * 2)) - (4 / 2)
		newParserTest(
			"complex",
			"5 + 3 * 2 - 4 / 2",
			newAstBuilder().addBinaryExpression(
				tokens.NewMinimal(tokens.MINUS, "-"),
				ast.NewBinaryExpr(
					tokens.NewMinimal(tokens.PLUS, "+"),
					ast.NewIntegerExpr(5),
					ast.NewBinaryExpr(
						tokens.NewMinimal(tokens.STAR, "*"),
						ast.NewIntegerExpr(3),
						ast.NewIntegerExpr(2),
					),
				),
				ast.NewBinaryExpr(
					tokens.NewMinimal(tokens.SLASH, "/"),
					ast.NewIntegerExpr(4),
					ast.NewIntegerExpr(2),
				),
			).toProgram(),
		),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParserTest(t, tt)
		})
	}
}

func TestParser_GroupedExpressions(t *testing.T) {
	tests := []parserTestItem{
		newParserTest("simple grouped", "(5)", newAstBuilder().addGroupedExpression(ast.NewIntegerExpr(5)).toProgram()),
		newParserTest(
			"grouped + binary",
			"(5 + 3)",
			newAstBuilder().addGroupedExpression(
				ast.NewBinaryExpr(tokens.NewMinimal(tokens.PLUS, "+"), ast.NewIntegerExpr(5), ast.NewIntegerExpr(3)),
			).toProgram(),
		),
		newParserTest(
			"nested grouping",
			"((5 + 3) * 2) / 4",
			newAstBuilder().addBinaryExpression(
				tokens.NewMinimal(tokens.SLASH, "/"),
				ast.NewGroupedExpr(
					ast.NewBinaryExpr(
						tokens.NewMinimal(tokens.STAR, "*"),
						ast.NewGroupedExpr(ast.NewBinaryExpr(tokens.NewMinimal(tokens.PLUS, "+"), ast.NewIntegerExpr(5), ast.NewIntegerExpr(3))),
						ast.NewIntegerExpr(2),
					),
				),
				ast.NewIntegerExpr(4),
			).toProgram(),
		),
		newParserTestFail(
			"missing closing paren",
			"(5 + 3",
			expectParseFailure("expected type of next token to be ), got EOF instead"),
		),
		newParserTestFail(
			"missing open paren",
			"5 + 3)",
			expectParseFailure("no prefix function found for ) token"),
		),
		newParserTestFail(
			"empty paren",
			"()",
			expectParseFailure("expression expected after ( token"),
		),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParserTest(t, tt)
		})
	}
}
