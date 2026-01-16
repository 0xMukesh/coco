package codegen

import (
	"fmt"

	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
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

func (cg *Codegen) setRuntimeFunc(name string, llvmFunc *ir.Func) {
	cg.runtimeFuncs[name] = llvmFunc
}

func (cg *Codegen) getGlobalStringLiteralDef(literal string) *ir.Global {
	if idx, exists := cg.stringIndices[literal]; exists {
		return cg.stringLiterals[idx]
	}

	name := fmt.Sprintf(".str.%d", len(cg.stringLiterals))
	globalDef := cg.module.NewGlobalDef(name, constant.NewCharArrayFromString(literal))

	cg.stringIndices[literal] = len(cg.stringLiterals)
	cg.stringLiterals = append(cg.stringLiterals, globalDef)

	return globalDef
}
