package codegen

import (
	"fmt"
	"strconv"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *Codegen) generateExpression(expr ast.Expression) (val value.Value, err error) {
	if expr.GetType() == nil {
		return nil, cg.propagateOrWrapError(nil, expr, "expression has no type")
	}

	switch e := expr.(type) {
	case *ast.IntegerExpression:
		val = constant.NewInt(types.I64, e.Value)
	case *ast.FloatExpression:
		val = constant.NewFloat(types.Double, e.Value)
	case *ast.BooleanExpression:
		val = constant.NewBool(e.Value)
	case *ast.StringExpression:
		return cg.generateStringExpression(e)
	case *ast.IdentifierExpression:
		return cg.generateIdentifierExpression(e)
	case *ast.BinaryExpression:
		return cg.generateBinaryExpression(e)
	case *ast.CallExpression:
		return cg.generateCallExpression(e)
	case *ast.IfExpression:
		return cg.generateIfExpression(e)
	case *ast.GroupedExpression:
		return cg.generateExpression(e.Expr)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", e)
	}

	if coercible, ok := expr.(ast.CoercibleExpression); ok {
		if targetType := coercible.GetCoercion(); targetType != nil {
			val, err = cg.applyCoercion(val, expr.GetType(), targetType)
			if err != nil {
				return nil, cg.propagateOrWrapError(err, expr, "failed to apply type coercion: %s", err.Error())
			}
		}
	}

	return val, nil
}

func (cg *Codegen) generateStringExpression(expr *ast.StringExpression) (value.Value, error) {
	str, err := strconv.Unquote(expr.Value)
	if err != nil {
		return nil, cg.propagateOrWrapError(nil, expr, "invalid string: %s", err.Error())
	}

	return cg.getOrCreateStringLiteral(str, ""), nil
}

func (cg *Codegen) generateIdentifierExpression(expr *ast.IdentifierExpression) (value.Value, error) {
	variable, exists := cg.scope.Get(expr.Literal)
	if !exists {
		return nil, cg.propagateOrWrapError(nil, expr, "unknown variable %q", expr.Literal)
	}

	return cg.builder.NewLoad(variable.alloca.ElemType, variable.alloca), nil
}

func (cg *Codegen) generateBinaryExpression(expr *ast.BinaryExpression) (value.Value, error) {
	left, err := cg.generateExpression(expr.Left)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr.Left, "failed to codegen left operand: %s", err.Error())
	}

	right, err := cg.generateExpression(expr.Right)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr.Right, "failed to codegen right operand: %s", err.Error())
	}

	// boolean comparision
	if expr.GetType().Equals(cotypes.BoolType{}) {
		switch expr.Operator.Type {
		case tokens.EQUALS:
			return cg.builder.NewICmp(enum.IPredEQ, left, right), nil
		case tokens.NOT_EQUALS:
			return cg.builder.NewICmp(enum.IPredNE, left, right), nil
		}
	}

	// integer arithemtic
	if expr.GetType().Equals(cotypes.IntType{}) {
		switch expr.Operator.Type {
		case tokens.PLUS:
			return cg.builder.NewAdd(left, right), nil
		case tokens.MINUS:
			return cg.builder.NewSub(left, right), nil
		case tokens.STAR:
			return cg.builder.NewMul(left, right), nil
		case tokens.SLASH:
			return cg.builder.NewSDiv(left, right), nil
		default:
			return nil, cg.propagateOrWrapError(nil, expr, "cannot perform %q operation on %q and %q types", expr.Operator, left.Type(), right.Type())
		}
	}

	// float arithemtic
	if expr.GetType().Equals(cotypes.FloatType{}) {
		switch expr.Operator.Type {
		case tokens.PLUS:
			return cg.builder.NewFAdd(left, right), nil
		case tokens.MINUS:
			return cg.builder.NewFSub(left, right), nil
		case tokens.STAR:
			return cg.builder.NewFMul(left, right), nil
		case tokens.SLASH:
			return cg.builder.NewFDiv(left, right), nil
		default:
			return nil, cg.propagateOrWrapError(nil, expr, "cannot perform %q operation on %q and %q types", expr.Operator.Type, left.Type(), right.Type())
		}
	}

	// integer comparison
	if left.Type().Equal(types.I64) && right.Type().Equal(types.I64) && expr.GetType().Equals(cotypes.BoolType{}) {
		switch expr.Operator.Type {
		case tokens.LESS_THAN:
			return cg.builder.NewICmp(enum.IPredSLT, left, right), nil
		case tokens.GREATER_THAN:
			return cg.builder.NewICmp(enum.IPredSGT, left, right), nil
		case tokens.LESS_THAN_EQUALS:
			return cg.builder.NewICmp(enum.IPredSLE, left, right), nil
		case tokens.GREATER_THAN_EQUALS:
			return cg.builder.NewICmp(enum.IPredSGE, left, right), nil
		case tokens.EQUALS:
			return cg.builder.NewICmp(enum.IPredEQ, left, right), nil
		case tokens.NOT_EQUALS:
			return cg.builder.NewICmp(enum.IPredNE, left, right), nil
		default:
			return nil, cg.propagateOrWrapError(nil, expr, "cannot perform %q operation on %q and %q types", expr.Operator.Type, left.Type(), right.Type())
		}
	}

	// float comparison
	if left.Type().Equal(types.Double) && right.Type().Equal(types.Double) && expr.GetType().Equals(cotypes.BoolType{}) {
		switch expr.Operator.Type {
		case tokens.LESS_THAN:
			return cg.builder.NewFCmp(enum.FPredOLT, left, right), nil
		case tokens.GREATER_THAN:
			return cg.builder.NewFCmp(enum.FPredOGT, left, right), nil
		case tokens.LESS_THAN_EQUALS:
			return cg.builder.NewFCmp(enum.FPredOLE, left, right), nil
		case tokens.GREATER_THAN_EQUALS:
			return cg.builder.NewFCmp(enum.FPredOGE, left, right), nil
		case tokens.EQUALS:
			return cg.builder.NewFCmp(enum.FPredOEQ, left, right), nil
		case tokens.NOT_EQUALS:
			return cg.builder.NewFCmp(enum.FPredONE, left, right), nil
		default:
			return nil, cg.propagateOrWrapError(nil, expr, "cannot perform %q operation on %q and %q types", expr.Operator.Type, left.Type(), right.Type())
		}
	}

	return nil, cg.propagateOrWrapError(nil, expr, "cannot perform %q operation on %q and %q types", expr.Operator.Type, left.Type(), right.Type())
}

func (cg *Codegen) generateCallExpression(expr *ast.CallExpression) (value.Value, error) {
	// TODO: at the moment, only builtin functions can be called
	if !expr.IsBuiltin {
		return nil, cg.propagateOrWrapError(nil, expr, "cannot call %q identifier", expr.Identifier)
	}

	if expr.IsBuiltin && expr.BuiltinKind == nil {
		return nil, cg.propagateOrWrapError(nil, expr, "function %q is marked as builtin but missing builtin kind", expr.Identifier)
	}

	switch *expr.BuiltinKind {
	case ast.BuiltinFuncPrint:
		return cg.generatePrintExpression(expr)
	case ast.BuiltinFuncExit:
		return cg.generateExitExpression(expr)
	case ast.BuiltinFuncIntCoercion:
		return cg.generateIntCoercionExpression(expr)
	case ast.BuiltinFuncFloatCoercion:
		return cg.generateFloatCoercionExpression(expr)
	default:
		return nil, cg.propagateOrWrapError(nil, expr, "unsupported builtin function: %s", expr.Identifier)
	}
}

func (cg *Codegen) generateIfExpression(expr *ast.IfExpression) (value.Value, error) {
	condition, err := cg.generateExpression(expr.Condition)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr.Condition, "failed to codegen if-branch condition: %s", err.Error())
	}

	trueBlock := cg.mainFn.NewBlock("")
	mergeBlock := cg.mainFn.NewBlock("")

	var falseBlock *ir.Block
	if expr.Alternative != nil {
		falseBlock = cg.mainFn.NewBlock("")
		cg.builder.NewCondBr(condition, trueBlock, falseBlock)
	} else {
		cg.builder.NewCondBr(condition, trueBlock, mergeBlock)
	}

	cg.builder = trueBlock
	if err := cg.generateStatement(expr.Consequence); err != nil {
		return nil, cg.propagateOrWrapError(err, expr.Consequence, "failed to codegen if-branch: %s", err.Error())
	}
	trueReturnVal := cg.blockReturnValue
	trueTerminated := cg.builder.Term != nil
	if !trueTerminated {
		cg.builder.NewBr(mergeBlock)
	}

	var falseReturnVal value.Value
	falseTerminated := false
	if expr.Alternative != nil {
		cg.builder = falseBlock
		if err := cg.generateStatement(expr.Alternative); err != nil {
			return nil, cg.propagateOrWrapError(err, expr.Alternative, "failed to codegen else-branch: %s", err.Error())
		}

		falseReturnVal = cg.blockReturnValue
		falseTerminated = cg.builder.Term != nil
		if !falseTerminated {
			cg.builder.NewBr(mergeBlock)
		}
	}

	allTerminated := expr.Alternative != nil && trueTerminated && falseTerminated
	if allTerminated {
		mergeBlock.NewUnreachable()
	}

	cg.builder = mergeBlock
	if expr.Alternative != nil && trueReturnVal != nil && falseReturnVal != nil && !allTerminated {
		cg.blockReturnValue = mergeBlock.NewPhi(
			ir.NewIncoming(trueReturnVal, trueBlock),
			ir.NewIncoming(falseReturnVal, falseBlock),
		)
	} else {
		cg.blockReturnValue = nil
	}

	return cg.blockReturnValue, nil
}
