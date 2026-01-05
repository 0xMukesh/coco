package eval

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/object"
	"github.com/0xmukesh/coco/internal/tokens"
)

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)

		if result != nil && (result.Type() == object.RETURN_OBJECT || result.Type() == object.ERROR_OBJECT) {
			return result
		}
	}

	return result
}

func evalUnaryExpression(operator tokens.TokenType, right object.Object) object.Object {
	switch operator {
	case tokens.BANG:
		return evalBangOperator(right)
	case tokens.MINUS:
		return evalMinusPrefixOperator(right)
	default:
		return newErrorObject("invalid operator for unary expression: %s", operator)
	}
}

func evalBangOperator(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	}

	return FALSE
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != object.INT_OBJECT || right.Type() != object.FLOAT_OBJECT {
		return newErrorObject("unknown operator: -%s", right.Type())
	}

	if right.Type() == object.INT_OBJECT {
		value := right.(*object.Integer).Value
		return &object.Integer{
			Value: -value,
		}
	} else {
		value := right.(*object.Float).Value
		return &object.Float{
			Value: -value,
		}
	}
}

func evalBinaryExpression(operator tokens.TokenType, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newErrorObject("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INT_OBJECT && right.Type() == object.INT_OBJECT:
		return evalIntegerBinaryExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJECT && right.Type() == object.FLOAT_OBJECT:
		return evalFloatBinaryExpression(operator, left, right)
	case operator == tokens.EQUALS:
		return nativeBoolToObjectBool(left == right)
	case operator == tokens.NOT_EQUALS:
		return nativeBoolToObjectBool(left != right)
	default:
		return newErrorObject("unknown operator: %s", operator)
	}
}

func evalIntegerBinaryExpression(operator tokens.TokenType, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case tokens.PLUS:
		return &object.Integer{
			Value: leftValue + rightValue,
		}
	case tokens.MINUS:
		return &object.Integer{
			Value: leftValue - rightValue,
		}
	case tokens.STAR:
		return &object.Integer{
			Value: leftValue * rightValue,
		}
	case tokens.SLASH:
		if rightValue == 0 {
			return newErrorObject("zero division")
		}

		return &object.Integer{
			Value: leftValue / rightValue,
		}
	case tokens.LESS_THAN:
		return &object.Boolean{
			Value: leftValue < rightValue,
		}
	case tokens.GREATER_THAN:
		return &object.Boolean{
			Value: leftValue > rightValue,
		}
	case tokens.LESS_THAN_EQUAL:
		return &object.Boolean{
			Value: leftValue <= rightValue,
		}
	case tokens.GREATER_THAN_EQUAL:
		return &object.Boolean{
			Value: leftValue >= rightValue,
		}
	default:
		return newErrorObject("invalid operator for binary expression: %s", operator)
	}
}

func evalFloatBinaryExpression(operator tokens.TokenType, left, right object.Object) object.Object {
	leftValue := left.(*object.Float).Value
	rightValue := right.(*object.Float).Value

	switch operator {
	case tokens.PLUS:
		return &object.Float{
			Value: leftValue + rightValue,
		}
	case tokens.MINUS:
		return &object.Float{
			Value: leftValue - rightValue,
		}
	case tokens.STAR:
		return &object.Float{
			Value: leftValue * rightValue,
		}
	case tokens.SLASH:
		if rightValue == 0 {
			return newErrorObject("zero division")
		}

		return &object.Float{
			Value: leftValue / rightValue,
		}
	case tokens.LESS_THAN:
		return &object.Boolean{
			Value: leftValue < rightValue,
		}
	case tokens.GREATER_THAN:
		return &object.Boolean{
			Value: leftValue > rightValue,
		}
	case tokens.LESS_THAN_EQUAL:
		return &object.Boolean{
			Value: leftValue <= rightValue,
		}
	case tokens.GREATER_THAN_EQUAL:
		return &object.Boolean{
			Value: leftValue >= rightValue,
		}
	default:
		return newErrorObject("invalid operator for binary expression: %s", operator)
	}
}

func evalIfExpression(expression *ast.IfExpression) object.Object {
	condition := Eval(expression.Condition)

	if isTruthy(condition) {
		return Eval(expression.Consequence)
	} else if expression.Alternative != nil {
		return Eval(expression.Alternative)
	} else {
		return NULL
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToObjectBool(node.Value)
	case *ast.UnaryExpression:
		right := Eval(node.Expression)
		if isError(right) {
			return right
		}

		return evalUnaryExpression(node.Token.Type, right)
	case *ast.BinaryExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalBinaryExpression(node.Token.Type, left, right)
	case *ast.ReturnStatement:
		val := Eval(node.Expression)
		if isError(val) {
			return val
		}

		return &object.Return{
			Value: val,
		}
	}

	return nil
}
