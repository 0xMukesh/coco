package codegen

import (
	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

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

func (cg *Codegen) getRuntimePrintFuncByType(t cotypes.Type) string {
	switch t.(type) {
	case cotypes.IntType:
		return "__coco_print_int"
	case cotypes.FloatType:
		return "__coco_print_float"
	case cotypes.BoolType:
		return "__coco_print_bool"
	default:
		return ""
	}
}

func (cg *Codegen) setRuntimeFunc(name string, llvmFunc *ir.Func) {
	cg.runtimeFuncs[name] = llvmFunc
}
