package eval

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/object"
	"github.com/0xmukesh/coco/internal/tokens"
)

type Evalutor struct{}

func NewEvalutor() *Evalutor {
	return &Evalutor{}
}

func (e *Evalutor) Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToObjectBool(node.Value)
	case *ast.Identifier:
		return e.evalIdentifier(node, env)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)
	case *ast.BlockStatement:
		return e.evalStatements(node.Statements, env)
	case *ast.IfExpression:
		return e.evalIfExpression(node, env)
	case *ast.LetStatement:
		val := e.Eval(node.Value, env)
		if isError(val) {
			return val
		}

		env.Set(node.Identifier.Value, val)
	case *ast.UnaryExpression:
		right := e.Eval(node.Expression, env)
		if isError(right) {
			return right
		}

		return e.evalUnaryExpression(node.Token.Type, right)
	case *ast.BinaryExpression:
		left := e.Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := e.Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return e.evalBinaryExpression(node.Token.Type, left, right)
	case *ast.ReturnStatement:
		val := e.Eval(node.Expression, env)
		if isError(val) {
			return val
		}

		return &object.Return{
			Value: val,
		}
	case *ast.FunctionLiteral:
		parameters := node.Parameters
		body := node.Body

		return &object.Function{
			Parameters: parameters,
			Body:       body,
			Env:        env,
		}
	case *ast.CallExpression:
		return e.evalCallExpression(node, env)
	case *ast.AssignmentExpression:
		return e.evalAssignmentExpression(node, env)
	}

	return nil
}

func (e *Evalutor) evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = e.Eval(stmt, env)

		switch result := result.(type) {
		case *object.Return:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func (e *Evalutor) evalStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = e.Eval(stmt, env)

		if result != nil && (result.Type() == object.RETURN_OBJECT || result.Type() == object.ERROR_OBJECT) {
			return result
		}
	}

	return result
}

func (e *Evalutor) evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range expressions {
		evaled := e.Eval(expr, env)
		if isError(evaled) {
			return []object.Object{evaled}
		}

		result = append(result, evaled)
	}

	return result
}

func (e *Evalutor) evalUnaryExpression(operator tokens.TokenType, right object.Object) object.Object {
	switch operator {
	case tokens.BANG:
		return e.evalBangOperator(right)
	case tokens.MINUS:
		return e.evalMinusPrefixOperator(right)
	default:
		return newErrorObject("invalid operator for unary expression: %s", operator)
	}
}

func (e *Evalutor) evalBangOperator(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	}

	return FALSE
}

func (e *Evalutor) evalMinusPrefixOperator(right object.Object) object.Object {
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

func (e *Evalutor) evalBinaryExpression(operator tokens.TokenType, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newErrorObject("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INT_OBJECT && right.Type() == object.INT_OBJECT:
		return e.evalIntegerBinaryExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJECT && right.Type() == object.FLOAT_OBJECT:
		return e.evalFloatBinaryExpression(operator, left, right)
	case operator == tokens.EQUALS:
		return nativeBoolToObjectBool(left == right)
	case operator == tokens.NOT_EQUALS:
		return nativeBoolToObjectBool(left != right)
	default:
		return newErrorObject("unknown operator: %s", operator)
	}
}

func (e *Evalutor) evalIntegerBinaryExpression(operator tokens.TokenType, left, right object.Object) object.Object {
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

func (e *Evalutor) evalFloatBinaryExpression(operator tokens.TokenType, left, right object.Object) object.Object {
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

func (e *Evalutor) evalIfExpression(expression *ast.IfExpression, env *object.Environment) object.Object {
	condition := e.Eval(expression.Condition, env)

	if isTruthy(condition) {
		return e.Eval(expression.Consequence, env)
	} else if expression.Alternative != nil {
		return e.Eval(expression.Alternative, env)
	} else {
		return NULL
	}
}

func (e *Evalutor) evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}

	return newErrorObject("unknown identifier: %s", node.Value)
}

func (e *Evalutor) evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	functionObj := e.Eval(node.Function, env)
	if isError(functionObj) {
		return functionObj
	}

	function, ok := functionObj.(*object.Function)
	if !ok {
		return newErrorObject("invalid function object - %s", node.Function)
	}

	if len(function.Parameters) != len(node.Arguments) {
		return newErrorObject("invalid number of arguments for %s function. expected %d, got %d", node.Function, len(function.Parameters), len(node.Arguments))
	}

	args := e.evalExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	innerEnv := object.NewEnvironmentWithParent(function.Env)

	for i, param := range function.Parameters {
		innerEnv.Set(param.Value, args[i])
	}

	evaled := e.Eval(function.Body, innerEnv)
	returnObj, ok := evaled.(*object.Return)
	if ok {
		return returnObj.Value
	}

	return evaled
}

func (e *Evalutor) evalAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) object.Object {
	identifierObj := e.Eval(node.Identifier, env)
	if isError(identifierObj) {
		return identifierObj
	}

	value := e.Eval(node.Value, env)
	if isError(value) {
		return value
	}

	env.Set(node.Identifier.String(), value)
	identifierObj = e.Eval(node.Identifier, env)
	return identifierObj
}
