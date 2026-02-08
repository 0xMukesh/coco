package typechecker

import (
	"github.com/0xmukesh/coco/internal/ast"
	"github.com/0xmukesh/coco/internal/env"
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
	default:
		return tc.propagateOrWrapError(nil, s, "unknown statement type: %T", s)
	}
}

func (tc *TypeChecker) checkLetStatement(stmt *ast.LetStatement) error {
	varName := stmt.Identifier.String()
	if tc.env.Has(varName) {
		return tc.propagateOrWrapError(nil, stmt, "cannot redeclare variable: %s", varName)
	}

	varType, err := tc.checkExpression(stmt.Value)
	if err != nil {
		return tc.propagateOrWrapError(err, stmt, "failed to typecheck let statement value: %s", err.Error())
	}

	tc.env.Set(varName, varType)
	return nil
}

func (tc *TypeChecker) checkAssignmentStatement(stmt *ast.AssignmentStatement) error {
	varName := stmt.Identifier.String()
	_, exists := tc.env.Get(varName)
	if !exists {
		return tc.propagateOrWrapError(nil, stmt, "unknown identifier: %s", varName)
	}

	_, err := tc.checkExpression(stmt.Value)
	if err != nil {
		return tc.propagateOrWrapError(err, stmt, "failed to typecheck assignment statement value: %s", err.Error())
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
