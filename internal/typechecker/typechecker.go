package typechecker

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/env"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

type TypeEnvironment = *env.Environent[cotypes.Type]

type TypeChecker struct {
	env      TypeEnvironment
	builtins map[string]*BuiltinFuncInfo

	errors []error
}

func New() *TypeChecker {
	tc := &TypeChecker{
		env:      env.NewEnvironment[cotypes.Type](),
		builtins: make(map[string]*BuiltinFuncInfo),
		errors:   []error{},
	}

	tc.registerBuiltins()
	return tc
}

func (tc *TypeChecker) Transform(program *ast.Program) *ast.Program {
	for _, stmt := range program.Statements {
		// process through all the statements and let errors collect in the background
		tc.checkStatement(stmt)
	}

	return program
}

func (tc *TypeChecker) HasErrors() bool {
	return len(tc.errors) > 0
}

func (tc *TypeChecker) Errors() []error {
	return tc.errors
}
