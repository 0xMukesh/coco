package codegen

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/tokens"
	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Variable struct {
	alloca  *ir.InstAlloca
	varType cotypes.Type
}

type Codegen struct {
	module         *ir.Module
	fn             *ir.Func
	builder        *ir.Block
	ret            value.Value
	variables      map[string]*Variable
	runtimeFuncs   map[string]*ir.Func
	stringLiterals []*ir.Global
	stringIndices  map[string]int

	errors []error
}

func New() *Codegen {
	module := ir.NewModule()
	fn := module.NewFunc("main", types.I32)
	builder := fn.NewBlock("")

	cg := &Codegen{
		module:         module,
		fn:             fn,
		builder:        builder,
		variables:      make(map[string]*Variable),
		errors:         make([]error, 0),
		runtimeFuncs:   make(map[string]*ir.Func),
		stringLiterals: make([]*ir.Global, 0),
		stringIndices:  make(map[string]int),
	}

	return cg
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
		return nil, cg.addErrorAtNode(expr, "undefined variable %q", expr.Literal)
	}

	return cg.builder.NewLoad(variable.alloca.ElemType, variable.alloca), nil
}

func (cg *Codegen) generateCallExpression(expr *ast.CallExpression) (value.Value, error) {
	// TODO: only builtin functions (just print) are supported
	if !expr.IsBuiltin {
		return nil, cg.addErrorAtNode(expr, "cannot call %q identifier", expr.Identifier.String())
	}

	if expr.IsBuiltin && expr.BuiltinKind == nil {
		return nil, cg.addErrorAtNode(expr, "function %q is marked as builtin but missing builtin kind", expr.Identifier.String())
	}

	switch *expr.BuiltinKind {
	case ast.BuiltinFuncPrint:
		return cg.generatePrintExpression(expr)
	default:
		return nil, cg.addErrorAtNode(expr, "unsupported builtin function %q", expr.Identifier.String())
	}
}

func (cg *Codegen) generatePrintExpression(expr *ast.CallExpression) (value.Value, error) {
	funcName := expr.Identifier.String()
	printFunc, exists := cg.runtimeFuncs[funcName]
	if !exists {
		printFunc = cg.module.NewFunc("printf", types.I32, ir.NewParam("format", types.NewPointer(types.I8)))
		printFunc.Sig.Variadic = true
		cg.setRuntimeFunc(funcName, printFunc)
	}

	// TODO: only integers are supported by print function
	for _, arg := range expr.Arguments {
		formatStrValue := ""
		var toPrintValue value.Value = nil

		switch a := arg.(type) {
		case *ast.IntegerExpression:
			formatStrValue = "%d"
			toPrintValue = constant.NewInt(types.I64, a.Value)
		}

		if len(formatStrValue) == 0 || toPrintValue == nil {
			return nil, cg.addErrorAtNode(expr, "unsupported argument type for print expression - %T", arg.GetType())
		}

		formatStrValue += "\n\x00"
		formatStrGlobalDef := cg.getGlobalStringLiteralDef(formatStrValue)
		formatStrGlobalDef.Immutable = true
		formatStrGlobalDef.Linkage = enum.LinkagePrivate
		formatStrGlobalDef.UnnamedAddr = enum.UnnamedAddrUnnamedAddr

		formatStrPtr := constant.NewGetElementPtr(
			types.NewArray(uint64(len(formatStrValue)), types.I8),
			formatStrGlobalDef,
			constant.NewInt(types.I64, 0),
			constant.NewInt(types.I64, 0),
		)

		cg.builder.NewCall(printFunc, formatStrPtr, toPrintValue)
	}

	return nil, nil
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
	case *ast.CallExpression:
		return cg.generateCallExpression(e)
	default:
		return nil, cg.addErrorAtNode(expr, "unsupported expression type")
	}
}

func (cg *Codegen) generateLetStatement(stmt *ast.LetStatement) error {
	varName := stmt.Identifier.String()
	_, exists := cg.variables[varName]
	if exists {
		return cg.addErrorAtNode(stmt, "cannot redeclare %q variable", varName)
	}

	initValue, err := cg.generateExpression(stmt.Value)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to generate let statement value: %s", err.Error())
	}

	varType := stmt.Value.GetType()
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
