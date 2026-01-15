package ast

import "github.com/0xmukesh/coco/internal/tokens"

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
