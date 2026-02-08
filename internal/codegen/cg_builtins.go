package codegen

import (
	"strings"

	"github.com/0xmukesh/coco/internal/ast"
	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *Codegen) generatePrintExpression(expr *ast.CallExpression) (value.Value, error) {
	funcName := expr.Identifier.String()
	printFunc, ok := cg.runtimeFuncs[funcName]
	if !ok {
		printFunc = cg.setupPrintFunction()
	}

	var fmtStr strings.Builder
	var boolArgs []value.Value // slice containing the pre-processed string literal version of booleans i.e. "true" for 1 and "false" for 0

	for i, arg := range expr.Arguments {
		if i > 0 {
			fmtStr.WriteString(" ")
		}

		switch arg.GetType().(type) {
		case cotypes.IntType:
			fmtStr.WriteString("%ld")
		case cotypes.FloatType:
			fmtStr.WriteString("%g")
		case cotypes.StringType:
			fmtStr.WriteString("%s")
		case cotypes.BoolType:
			fmtStr.WriteString("%s")

			trueStr, ok := cg.globalDefs[TRUE_STR_GLOBAL_DEF_NAME]
			if !ok {
				trueStr = cg.getOrCreateStringLiteral("true", TRUE_STR_GLOBAL_DEF_NAME)
			}
			falseStr, ok := cg.globalDefs[FALSE_STR_GLOBAL_DEF_NAME]
			if !ok {
				falseStr = cg.getOrCreateStringLiteral("false", FALSE_STR_GLOBAL_DEF_NAME)
			}

			truePtr, err := cg.getStringLiteralPointer(trueStr)
			if err != nil {
				return nil, cg.propagateOrWrapError(err, expr, "failed to codegen string literal pointer: true")
			}
			falsePtr, err := cg.getStringLiteralPointer(falseStr)
			if err != nil {
				return nil, cg.propagateOrWrapError(err, expr, "failed to codegen string literal pointer: false")
			}

			boolValue, err := cg.generateExpression(arg)
			if err != nil {
				return nil, cg.propagateOrWrapError(err, expr, "failed to codegen arg: %q", arg)
			}

			strPtr := cg.builder.NewSelect(boolValue, truePtr, falsePtr)
			boolArgs = append(boolArgs, strPtr)
		}
	}

	fmtStr.WriteString("\n")

	fmtStrGlobalDef := cg.getOrCreateStringLiteral(fmtStr.String(), "")
	fmtPtr, err := cg.getStringLiteralPointer(fmtStrGlobalDef)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to codegen string literal pointer: %s", fmtStr.String())
	}

	args := []value.Value{fmtPtr}
	boolIdx := 0

	for _, arg := range expr.Arguments {
		switch arg.GetType().(type) {
		case cotypes.BoolType:
			args = append(args, boolArgs[boolIdx])
			boolIdx++
		default:
			v, err := cg.generateExpression(arg)
			if err != nil {
				return nil, cg.propagateOrWrapError(err, expr, "failed to codegen arg: %s", arg)
			}

			// if it is string literal then pass then pointer of it
			if arg.GetType().Equals(cotypes.StringType{}) {
				v, err = cg.getStringLiteralPointer(v.(*ir.Global))
				if err != nil {
					return nil, cg.propagateOrWrapError(err, expr, "failed to codegen string literal pointer: %s", arg.String())
				}
			}

			args = append(args, v)
		}
	}

	cg.builder.NewCall(printFunc, args...)
	return nil, nil
}

func (cg *Codegen) generateExitExpression(expr *ast.CallExpression) (value.Value, error) {
	exitCodeArg := expr.Arguments[0]
	exitVal, err := cg.generateExpression(exitCodeArg)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to codegen arg: %s", exitCodeArg)
	}

	cg.builder.NewRet(exitVal)
	return nil, nil
}

func (cg *Codegen) generateIntCoercionExpression(expr *ast.CallExpression) (value.Value, error) {
	intArg := expr.Arguments[0]
	val, err := cg.generateExpression(intArg)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to codegen arg: %s", intArg)
	}

	// if value isn't int64 then only run the FPToSI (floating-point to signed integer) instruction
	if !val.Type().Equal(types.I64) {
		return cg.builder.NewFPToSI(val, types.I64), nil
	}

	return val, nil
}

func (cg *Codegen) generateFloatCoercionExpression(expr *ast.CallExpression) (value.Value, error) {
	floatArg := expr.Arguments[0]
	val, err := cg.generateExpression(floatArg)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to codegen arg: %s", floatArg)
	}

	// if value isn't double then only run the SIToFP (signed integer to floating-point) instruction
	if !val.Type().Equal(types.Double) {
		return cg.builder.NewSIToFP(val, types.Double), nil
	}

	return val, nil
}
