package ast

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
