package typechecker

import (
	"github.com/0xmukesh/coco/internal/ast"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

type builtinsInfo struct {
	name    string
	kind    ast.BuiltinsKind
	checker func(*ast.CallExpression) cotypes.Type
}

func (tc *TypeChecker) registerBuiltins() {
	tc.builtins["print"] = &builtinsInfo{
		name:    "print",
		kind:    ast.BuiltinFuncPrint,
		checker: tc.checkPrintBuiltin,
	}
}
