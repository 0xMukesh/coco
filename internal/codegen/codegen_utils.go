package codegen

import (
	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

var TRUE_GLOBAL_DEF_NAME = "__coco_true"
var FALSE_GLOBAL_DEF_NAME = "__coco_false"

func (cg *Codegen) typeToLlvm(t cotypes.Type) (types.Type, error) {
	switch t.(type) {
	case cotypes.IntType:
		return types.I64, nil
	case cotypes.FloatType:
		return types.Double, nil
	case cotypes.BoolType:
		return types.I1, nil
	default:
		return nil, cg.addError("unsupported type - %v", t)
	}
}

func (cg *Codegen) setupPrintfRuntimeFunc() *ir.Func {
	printfFunc := cg.module.NewFunc("printf", types.I32, ir.NewParam("fmt", types.NewPointer(types.I8)))
	printfFunc.Sig.Variadic = true
	cg.runtimeFuncs["print"] = printfFunc

	return printfFunc
}

func (cg *Codegen) setupTrueGlobalDef() *ir.Global {
	trueStr := cg.module.NewGlobalDef(TRUE_GLOBAL_DEF_NAME, constant.NewCharArrayFromString("true\x00"))
	trueStr.Immutable = true
	trueStr.Linkage = enum.LinkagePrivate
	cg.globalDefs[TRUE_GLOBAL_DEF_NAME] = trueStr

	return trueStr
}

func (cg *Codegen) setupFalseStrGlobalDef() *ir.Global {
	falseStr := cg.module.NewGlobalDef(FALSE_GLOBAL_DEF_NAME, constant.NewCharArrayFromString("false\x00"))
	falseStr.Immutable = true
	falseStr.Linkage = enum.LinkagePrivate
	cg.globalDefs[FALSE_GLOBAL_DEF_NAME] = falseStr

	return falseStr
}
