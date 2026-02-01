package codegen

import (
	"fmt"
	"strings"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/env"
	"github.com/0xmukesh/coco/internal/tokens"
	cotypes "github.com/0xmukesh/coco/internal/types"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Scope = *env.Environent[ScopeItem]
type ScopeItem struct {
	alloca *ir.InstAlloca
	typ    cotypes.Type
}

type Codegen struct {
	module  *ir.Module
	mainFn  *ir.Func
	builder *ir.Block
	ret     value.Value

	scope        Scope
	runtimeFuncs map[string]*ir.Func
	globalDefs   map[string]*ir.Global

	nameCounter int
	errors      []error
}

func New() *Codegen {
	module := ir.NewModule()
	mainFn := module.NewFunc("main", types.I32)
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
	case *ast.BlockStatement:
		previousScope := cg.scope
		cg.scope = env.NewEnvironmentWithParent(previousScope)

		for _, stmt := range s.Statements {
			if err := cg.generateStatement(stmt); err != nil {
				return err
			}
		}

		cg.scope = previousScope
		return nil
	}

	return nil
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
	case *ast.IfExpression:
		return cg.generateIfExpression(e)
	default:
		return nil, cg.addErrorAtNode(expr, "unsupported expression type")
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
	variable, exists := cg.scope.Get(expr.Literal)
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
	funcName := expr.Identifier.String()
	printfFunc, ok := cg.runtimeFuncs[funcName]
	if !ok {
		printfFunc = cg.setupPrintfRuntimeFunc()
	}

	var fmtStr strings.Builder
	var boolArgs []value.Value

	for i, arg := range expr.Arguments {
		if i > 0 {
			fmtStr.WriteString(" ")
		}

		switch arg.GetType().(type) {
		case cotypes.IntType:
			fmtStr.WriteString("%ld")
		case cotypes.FloatType:
			fmtStr.WriteString("%g")
		case cotypes.BoolType:
			fmtStr.WriteString("%s")

			trueStr, ok := cg.globalDefs[TRUE_GLOBAL_DEF_NAME]
			if !ok {
				trueStr = cg.setupTrueGlobalDef()
			}

			falseStr, ok := cg.globalDefs[FALSE_GLOBAL_DEF_NAME]
			if !ok {
				falseStr = cg.setupFalseStrGlobalDef()
			}

			truePtr := cg.builder.NewGetElementPtr(
				types.NewArray(uint64(len("true\x00")), types.I8),
				trueStr,
				constant.NewInt(types.I64, 0),
				constant.NewInt(types.I64, 0),
			)

			falsePtr := cg.builder.NewGetElementPtr(
				types.NewArray(uint64(len("false\x00")), types.I8),
				falseStr,
				constant.NewInt(types.I64, 0),
				constant.NewInt(types.I64, 0),
			)

			boolValue, err := cg.generateExpression(arg)
			if err != nil {
				return nil, err
			}

			strPtr := cg.builder.NewSelect(boolValue, truePtr, falsePtr)
			boolArgs = append(boolArgs, strPtr)
		}
	}
	fmtStr.WriteString("\n\x00")

	fmtGlobalDef := cg.module.NewGlobalDef(fmt.Sprintf(".fmt.%d", cg.nameCounter), constant.NewCharArrayFromString(fmtStr.String()))
	fmtGlobalDef.Immutable = true
	fmtGlobalDef.Linkage = enum.LinkagePrivate
	fmtGlobalDef.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	cg.nameCounter++

	fmtPtr := cg.builder.NewGetElementPtr(
		types.NewArray(uint64(len(fmtStr.String())), types.I8),
		fmtGlobalDef,
		constant.NewInt(types.I64, 0),
		constant.NewInt(types.I64, 0),
	)

	args := []value.Value{fmtPtr}
	boolIdx := 0
	for _, arg := range expr.Arguments {
		switch arg.GetType().(type) {
		case cotypes.BoolType:
			// use precomputed boolean str value
			args = append(args, boolArgs[boolIdx])
			boolIdx++
		default:
			v, err := cg.generateExpression(arg)
			if err != nil {
				return nil, err
			}

			args = append(args, v)
		}
	}

	cg.builder.NewCall(printfFunc, args...)
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

func (cg *Codegen) generateIfExpression(expr *ast.IfExpression) (value.Value, error) {
	condition, err := cg.generateExpression(expr.Condition)
	if err != nil {
		return nil, cg.propagateOrWrapError(err, expr, "failed to generate value for if-branch condition: %s", err.Error())
	}

	ifTrue := cg.mainFn.NewBlock("")
	ifFalse := cg.mainFn.NewBlock("")
	merge := cg.mainFn.NewBlock("")

	cg.builder.NewCondBr(condition, ifTrue, ifFalse)

	cg.builder = ifTrue
	cg.generateStatement(expr.Consequence)
	cg.builder.NewBr(merge)

	cg.builder = ifFalse
	cg.generateStatement(expr.Alternative)
	cg.builder.NewBr(merge)

	cg.builder = merge
	return nil, nil
}

func (cg *Codegen) generateLetStatement(stmt *ast.LetStatement) error {
	varName := stmt.Identifier.String()
	exists := cg.scope.Has(varName)
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
	cg.builder.NewStore(initValue, alloca)

	cg.scope.Set(varName, ScopeItem{
		alloca: alloca,
		typ:    varType,
	})
	return nil
}

func (cg *Codegen) generateAssignmentStatement(stmt *ast.AssignmentStatement) error {
	varName := stmt.Identifier.String()
	variable, exists := cg.scope.Get(varName)
	if !exists {
		return cg.addErrorAtNode(stmt, "cannot assign to undefined variable: %s", varName)
	}

	newValue, err := cg.generateExpression(stmt.Value)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to generate value for assignment statement: %s", err.Error())
	}
	newType := stmt.Value.GetType()

	if !variable.typ.Equals(newType) {
		return cg.addErrorAtNode(stmt, "cannot assign %s type to variable of type %s", newType, variable.typ)
	}

	cg.builder.NewStore(newValue, variable.alloca)
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
