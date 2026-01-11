package typechecker

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/0xmukesh/coco/internal/types"
)

type TypeChecker struct {
	env    *types.TypeEnvironment
	errors []string
}

func New(env *types.TypeEnvironment) *TypeChecker {
	return &TypeChecker{
		env:    env,
		errors: []string{},
	}
}

func (tc *TypeChecker) checkBinaryExpression(expr *ast.BinaryExpression) types.Type {
	leftType := tc.checkExpression(expr.Left)
	rightType := tc.checkExpression(expr.Right)
	leftTypeCategory := types.GetTypeCategory(leftType)
	rightTypeCategory := types.GetTypeCategory(rightType)

	op := expr.Operator.Type
	isComparisonOperator := op == tokens.LESS_THAN || op == tokens.GREATER_THAN || op == tokens.LESS_THAN_EQUALS || op == tokens.GREATER_THAN_EQUALS || op == tokens.EQUALS || op == tokens.NOT_EQUALS

	// numeric types (int, float)
	if leftTypeCategory == types.CategoryNumeric && rightTypeCategory == types.CategoryNumeric {
		// arithmetic operators
		if op == tokens.PLUS || op == tokens.MINUS || op == tokens.STAR || op == tokens.SLASH || op == tokens.DOUBLE_STAR {
			// if it is numeric arithmetic and either one of them is float, then result type is float
			// and the one which is integer is converted to float expression
			if leftType.Equals(types.FloatType{}) || rightType.Equals(types.FloatType{}) {
				if leftIntLit, ok := expr.Left.(*ast.IntegerExpression); ok {
					expr.Left = &ast.FloatExpression{
						Token: leftIntLit.Token,
						Value: float64(leftIntLit.Value),
						Type:  types.FloatType{},
					}
				}

				if rightIntLit, ok := expr.Right.(*ast.IntegerExpression); ok {
					expr.Right = &ast.FloatExpression{
						Token: rightIntLit.Token,
						Value: float64(rightIntLit.Value),
						Type:  types.FloatType{},
					}
				}

				return expr.SetType(types.FloatType{})
			} else {
				return expr.SetType(types.IntType{})
			}
		}

		// comparison operators
		if isComparisonOperator {
			return expr.SetType(types.BoolType{})
		}
	}

	// strings
	if leftType.Equals(types.StringType{}) && rightType.Equals(types.StringType{}) {
		// string concatenation
		if op == tokens.PLUS {
			return expr.SetType(types.StringType{})
		}

		// lexicographical comparison
		if isComparisonOperator {
			return expr.SetType(types.BoolType{})
		}
	}

	// bools
	if leftType.Equals(types.BoolType{}) && rightType.Equals(types.BoolType{}) {
		if isComparisonOperator {
			return expr.SetType(types.BoolType{})
		}
	}

	tc.addError("cannot perform %s operation on %s and %s", op, leftType, rightType)
	return leftType
}

func (tc *TypeChecker) checkExpression(expr ast.Expression) types.Type {
	switch e := expr.(type) {
	case *ast.IntegerExpression:
		return types.IntType{}
	case *ast.StringExpression:
		return types.StringType{}
	case *ast.BooleanExpression:
		return types.BoolType{}
	case *ast.BinaryExpression:
		return tc.checkBinaryExpression(e)
	default:
		tc.addError("unknown expression of type %T", expr)
		return nil
	}
}

func (tc *TypeChecker) checkStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		tc.checkExpression(s.Expr)
	}
}

func (tc *TypeChecker) Transform(program *ast.Program) *ast.Program {
	for _, stmt := range program.Statements {
		tc.checkStatement(stmt)
	}

	return program
}

func (tc *TypeChecker) addError(format string, a ...any) {
	tc.errors = append(tc.errors, fmt.Sprintf(format, a...))
}

func (tc *TypeChecker) Errors() []string {
	return tc.errors
}
func (tc *TypeChecker) HasErrors() bool {
	return len(tc.errors) > 0
}
