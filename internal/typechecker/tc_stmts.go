package typechecker

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/env"
	cotypes "github.com/0xmukesh/coco/internal/types"
)

func (tc *TypeChecker) checkStatement(stmt ast.Statement) (err error) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		_, err = tc.checkExpression(s.Expr)
		return err
	case *ast.LetStatement:
		return tc.checkLetStatement(s)
	case *ast.AssignmentStatement:
		return tc.checkAssignmentStatement(s)
	case *ast.BlockStatement:
		return tc.checkBlockStatement(s)
	case *ast.ReturnStatement:
		_, err = tc.checkExpression(s.Expr)
		return err
	case *ast.WhileStatement:
		return tc.checkWhileStatement(s)
	case *ast.BreakStatement:
		return nil
	default:
		return tc.propagateOrWrapError(nil, s, "unknown statement type: %T", s)
	}
}

func (tc *TypeChecker) checkLetStatement(stmt *ast.LetStatement) error {
	varName := stmt.Identifier.String()
	if tc.env.Has(varName) {
		return tc.propagateOrWrapError(nil, stmt.Identifier, "cannot redeclare variable: %s", varName)
	}

	varType, err := tc.checkExpression(stmt.Value)
	if err != nil {
		return tc.propagateOrWrapError(err, stmt.Value, "failed to typecheck let statement value: %s", err.Error())
	}

	if varType.Equals(cotypes.VoidType{}) {
		return tc.propagateOrWrapError(err, stmt.Identifier, "cannot bind variable %q to 'void' type", varName)
	}

	tc.env.Set(varName, varType)
	return nil
}

func (tc *TypeChecker) checkAssignmentStatement(stmt *ast.AssignmentStatement) error {
	varName := stmt.Identifier.String()
	originalType, exists := tc.env.Get(varName)
	if !exists {
		return tc.propagateOrWrapError(nil, stmt.Identifier, "unknown identifier: %s", varName)
	}

	newType, err := tc.checkExpression(stmt.Value)
	if err != nil {
		return tc.propagateOrWrapError(err, stmt.Value, "failed to typecheck assignment statement value: %s", err.Error())
	}

	if !originalType.Equals(newType) {
		return tc.propagateOrWrapError(err, stmt.Value, "cannot assign value of type %q to value of type %q", originalType, newType)
	}

	return nil
}

func (tc *TypeChecker) checkBlockStatement(stmt *ast.BlockStatement) error {
	tc.env = env.NewEnvironmentWithParent(tc.env)

	for _, s := range stmt.Statements {
		if err := tc.checkStatement(s); err != nil {
			return tc.propagateOrWrapError(err, s, "failed to typecheck statement: %s", err.Error())
		}
	}

	tc.env = tc.env.Parent()
	return nil
}

func (tc *TypeChecker) checkWhileStatement(stmt *ast.WhileStatement) error {
	conditionType, err := tc.checkExpression(stmt.Condition)
	if err != nil {
		return tc.propagateOrWrapError(err, stmt.Condition, "failed to typecheck expression: %s", err.Error())
	}

	if !conditionType.Equals(cotypes.BoolType{}) {
		return tc.propagateOrWrapError(nil, stmt.Condition, "non-boolean expression in while condition")
	}

	if err := tc.checkStatement(stmt.Body); err != nil {
		return tc.propagateOrWrapError(err, stmt.Body, "failed to typecheck statement: %s", err.Error())
	}

	return nil
}
