package typechecker

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

func (tc *TypeChecker) checkExpression(expr ast.Expression) (t cotypes.Type, err error) {
	switch e := expr.(type) {
	case *ast.IntegerExpression:
		t = cotypes.IntType{}
	case *ast.FloatExpression:
		t = cotypes.FloatType{}
	case *ast.BooleanExpression:
		t = cotypes.BoolType{}
	case *ast.StringExpression:
		t = cotypes.StringType{}
	case *ast.IdentifierExpression:
		identType, found := tc.env.Get(e.String())

		if !found {
			err = tc.propagateOrWrapError(nil, expr, "unknown identifier: %s", e.String())
		} else {
			t = identType
		}
	case *ast.BinaryExpression:
		t, err = tc.checkBinaryExpression(e)
	case *ast.CallExpression:
		t, err = tc.checkCallExpression(e)
	case *ast.GroupedExpression:
		t, err = tc.checkExpression(e.Expr)
	case *ast.IfExpression:
		t, err = tc.checkIfExpression(e)
	default:
		err = tc.propagateOrWrapError(nil, expr, "unknown expression type: %T", e)
	}

	if err != nil {
		tc.propagateOrWrapError(err, expr, "failed to typecheck expression of type %T", expr)
	} else {
		expr.SetType(t)
	}

	return t, err
}

func (tc *TypeChecker) checkBinaryExpression(expr *ast.BinaryExpression) (cotypes.Type, error) {
	leftType, err := tc.checkExpression(expr.Left)
	if err != nil {
		return nil, tc.propagateOrWrapError(err, expr, "failed to type check left operand: %s", err.Error())
	}

	rightType, err := tc.checkExpression(expr.Right)
	if err != nil {
		return nil, tc.propagateOrWrapError(err, expr, "failed to type check right operand: %s", err.Error())
	}

	leftTypeCategory := cotypes.GetTypeCategory(leftType)
	rightTypeCategory := cotypes.GetTypeCategory(rightType)

	if leftType == nil || rightType == nil {
		return nil, err
	}

	op := expr.Operator.Type
	isComparisonOperator := op == tokens.LESS_THAN ||
		op == tokens.GREATER_THAN ||
		op == tokens.LESS_THAN_EQUALS ||
		op == tokens.GREATER_THAN_EQUALS ||
		op == tokens.EQUALS ||
		op == tokens.NOT_EQUALS

	// numeric types (int, float)
	if leftTypeCategory == cotypes.CategoryNumeric && rightTypeCategory == cotypes.CategoryNumeric {
		leftIsLiteral := tc.isLiteralExpression(expr.Left)
		rightIsLiteral := tc.isLiteralExpression(expr.Right)

		leftIsInt := leftType.Equals(cotypes.IntType{})
		leftIsFloat := leftType.Equals(cotypes.FloatType{})
		rightIsInt := rightType.Equals(cotypes.IntType{})
		rightIsFloat := rightType.Equals(cotypes.FloatType{})

		// arithmetic operators
		if op == tokens.PLUS || op == tokens.MINUS || op == tokens.STAR || op == tokens.SLASH || op == tokens.DOUBLE_STAR {
			// if type of both the operands are the same, then the final type of the expression would be equal to the same type
			if leftType.Equals(rightType) {
				return expr.SetType(leftType), nil
			}

			// if left and right have different type i.e. int and float then the same type would be float
			// and coercions would be used to convert int to float
			if (leftIsInt && rightIsFloat) || (leftIsFloat && rightIsInt) {
				if leftIsInt && leftIsLiteral {
					if coercible, ok := expr.Left.(ast.CoercibleExpression); ok {
						t := expr.SetType(cotypes.FloatType{})
						coercible.SetCoercion(cotypes.FloatType{})
						return t, nil
					}
				}

				if rightIsInt && rightIsLiteral {
					if coercible, ok := expr.Right.(ast.CoercibleExpression); ok {
						t := expr.SetType(cotypes.FloatType{})
						coercible.SetCoercion(cotypes.FloatType{})
						return t, nil
					}
				}

				return nil, tc.propagateOrWrapError(nil, expr, "cannot perform %s operation on %s and %s (implict coercion works for literals only)", op, leftType, rightType)
			}
		}

		// comparison operators
		if isComparisonOperator {
			return expr.SetType(cotypes.BoolType{}), err
		}
	}

	// strings
	if leftType.Equals(cotypes.StringType{}) && rightType.Equals(cotypes.StringType{}) {
		// string concatenation
		if op == tokens.PLUS {
			return expr.SetType(cotypes.StringType{}), err
		}

		// lexicographical comparison
		if isComparisonOperator {
			return expr.SetType(cotypes.BoolType{}), err
		}
	}

	// bools
	if leftType.Equals(cotypes.BoolType{}) && rightType.Equals(cotypes.BoolType{}) {
		if isComparisonOperator {
			return expr.SetType(cotypes.BoolType{}), err
		}
	}

	return nil, tc.propagateOrWrapError(nil, expr, "cannot perform %s operation on %s and %s types", op, leftType, rightType)
}

func (tc *TypeChecker) checkCallExpression(expr *ast.CallExpression) (cotypes.Type, error) {
	// TODO: at the moment, only builtin functions are allowed
	if builtin, isBuiltin := tc.builtins[expr.Identifier.String()]; isBuiltin {
		expr.IsBuiltin = true
		expr.BuiltinKind = &builtin.kind

		return builtin.checker(expr)
	}

	return nil, tc.propagateOrWrapError(nil, expr, "cannot call %s identifier", expr.Identifier.String())
}

func (tc *TypeChecker) checkIfExpression(expr *ast.IfExpression) (cotypes.Type, error) {
	conditionType, err := tc.checkExpression(expr.Condition)
	if err != nil {
		return nil, tc.propagateOrWrapError(err, expr, "failed to type check if branch condition expression: %s", err.Error())
	}

	if !conditionType.Equals(cotypes.BoolType{}) {
		return nil, tc.propagateOrWrapError(nil, expr, "non-boolean condition if if expression")
	}

	var (
		consequenceReturnType cotypes.Type = nil
		alternativeReturnType cotypes.Type = nil
	)

	for _, stmt := range expr.Consequence.Statements {
		tc.checkStatement(expr.Consequence)
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if ok {
			consequenceReturnType = returnStmt.Expr.GetType()
		}
	}

	if expr.Alternative != nil {
		for _, stmt := range expr.Alternative.Statements {
			tc.checkStatement(expr.Alternative)
			returnStmt, ok := stmt.(*ast.ReturnStatement)

			if ok {
				alternativeReturnType = returnStmt.Expr.GetType()
			}
		}
	}

	if consequenceReturnType != nil && alternativeReturnType != nil && consequenceReturnType.Equals(alternativeReturnType) {
		return consequenceReturnType, nil
	}

	if consequenceReturnType == nil && alternativeReturnType == nil {
		return cotypes.VoidType{}, nil
	}

	return nil, tc.propagateOrWrapError(nil, expr, "return types of if and else blocks need to be equal. got %q from if-block and %q from else-block", consequenceReturnType, alternativeReturnType)
}
