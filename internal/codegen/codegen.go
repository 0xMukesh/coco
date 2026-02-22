package codegen

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/env"
	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Scope = *env.Environent[ScopeItem]
type ScopeItem struct {
	alloca *ir.InstAlloca
	typ    cotypes.Type
}

type Codegen struct {
	module *ir.Module
	mainFn *ir.Func

	builder *ir.Block
	// used to transfer the return value among statements and expressions, mainly
	blockReturnValue value.Value
	// used for codegen of break statements
	loopExitBlock *ir.Block

	scope        Scope
	runtimeFuncs map[string]*ir.Func
	globalDefs   map[string]*ir.Global

	nameCounter int
	errors      []error
}

func New() *Codegen {
	module := ir.NewModule()
	mainFn := module.NewFunc("main", types.I64)
	builder := mainFn.NewBlock("")

	cg := &Codegen{
		module:       module,
		mainFn:       mainFn,
		builder:      builder,
		scope:        env.NewEnvironment[ScopeItem](),
		runtimeFuncs: make(map[string]*ir.Func),
		globalDefs:   make(map[string]*ir.Global),
		errors:       make([]error, 0),
	}

	return cg
}

func (cg *Codegen) Generate(program *ast.Program) *ir.Module {
	for _, stmt := range program.Statements {
		// don't check for errors and process through the entire ast to collect all the errors
		cg.generateStatement(stmt)
	}

	if cg.HasErrors() {
		return nil
	}

	if cg.builder.Term == nil {
		cg.builder.NewRet(constant.NewInt(types.I64, 0))
	}

	return cg.module
}

func (cg *Codegen) EmitIR() string {
	return cg.module.String()
}

func (cg *Codegen) HasErrors() bool {
	return len(cg.errors) > 0
}

func (cg *Codegen) Errors() []error {
	return cg.errors
}
