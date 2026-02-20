package codegen

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/env"
)

func (cg *Codegen) generateStatement(stmt ast.Statement) (err error) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		_, err = cg.generateExpression(s.Expr)
		return err
	case *ast.LetStatement:
		return cg.generateLetStatement(s)
	case *ast.AssignmentStatement:
		return cg.generateAssignmentStatement(s)
	case *ast.BlockStatement:
		return cg.generateBlockStatement(s)
	case *ast.ReturnStatement:
		return cg.generateReturnStatement(s)
	case *ast.WhileStatement:
		return cg.generateWhileStatement(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", s)
	}
}

func (cg *Codegen) generateLetStatement(stmt *ast.LetStatement) error {
	varName := stmt.Identifier.String()
	exists := cg.scope.Has(varName)
	if exists {
		return cg.propagateOrWrapError(nil, stmt, "cannot redeclare %q variable", varName)
	}

	varType := stmt.Value.GetType()

	initValue, err := cg.generateExpression(stmt.Value)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to generate value for let statement: %s", err.Error())
	}

	alloca := cg.builder.NewAlloca(initValue.Type())
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
		return cg.propagateOrWrapError(nil, stmt, "cannot assign to undefined variable: %s", varName)
	}

	newValue, err := cg.generateExpression(stmt.Value)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to codegen value for assignment statement: %s", err.Error())
	}
	newType := stmt.Value.GetType()

	if !variable.typ.Equals(newType) {
		return cg.propagateOrWrapError(nil, stmt, "cannot assign %q type to variable of type %q", newType, variable.typ)
	}

	cg.builder.NewStore(newValue, variable.alloca)
	return nil
}

func (cg *Codegen) generateBlockStatement(stmt *ast.BlockStatement) error {
	previousScope := cg.scope
	cg.scope = env.NewEnvironmentWithParent(previousScope)

	for _, s := range stmt.Statements {
		if err := cg.generateStatement(s); err != nil {
			return err
		}
	}

	cg.scope = previousScope
	return nil
}

func (cg *Codegen) generateReturnStatement(stmt *ast.ReturnStatement) error {
	val, err := cg.generateExpression(stmt.Expr)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt, "failed to codegen value for return statement: %s", err.Error())
	}

	cg.blockReturnValue = val
	return nil
}

func (cg *Codegen) generateWhileStatement(stmt *ast.WhileStatement) error {
	conditionBlock := cg.mainFn.NewBlock("")
	bodyBlock := cg.mainFn.NewBlock("")
	exitBlock := cg.mainFn.NewBlock("")

	cg.builder.NewBr(conditionBlock)

	cg.builder = conditionBlock
	val, err := cg.generateExpression(stmt.Condition)
	if err != nil {
		return cg.propagateOrWrapError(err, stmt.Condition, "failed to codegen while statement's condition: %s", err.Error())
	}
	cg.builder.NewCondBr(val, bodyBlock, exitBlock)

	cg.builder = bodyBlock
	if err := cg.generateStatement(stmt.Body); err != nil {
		return cg.propagateOrWrapError(err, stmt.Body, "failed to codegen while statement's body: %s", err.Error())
	}
	cg.builder.NewBr(conditionBlock)

	cg.builder = exitBlock

	return nil
}
