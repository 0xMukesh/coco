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
	module       *ir.Module
	fn           *ir.Func
	builder      *ir.Block
	ret          value.Value
	variables    map[string]*Variable
	runtimeFuncs map[string]*ir.Func

	errors []error
}

func New() *Codegen {
	module := ir.NewModule()
	fn := module.NewFunc("main", types.I32)
	builder := fn.NewBlock("")

	cg := &Codegen{
		module:       module,
		fn:           fn,
		builder:      builder,
		variables:    make(map[string]*Variable),
		errors:       make([]error, 0),
		runtimeFuncs: make(map[string]*ir.Func),
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

	// integer arithmetic
	if expr.GetType().Equals(cotypes.IntType{}) {
		switch expr.Operator.Type {
		case tokens.PLUS:
			return cg.builder.NewAdd(left, right), nil
		case tokens.MINUS:
			return cg.builder.NewSub(left, right), nil
		case tokens.STAR:
			return cg.builder.NewMul(left, right), nil
		case tokens.SLASH:
			return cg.builder.NewSDiv(left, right), nil
		default:
			return nil, cg.addErrorAtNode(expr, "cannot perform %s operation", expr.Operator.Type)
		}
	}

	// float arithmetic
	if expr.GetType().Equals(cotypes.FloatType{}) {
		switch expr.Operator.Type {
		case tokens.PLUS:
			return cg.builder.NewFAdd(left, right), nil
		case tokens.MINUS:
			return cg.builder.NewFSub(left, right), nil
		case tokens.STAR:
			return cg.builder.NewMul(left, right), nil
		case tokens.SLASH:
			return cg.builder.NewFDiv(left, right), nil
		default:
			return nil, cg.addErrorAtNode(expr, "cannot perform %s operation", expr.Operator.Type)
		}
	}

	// integer comparison
	if left.Type().Equal(types.I64) && right.Type().Equal(types.I64) && expr.GetType().Equals(cotypes.BoolType{}) {
		switch expr.Operator.Type {
		case tokens.LESS_THAN:
			return cg.builder.NewICmp(enum.IPredSLT, left, right), nil
		case tokens.GREATER_THAN:
			return cg.builder.NewICmp(enum.IPredSGT, left, right), nil
		case tokens.LESS_THAN_EQUALS:
			return cg.builder.NewICmp(enum.IPredSLE, left, right), nil
		case tokens.GREATER_THAN_EQUALS:
			return cg.builder.NewICmp(enum.IPredSGE, left, right), nil
		case tokens.EQUALS:
			return cg.builder.NewICmp(enum.IPredEQ, left, right), nil
		case tokens.NOT_EQUALS:
			return cg.builder.NewICmp(enum.IPredNE, left, right), nil
		default:
			return nil, cg.addErrorAtNode(expr, "cannot perform %s operation", expr.Operator.Type)
		}
	}

	// float comparison
	if left.Type().Equal(types.Double) && right.Type().Equal(types.Double) && expr.GetType().Equals(cotypes.BoolType{}) {
		switch expr.Operator.Type {
		case tokens.LESS_THAN:
			return cg.builder.NewFCmp(enum.FPredOLT, left, right), nil
		case tokens.GREATER_THAN:
			return cg.builder.NewFCmp(enum.FPredOGT, left, right), nil
		case tokens.LESS_THAN_EQUALS:
			return cg.builder.NewFCmp(enum.FPredOLE, left, right), nil
		case tokens.GREATER_THAN_EQUALS:
			return cg.builder.NewFCmp(enum.FPredOGE, left, right), nil
		case tokens.EQUALS:
			return cg.builder.NewFCmp(enum.FPredOEQ, left, right), nil
		case tokens.NOT_EQUALS:
			return cg.builder.NewFCmp(enum.FPredONE, left, right), nil
		default:
			return nil, cg.addErrorAtNode(expr, "cannot perform %s operation", expr.Operator.Type)
		}
	}

	return nil, cg.addErrorAtNode(expr, "cannot perform %s operation", expr.Operator.Type)
}

func (cg *Codegen) generateIdentifier(expr *ast.IdentifierExpression) (value.Value, error) {
	variable, exists := cg.variables[expr.Literal]
	if !exists {
		return nil, cg.addErrorAtNode(expr, "undefined variable %q", expr.Literal)
	}

	return cg.builder.NewLoad(variable.alloca.ElemType, variable.alloca), nil
}

func (cg *Codegen) generateCallExpression(expr *ast.CallExpression) (value.Value, error) {
	if !expr.IsBuiltin {
		return nil, cg.addErrorAtNode(expr, "cannot call %q identifier", expr.Identifier.String())
	}

	if expr.IsBuiltin && expr.BuiltinKind == nil {
		return nil, cg.addErrorAtNode(expr, "function %q is marked as builtin but missing builtin kind", expr.Identifier.String())
	}

	switch *expr.BuiltinKind {
	case ast.BuiltinFuncPrint:
		return cg.generatePrintExpression(expr)
	case ast.BuiltinFuncExit:
		return cg.generateExitExpression(expr)
	case ast.BuiltinFuncInt:
		return cg.generateIntExpression(expr)
	case ast.BuiltinFuncFloat:
		return cg.generateFloatExpression(expr)
	default:
		return nil, cg.addErrorAtNode(expr, "unsupported builtin function %q", expr.Identifier.String())
	}
}

func (cg *Codegen) generatePrintExpression(expr *ast.CallExpression) (value.Value, error) {
	for _, arg := range expr.Arguments {
		printFuncName := cg.getRuntimePrintFuncByType(arg.GetType())
		if printFuncName == "" {
			return nil, cg.addErrorAtNode(expr, "unsupported argument type for print expression - %T", arg.GetType())
		}

		printFunc, exists := cg.runtimeFuncs[printFuncName]
		if !exists {
			var param *ir.Param

			switch arg.GetType().(type) {
			case cotypes.IntType:
				param = ir.NewParam("value", types.I64)
			case cotypes.FloatType:
				param = ir.NewParam("value", types.Double)
			case cotypes.BoolType:
				param = ir.NewParam("value", types.I1)
			}

			if param == nil {
				return nil, cg.addErrorAtNode(expr, "unsupported argument type for print expression - %T", arg.GetType())
			}

			printFunc = cg.module.NewFunc(printFuncName, types.Void, param)
			cg.setRuntimeFunc(printFuncName, printFunc)
		}

		if printFunc == nil {
			return nil, cg.addErrorAtNode(expr, "unsupported argument type for print expression - %T", arg.GetType())
		}

		toPrintValue, err := cg.generateExpression(arg)
		if err != nil {
			return nil, cg.propagateOrWrapError(err, expr, "failed to generate value for print call expression argument - %T", arg.GetType())
		}

		if arg.GetType().Equals(cotypes.BoolType{}) {
			toPrintValue = cg.builder.NewZExt(toPrintValue, types.I64)
		}

		cg.builder.NewCall(printFunc, toPrintValue)
	}

	return nil, nil
}

func (cg *Codegen) generateExitExpression(expr *ast.CallExpression) (value.Value, error) {
	exitVal, err := cg.generateExpression(expr.Arguments[0])
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to generate value for exit call expression argument: %s", err.Error())
	}

	if exitVal.Type() == types.I64 {
		cg.ret = cg.builder.NewTrunc(exitVal, types.I32)
	} else {
		cg.ret = exitVal
	}

	return nil, nil
}

func (cg *Codegen) generateIntExpression(expr *ast.CallExpression) (value.Value, error) {
	val, err := cg.generateExpression(expr.Arguments[0])
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to generate value for int call expression argument: %s", err.Error())
	}

	if !val.Type().Equal(types.I64) {
		return cg.builder.NewFPToSI(val, types.I64), nil
	}

	return val, nil
}

func (cg *Codegen) generateFloatExpression(expr *ast.CallExpression) (value.Value, error) {
	val, err := cg.generateExpression(expr.Arguments[0])
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to generate value for float call expression argument: %s", err.Error())
	}

	if !val.Type().Equal(types.Double) {
		return cg.builder.NewSIToFP(val, types.Double), nil
	}

	return val, nil
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
	case *ast.BooleanExpression:
		return constant.NewBool(e.Value), nil
	case *ast.IdentifierExpression:
		return cg.generateIdentifier(e)
	case *ast.BinaryExpression:
		return cg.generateBinaryExpression(e)
	case *ast.CallExpression:
		return cg.generateCallExpression(e)
	case *ast.GroupedExpression:
		return cg.generateExpression(e.Expr)
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
		return cg.propagateOrWrapError(err, stmt, "failed to generate value for let statement: %s", err.Error())
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
		return cg.propagateOrWrapError(err, stmt, "failed to generate value for assignment statement: %s", err.Error())
	}
	newType := stmt.Value.GetType()

	if !variable.varType.Equals(newType) {
		return cg.addErrorAtNode(stmt, "cannot assign %s type to variable of type %s", newType, variable.varType)
	}

	cg.builder.NewStore(newValue, variable.alloca)
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
