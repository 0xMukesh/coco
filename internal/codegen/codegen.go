package codegen

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Codegen struct {
	module  *ir.Module
	fn      *ir.Func
	builder *ir.Block
	ret     value.Value
}

func New() *Codegen {
	module := ir.NewModule()
	fn := module.NewFunc("main", types.I32)
	builder := fn.NewBlock("")

	return &Codegen{
		module:  module,
		fn:      fn,
		builder: builder,
	}
}

func (cg *Codegen) generateBinaryExpression(expr *ast.BinaryExpression) value.Value {
	left := cg.generateExpression(expr.Left)
	right := cg.generateExpression(expr.Right)

	// TODO: only int types are handled
	switch expr.Operator.Type {
	case tokens.PLUS:
		return cg.builder.NewAdd(left, right)
	case tokens.MINUS:
		return cg.builder.NewAdd(left, right)
	case tokens.STAR:
		return cg.builder.NewAdd(left, right)
	case tokens.SLASH:
		return cg.builder.NewAdd(left, right)
	default:
		return nil
	}
}

func (cg *Codegen) generateExpression(expr ast.Expression) value.Value {
	if expr.GetType() == nil {
		panic(fmt.Sprintf("expression has no type: %s", expr))
	}

	switch e := expr.(type) {
	case *ast.IntegerExpression:
		return constant.NewInt(types.I64, e.Value)
	case *ast.BinaryExpression:
		return cg.generateBinaryExpression(e)
	default:
		return nil
	}
}

func (cg *Codegen) generateStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		cg.generateExpression(s.Expr)
	case *ast.ExitStatement:
		exitVal := cg.generateExpression(s.Expr)

		if exitVal.Type() == types.I64 {
			cg.ret = cg.builder.NewTrunc(exitVal, types.I32)
		} else {
			cg.ret = exitVal
		}
	}
}

func (cg *Codegen) Generate(program *ast.Program) *ir.Module {
	for _, stmt := range program.Statements {
		cg.generateStatement(stmt)
	}

	if cg.ret != nil {
		cg.builder.NewRet(cg.ret)
	} else {
		cg.builder.NewRet(constant.NewInt(types.I32, 0))
	}

	return cg.module
}

func (cg *Codegen) EmitIR() string {
	return cg.module.String()
}
