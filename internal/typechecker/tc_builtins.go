package typechecker

import (
	"github.com/0xmukesh/coco/internal/ast"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

type BuiltinFuncInfo struct {
	name    string
	kind    ast.BuiltinsKind
	checker func(*ast.CallExpression) (t cotypes.Type, err error)
}

func (tc *TypeChecker) registerBuiltins() {
	tc.builtins["print"] = &BuiltinFuncInfo{
		name:    "print",
		kind:    ast.BuiltinFuncPrint,
		checker: tc.checkPrintBuiltin,
	}
	tc.builtins["exit"] = &BuiltinFuncInfo{
		name:    "exit",
		kind:    ast.BuiltinFuncExit,
		checker: tc.checkExitBuiltin,
	}
	tc.builtins["int"] = &BuiltinFuncInfo{
		name:    "int",
		kind:    ast.BuiltinFuncIntCoercion,
		checker: tc.checkIntCoercionBuiltin,
	}
	tc.builtins["float"] = &BuiltinFuncInfo{
		name:    "float",
		kind:    ast.BuiltinFuncFloatCoercion,
		checker: tc.checkFloatCoercionBuiltin,
	}
}

func (tc *TypeChecker) checkPrintBuiltin(expr *ast.CallExpression) (cotypes.Type, error) {
	for i, arg := range expr.Arguments {
		argType, err := tc.checkExpression(arg)
		if err != nil {
			return nil, tc.propagateOrWrapError(err, expr, "failed to type check print func arg at %d idx: %s", i, err.Error())
		}

		arg.SetType(argType)

		// TODO: at the moment, only ints, floats, bools and strings are allowed to be printed
		if !argType.Equals(cotypes.IntType{}) && !argType.Equals(cotypes.FloatType{}) && !argType.Equals(cotypes.BoolType{}) && !argType.Equals(cotypes.StringType{}) {
			return nil, tc.propagateOrWrapError(nil, expr, "invalid argument at %d idx to print", i)
		}
	}

	return cotypes.VoidType{}, nil
}

func (tc *TypeChecker) checkExitBuiltin(expr *ast.CallExpression) (cotypes.Type, error) {
	if len(expr.Arguments) != 1 {
		return nil, tc.propagateOrWrapError(nil, expr, "too many arguments. expected one argument, got %d arguments", len(expr.Arguments))
	}

	exitCode, err := tc.checkExpression(expr.Arguments[0])
	if err != nil {
		return nil, tc.propagateOrWrapError(err, expr, "failed to type check exit func arg: %s", err.Error())
	}

	if !exitCode.Equals(cotypes.IntType{}) {
		return nil, tc.propagateOrWrapError(nil, expr, "expected exit code to be of type int, got %s", exitCode.String())
	}

	return cotypes.VoidType{}, nil
}

func (tc *TypeChecker) checkIntCoercionBuiltin(expr *ast.CallExpression) (cotypes.Type, error) {
	if len(expr.Arguments) != 1 {
		return nil, tc.propagateOrWrapError(nil, expr, "too many arguments. expected one argument, got %d arguments", len(expr.Arguments))
	}

	valType, err := tc.checkExpression(expr.Arguments[0])
	if err != nil {
		return nil, tc.propagateOrWrapError(err, expr, "failed to type check int func arg: %s", err.Error())
	}

	if cotypes.GetTypeCategory(valType) != cotypes.CategoryNumeric {
		return nil, tc.propagateOrWrapError(nil, expr, "cannot convert %s to int", valType)
	}

	return cotypes.IntType{}, nil
}

func (tc *TypeChecker) checkFloatCoercionBuiltin(expr *ast.CallExpression) (cotypes.Type, error) {
	if len(expr.Arguments) != 1 {
		return nil, tc.propagateOrWrapError(nil, expr, "too many arguments. expected one argument, got %d arguments", len(expr.Arguments))
	}

	valType, err := tc.checkExpression(expr.Arguments[0])
	if err != nil {
		return nil, tc.propagateOrWrapError(err, expr, "failed to type check float func arg: %s", err.Error())
	}

	if cotypes.GetTypeCategory(valType) != cotypes.CategoryNumeric {
		return nil, tc.propagateOrWrapError(nil, expr, "cannot convert %s to float", valType)
	}

	return cotypes.FloatType{}, nil
}
