package codegen

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	alloca  *ir.InstAlloca
	varType cotypes.Type
}

type Codegen struct {
	module    *ir.Module
	fn        *ir.Func
	builder   *ir.Block
	ret       value.Value
	variables map[string]*Variable
	errors    []error
}

func New() *Codegen {
	module := ir.NewModule()
	fn := module.NewFunc("main", types.I32)
	builder := fn.NewBlock("")

	return &Codegen{
		module:    module,
		fn:        fn,
		builder:   builder,
		variables: make(map[string]*Variable),
		errors:    make([]error, 0),
	}
}

func (cg *Codegen) addErrorAtNode(node ast.Node, msg string, args ...any) error {
	err := &CodegenError{
		message: fmt.Sprintf(msg, args...),
		node:    node,
	}

	cg.errors = append(cg.errors, err)
	return err
}

func (cg *Codegen) addError(msg string, args ...any) error {
	err := fmt.Errorf(msg, args...)
	cg.errors = append(cg.errors, err)
	return err
}

func (cg *Codegen) propagateOrWrapError(err error, node ast.Node, msg string, args ...any) error {
	if isCodegenError(err) {
		return err
	}

	return cg.addErrorAtNode(node, msg, args...)
}

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

func (cg *Codegen) llvmToType(t types.Type) (cotypes.Type, error) {
	switch t {
	case types.I64:
		return cotypes.IntType{}, nil
	case types.Double:
		return cotypes.FloatType{}, nil
	case types.I1:
		return cotypes.BoolType{}, nil
	default:
		return nil, cg.addError("unsupported LLVM type - %T", t)
	}
}

func (cg *Codegen) generateBinaryExpression(expr *ast.BinaryExpression) (value.Value, error) {
	left, err := cg.generateExpression(expr.Left)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to generate left operand: %s", err.Error())
	}

	right, err := cg.generateExpression(expr.Right)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to generate right operand: %s", err.Error())
	}

	isFloat := expr.GetType().Equals(cotypes.FloatType{})

	switch expr.Operator.Type {
	case tokens.PLUS:
		if isFloat {
			return cg.builder.NewFAdd(left, right), nil
		}
		return cg.builder.NewAdd(left, right), nil
	case tokens.MINUS:
		if isFloat {
			return cg.builder.NewFSub(left, right), nil
		}
		return cg.builder.NewSub(left, right), nil
	case tokens.STAR:
		if isFloat {
			return cg.builder.NewFMul(left, right), nil
		}
		return cg.builder.NewMul(left, right), nil
	case tokens.SLASH:
		if isFloat {
			return cg.builder.NewFDiv(left, right), nil
		}
		return cg.builder.NewSDiv(left, right), nil
	default:
		return nil, cg.addErrorAtNode(expr, "unknown binary operator: %s", expr.Operator.Type)
	}
}

func (cg *Codegen) generateIdentifier(expr *ast.IdentifierExpression) (value.Value, error) {
	variable, exists := cg.variables[expr.Literal]
	if !exists {
		return nil, cg.addErrorAtNode(expr, "undefined variable '%s'", expr.Literal)
	}

	return cg.builder.NewLoad(variable.alloca.ElemType, variable.alloca), nil
}

func (cg *Codegen) generateExpression(expr ast.Expression) (value.Value, error) {
	if expr.GetType() == nil {
		return nil, cg.addErrorAtNode(expr, "expression has no type")
	}

	switch e := expr.(type) {
	case *ast.IntegerExpression:
		return constant.NewInt(types.I64, e.Value), nil
	case *ast.FloatExpression:
		return constant.NewFloat(types.Double, e.Value), nil
	case *ast.IdentifierExpression:
		return cg.generateIdentifier(e)
	case *ast.BinaryExpression:
		return cg.generateBinaryExpression(e)
	default:
		return nil, cg.addErrorAtNode(expr, "unsupported expression type")
	}
}

func (cg *Codegen) generateLetStatement(stmt *ast.LetStatement) error {
	initValue, err := cg.generateExpression(stmt.Value)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to generate let statement value: %s", err.Error())
	}

	varType := stmt.Value.GetType()
	varName := stmt.Identifier.String()
	llvmType, err := cg.typeToLlvm(varType)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to retrieve llvm equivalent type: %s", err.Error())
	}

	alloca := cg.builder.NewAlloca(llvmType)
	alloca.SetName(varName)

	cg.builder.NewStore(initValue, alloca)

	cg.variables[varName] = &Variable{
		alloca:  alloca,
		varType: varType,
	}
	return nil
}

func (cg *Codegen) generateAssignmentStatement(stmt *ast.AssignmentStatement) error {
	varName := stmt.Identifier.String()
	variable, exists := cg.variables[varName]
	if !exists {
		return cg.addErrorAtNode(stmt, "cannot assign to undefined variable: %s", varName)
	}

	newValue, err := cg.generateExpression(stmt.Value)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to generate assignment statement value: %s", err.Error())
	}
	newType := stmt.Value.GetType()

	if !variable.varType.Equals(newType) {
		return cg.addErrorAtNode(stmt, "cannot assign %s type to variable of type %s", newType, variable.varType)
	}

	cg.builder.NewStore(newValue, variable.alloca)
	return nil
}

func (cg *Codegen) generateExitStatement(stmt *ast.ExitStatement) error {
	exitVal, err := cg.generateExpression(stmt.Expr)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to generate exit value: %s", err.Error())
	}

	if exitVal.Type() == types.I64 {
		cg.ret = cg.builder.NewTrunc(exitVal, types.I32)
	} else {
		cg.ret = exitVal
	}

	return nil
}

func (cg *Codegen) generateStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		_, err := cg.generateExpression(s.Expr)
		if err != nil {
			return err
		}
	case *ast.LetStatement:
		return cg.generateLetStatement(s)
	case *ast.AssignmentStatement:
		return cg.generateAssignmentStatement(s)
	case *ast.ExitStatement:
		return cg.generateExitStatement(s)
	}

	return nil
}

func (cg *Codegen) Generate(program *ast.Program) *ir.Module {
	for _, stmt := range program.Statements {
		if err := cg.generateStatement(stmt); err != nil {
			continue
		}
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

func (cg *Codegen) HasErrors() bool {
	return len(cg.errors) > 0
}

func (cg *Codegen) Errors() []error {
	return cg.errors
}
