package typechecker

import (
	"github.com/0xmukesh/coco/internal/ast"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

type builtinsInfo struct {
	name    string
	kind    ast.BuiltinsKind
	checker func(*ast.CallExpression) (t cotypes.Type, err error)
}

func (tc *TypeChecker) registerBuiltins() {
	tc.builtins["print"] = &builtinsInfo{
		name:    "print",
		kind:    ast.BuiltinFuncPrint,
		checker: tc.checkPrintBuiltin,
	}
	tc.builtins["exit"] = &builtinsInfo{
		name:    "exit",
		kind:    ast.BuiltinFuncExit,
		checker: tc.checkExitBuiltin,
	}
	tc.builtins["int"] = &builtinsInfo{
		name:    "int",
		kind:    ast.BuiltinFuncInt,
		checker: tc.checkIntBuiltin,
	}
	tc.builtins["float"] = &builtinsInfo{
		name:    "float",
		kind:    ast.BuiltinFuncFloat,
		checker: tc.checkFloatBuiltin,
	}
}
