package ast

import "github.com/0xmukesh/coco/internal/tokens"

type BuiltinsKind int

const (
	BuiltinFuncPrint BuiltinsKind = iota
	BuiltinFuncExit
)

func NewIntegerExpr(value int64) Expression {
	return &IntegerExpression{
		Value: value,
	}
}

func NewFloatExpr(value float64) Expression {
	return &FloatExpression{
		Value: value,
	}
}

func NewBooleanExpr(value bool) Expression {
	return &BooleanExpression{
		Value: value,
	}
}

func NewStringExpr(value string) Expression {
	return &StringExpression{
		Value: value,
	}
}

func NewIdentifierExpr(literal string) Expression {
	return &IdentifierExpression{
		Literal: literal,
	}
}

func NewUnaryExpr(operator tokens.Token, expr Expression) Expression {
	return &UnaryExpression{
		Token: operator,
		Expr:  expr,
	}
}

func NewBinaryExpr(operator tokens.Token, left, right Expression) Expression {
	return &BinaryExpression{
		Operator: operator,
		Left:     left,
		Right:    right,
	}
}

func NewGroupedExpr(expr Expression) Expression {
	return &GroupedExpression{
		Expr: expr,
	}
}

func WrapExprsAsStmts(exprs []Expression) []Statement {
	stmts := []Statement{}

	for _, e := range exprs {
		stmts = append(stmts, &ExpressionStatement{
			Expr: e,
		})
	}

	return stmts
}

func NewIfExpr(condition Expression, consequence []Statement, alternative []Statement) Expression {
	expr := &IfExpression{
		Condition: condition,
		Consequence: &BlockStatement{
			Statements: consequence,
		},
	}

	if alternative != nil {
		expr.Alternative = &BlockStatement{
			Statements: alternative,
		}
	}

	return expr
}
