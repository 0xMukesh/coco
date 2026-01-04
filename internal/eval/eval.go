package eval

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/object"
	"github.com/0xmukesh/coco/internal/tokens"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{
		Value: true,
	}
	FALSE = &object.Boolean{
		Value: false,
	}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToObjectBool(node.Value)
	case *ast.UnaryExpression:
		right := Eval(node.Expression)
		return evalUnaryExpression(node.Token.Type, right)
	case *ast.BinaryExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)

		return evalBinaryExpression(node.Token.Type, left, right)
	}

	return nil
}

func evalStatements(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)
	}

	return result
}

func nativeBoolToObjectBool(value bool) object.Object {
	if value == true {
		return TRUE
	} else {
		return FALSE
	}
}

func evalUnaryExpression(operator tokens.TokenType, right object.Object) object.Object {
	switch operator {
	case tokens.BANG:
		return evalBangOperator(right)
	case tokens.MINUS:
		return evalMinusPrefixOperator(right)
	default:
		return NULL // FIXME: add proper error handling
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
	if right.Type() != object.OBJECT_INT {
		return NULL // FIXME: add proper error handling
	}

	value := right.(*object.Integer).Value
	return &object.Integer{
		Value: -value,
	}
}

func evalBinaryExpression(operator tokens.TokenType, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.OBJECT_INT && right.Type() == object.OBJECT_INT:
		return evalIntegerBinaryExpression(operator, left, right)
	case operator == tokens.EQUALS:
		return nativeBoolToObjectBool(left == right)
	case operator == tokens.NOT_EQUALS:
		return nativeBoolToObjectBool(left != right)
	default:
		return NULL
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
		return NULL
	}
}
