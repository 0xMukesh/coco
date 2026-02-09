package codegen

import (
	"fmt"
	"maps"
	"slices"

	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/0xmukesh/coco/internal/utils"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	TRUE_STR_GLOBAL_DEF_NAME  = "__coco_true"
	FALSE_STR_GLOBAL_DEF_NAME = "__coco_false"
)

func (cg *Codegen) applyCoercion(val value.Value, fromType, targetType cotypes.Type) (value.Value, error) {
	if fromType.Equals(cotypes.IntType{}) && targetType.Equals(cotypes.FloatType{}) {
		return cg.builder.NewSIToFP(val, types.Double), nil
	}

	if fromType.Equals(cotypes.FloatType{}) && targetType.Equals(cotypes.IntType{}) {
		return cg.builder.NewFPToSI(val, types.I64), nil
	}

	return val, fmt.Errorf("unsupported type coercion from %s to %s", fromType, targetType)
}

func (cg *Codegen) getOrCreateStringLiteral(str string, name string) *ir.Global {
	// before passing string literals to llvm, null terminator is added to the end of the string
	str += "\x00"

	globalDefs := slices.Collect(maps.Values(cg.globalDefs))
	strIdx := slices.IndexFunc(globalDefs, func(e *ir.Global) bool {
		return e.Init.Ident() == str
	})

	if strIdx != -1 {
		return globalDefs[strIdx]
	} else {
		if name == "" {
			name = fmt.Sprintf(".str.%d", cg.nameCounter)

		}

		globalDef := cg.module.NewGlobalDef(name, constant.NewCharArrayFromString(str))
		globalDef.Immutable = true
		globalDef.Linkage = enum.LinkagePrivate
		globalDef.UnnamedAddr = enum.UnnamedAddrUnnamedAddr

		cg.globalDefs[name] = globalDef
		cg.nameCounter++

		return globalDef
	}
}

func (cg *Codegen) getStringLiteralPointer(strGlobalDef *ir.Global) (*ir.InstGetElementPtr, error) {
	strValue, err := utils.LlvmStringToGoLiteral(strGlobalDef.Init.Ident())
	if err != nil {
		return nil, err
	}

	ptr := cg.builder.NewGetElementPtr(
		types.NewArray(uint64(len(strValue)), types.I8),
		strGlobalDef,
		constant.NewInt(types.I64, 0),
		constant.NewInt(types.I64, 0),
	)

	return ptr, nil
}

func (cg *Codegen) setupPrintFunction() *ir.Func {
	printFunc := cg.module.NewFunc("printf", types.I32, ir.NewParam("fmt", types.NewPointer(types.I8)))
	printFunc.Sig.Variadic = true
	cg.runtimeFuncs["print"] = printFunc

	return printFunc
}
