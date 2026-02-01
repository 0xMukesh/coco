package typechecker

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/env"
	"github.com/0xmukesh/coco/internal/tokens"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

type TypeEnvironment = *env.Environent[cotypes.Type]

type TypeChecker struct {
	env      TypeEnvironment
	builtins map[string]*builtinsInfo

	errors []error
}

func New() *TypeChecker {
	tc := &TypeChecker{
		env:      env.NewEnvironment[cotypes.Type](),
		builtins: make(map[string]*builtinsInfo),
		errors:   []error{},
	}

	tc.registerBuiltins()

	return tc
}

func (tc *TypeChecker) checkExpression(expr ast.Expression) (t cotypes.Type, err error) {
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
			err = fmt.Errorf("unknown identifier: %s", e.String())
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
		err = fmt.Errorf("unknown expression of type %T", expr)
	}

	if err != nil {
		tc.addError("%s", err.Error())
	} else {
		expr.SetType(t)
	}

	return
}

func (tc *TypeChecker) checkStatement(stmt ast.Statement) (err error) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		tc.checkExpression(s.Expr)
	case *ast.LetStatement:
		varName := s.Identifier.String()
		if tc.env.Has(varName) {
			tc.addError("cannot redeclare variable: %s", varName)
			return
		}

		varType, err := tc.checkExpression(s.Value)
		if err != nil {
			return err
		}

		tc.env.Set(s.Identifier.String(), varType)
	case *ast.AssignmentStatement:
		_, exists := tc.env.Get(s.Identifier.String())
		if !exists {
			err = fmt.Errorf("unknown identifier: %s", s.Identifier.String())
		}

		tc.checkExpression(s.Value)
	case *ast.BlockStatement:
		tc.env = env.NewEnvironmentWithParent(tc.env)
		for _, s := range s.Statements {
			tc.checkStatement(s)
		}

		tc.env = tc.env.Parent()
	}

	return
}

func (tc *TypeChecker) checkBinaryExpression(expr *ast.BinaryExpression) (t cotypes.Type, err error) {
	leftType, err := tc.checkExpression(expr.Left)
	if err != nil {
		return t, tc.propagateOrWrapError(err, expr, "failed to type check left operand: %s", err.Error())
	}

	rightType, err := tc.checkExpression(expr.Right)
	if err != nil {
		return t, tc.propagateOrWrapError(err, expr, "failed to type check right operand: %s", err.Error())
	}

	leftTypeCategory := cotypes.GetTypeCategory(leftType)
	rightTypeCategory := cotypes.GetTypeCategory(rightType)

	if leftType == nil || rightType == nil {
		return t, err
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

				return expr.SetType(cotypes.FloatType{}), err
			} else {
				return expr.SetType(cotypes.IntType{}), err
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

	err = fmt.Errorf("cannot perform %s operation on %s and %s", op, leftType, rightType)
	return
}

func (tc *TypeChecker) checkCallExpression(expr *ast.CallExpression) (t cotypes.Type, err error) {
	if builtin, isBuiltin := tc.builtins[expr.Identifier.String()]; isBuiltin {
		expr.IsBuiltin = true
		expr.BuiltinKind = &builtin.kind
		return builtin.checker(expr)
	}

	err = fmt.Errorf("cannot call %s identifier", expr.Identifier.String())
	return
}

func (tc *TypeChecker) checkIfExpression(expr *ast.IfExpression) (t cotypes.Type, err error) {
	conditionType, err := tc.checkExpression(expr.Condition)
	if err != nil {
		return t, tc.propagateOrWrapError(err, expr, "failed to type check if branch condition expression: %s", err.Error())
	}

	if !conditionType.Equals(cotypes.BoolType{}) {
		return t, fmt.Errorf("non-boolean condition if if expression")
	}

	tc.checkStatement(expr.Consequence)

	if expr.Alternative != nil {
		tc.checkStatement(expr.Alternative)
	}

	// TODO: need to proper handle the return type for if-expression after adding support for return statements
	return cotypes.VoidType{}, nil
}

func (tc *TypeChecker) checkPrintBuiltin(expr *ast.CallExpression) (t cotypes.Type, err error) {
	for i, arg := range expr.Arguments {
		argType, err := tc.checkExpression(arg)
		if err != nil {
			return t, tc.propagateOrWrapError(err, expr, "failed to type check print func arg at %d idx: %s", i, err.Error())
		}

		arg.SetType(argType)

		if !argType.Equals(cotypes.IntType{}) && !argType.Equals(cotypes.FloatType{}) && !argType.Equals(cotypes.BoolType{}) {
			return t, fmt.Errorf("invalid argument at %d idx to print", i)
		}
	}

	return cotypes.VoidType{}, nil
}

func (tc *TypeChecker) checkExitBuiltin(expr *ast.CallExpression) (t cotypes.Type, err error) {
	if len(expr.Arguments) != 1 {
		return t, fmt.Errorf("too many arguments. expected one argument, got %d arguments", len(expr.Arguments))
	}

	exitCode, err := tc.checkExpression(expr.Arguments[0])
	if err != nil {
		return t, tc.propagateOrWrapError(err, expr, "failed to type check exit func arg: %s", err.Error())
	}

	if !exitCode.Equals(cotypes.IntType{}) {
		return t, fmt.Errorf("expected exit code to be of type int, got %s", exitCode.String())
	}

	return cotypes.VoidType{}, nil
}

func (tc *TypeChecker) checkIntBuiltin(expr *ast.CallExpression) (t cotypes.Type, err error) {
	if len(expr.Arguments) != 1 {
		return t, fmt.Errorf("too many arguments. expected one argument, got %d arguments", len(expr.Arguments))
	}

	valType, err := tc.checkExpression(expr.Arguments[0])
	if err != nil {
		return t, tc.propagateOrWrapError(err, expr, "failed to type check int func arg: %s", err.Error())
	}

	if cotypes.GetTypeCategory(valType) != cotypes.CategoryNumeric {
		return t, fmt.Errorf("cannot convert %s to int", valType)
	}

	return cotypes.IntType{}, nil
}

func (tc *TypeChecker) checkFloatBuiltin(expr *ast.CallExpression) (t cotypes.Type, err error) {
	if len(expr.Arguments) != 1 {
		return t, fmt.Errorf("too many arguments. expected one argument, got %d arguments", len(expr.Arguments))
	}

	valType, err := tc.checkExpression(expr.Arguments[0])
	if err != nil {
		return t, tc.propagateOrWrapError(err, expr, "failed to type check float func arg: %s", err.Error())
	}

	if cotypes.GetTypeCategory(valType) != cotypes.CategoryNumeric {
		return t, fmt.Errorf("cannot convert %s to float", valType)
	}

	return cotypes.FloatType{}, nil
}

func (tc *TypeChecker) Transform(program *ast.Program) *ast.Program {
	for _, stmt := range program.Statements {
		tc.checkStatement(stmt)
	}

	return program
}

func (tc *TypeChecker) Errors() []error {
	return tc.errors
}

func (tc *TypeChecker) HasErrors() bool {
	return len(tc.errors) > 0
}
