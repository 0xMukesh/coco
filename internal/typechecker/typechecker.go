package typechecker

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

type TypeChecker struct {
	env      *cotypes.TypeEnvironment
	builtins map[string]*builtinsInfo
	errors   []string
}

func New(env *cotypes.TypeEnvironment) *TypeChecker {
	tc := &TypeChecker{
		env:      env,
		errors:   []string{},
		builtins: make(map[string]*builtinsInfo),
	}

	tc.registerBuiltins()

	return tc
}

func (tc *TypeChecker) checkBinaryExpression(expr *ast.BinaryExpression) cotypes.Type {
	leftType := tc.checkExpression(expr.Left)
	rightType := tc.checkExpression(expr.Right)
	leftTypeCategory := cotypes.GetTypeCategory(leftType)
	rightTypeCategory := cotypes.GetTypeCategory(rightType)

	if leftType == nil || rightType == nil {
		return nil
	}

	op := expr.Operator.Type
	isComparisonOperator := op == tokens.LESS_THAN || op == tokens.GREATER_THAN || op == tokens.LESS_THAN_EQUALS || op == tokens.GREATER_THAN_EQUALS || op == tokens.EQUALS || op == tokens.NOT_EQUALS

	// numeric types (int, float)
	if leftTypeCategory == cotypes.CategoryNumeric && rightTypeCategory == cotypes.CategoryNumeric {
		// arithmetic operators
		if op == tokens.PLUS || op == tokens.MINUS || op == tokens.STAR || op == tokens.SLASH || op == tokens.DOUBLE_STAR {
			// if it is numeric arithmetic and either one of them is float, then result type is float
			// and the one which is integer is converted to float expression
			if leftType.Equals(cotypes.FloatType{}) || rightType.Equals(cotypes.FloatType{}) {
				if leftIntLit, ok := expr.Left.(*ast.IntegerExpression); ok {
					expr.Left = &ast.FloatExpression{
						Token: leftIntLit.Token,
						Value: float64(leftIntLit.Value),
						Type:  cotypes.FloatType{},
					}
				}

				if rightIntLit, ok := expr.Right.(*ast.IntegerExpression); ok {
					expr.Right = &ast.FloatExpression{
						Token: rightIntLit.Token,
						Value: float64(rightIntLit.Value),
						Type:  cotypes.FloatType{},
					}
				}

				return expr.SetType(cotypes.FloatType{})
			} else {
				return expr.SetType(cotypes.IntType{})
			}
		}

		// comparison operators
		if isComparisonOperator {
			return expr.SetType(cotypes.BoolType{})
		}
	}

	// strings
	if leftType.Equals(cotypes.StringType{}) && rightType.Equals(cotypes.StringType{}) {
		// string concatenation
		if op == tokens.PLUS {
			return expr.SetType(cotypes.StringType{})
		}

		// lexicographical comparison
		if isComparisonOperator {
			return expr.SetType(cotypes.BoolType{})
		}
	}

	// bools
	if leftType.Equals(cotypes.BoolType{}) && rightType.Equals(cotypes.BoolType{}) {
		if isComparisonOperator {
			return expr.SetType(cotypes.BoolType{})
		}
	}

	tc.addError("cannot perform %s operation on %s and %s", op, leftType, rightType)
	return nil
}

func (tc *TypeChecker) checkCallExpression(expr *ast.CallExpression) cotypes.Type {
	if builtin, isBuiltin := tc.builtins[expr.Identifier.String()]; isBuiltin {
		expr.IsBuiltin = true
		expr.BuiltinKind = &builtin.kind
		return builtin.checker(expr)
	}

	// TODO: only builtin functions (just print) can be called
	tc.addError("cannot call %s identifier", expr.Identifier.String())
	return nil
}

func (tc *TypeChecker) checkPrintBuiltin(expr *ast.CallExpression) cotypes.Type {
	for i, arg := range expr.Arguments {
		argType := tc.checkExpression(arg)
		arg.SetType(argType)

		// TODO: only integers are supported for printing
		if !argType.Equals(cotypes.IntType{}) {
			tc.addError("invalid argument at %d idx to print", i)
			return nil
		}
	}

	return cotypes.VoidType{}
}

func (tc *TypeChecker) checkExpression(expr ast.Expression) cotypes.Type {
	var t cotypes.Type = nil

	switch e := expr.(type) {
	case *ast.IntegerExpression:
		t = cotypes.IntType{}
	case *ast.FloatExpression:
		t = cotypes.FloatType{}
	case *ast.StringExpression:
		t = cotypes.StringType{}
	case *ast.BooleanExpression:
		t = cotypes.BoolType{}
	case *ast.IdentifierExpression:
		identType, found := tc.env.Get(e.String())
		if !found {
			tc.addError("unknown identifier: %s", e.String())
		} else {
			t = identType
		}
	case *ast.BinaryExpression:
		t = tc.checkBinaryExpression(e)
	case *ast.CallExpression:
		t = tc.checkCallExpression(e)
	default:
		tc.addError("unknown expression of type %T", expr)
	}

	expr.SetType(t)
	return t
}

func (tc *TypeChecker) checkExitStatement(stmt *ast.ExitStatement) {
	exprType := tc.checkExpression(stmt.Expr)

	if exprType != nil && !exprType.Equals(cotypes.IntType{}) {
		tc.addError("expected exit code to an integer, got %q", exprType)
	}
}

func (tc *TypeChecker) checkStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		tc.checkExpression(s.Expr)
	case *ast.LetStatement:
		varType := tc.checkExpression(s.Value)
		tc.env.Set(s.Identifier.String(), varType)
	case *ast.AssignmentStatement:
		_, exists := tc.env.Get(s.Identifier.String())
		if !exists {
			tc.addError("unknown identifier: %s", s.Identifier.String())
		}

		tc.checkExpression(s.Value)
	case *ast.ExitStatement:
		tc.checkExitStatement(s)
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
