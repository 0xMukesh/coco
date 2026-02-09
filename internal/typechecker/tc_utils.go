package typechecker

import "github.com/0xmukesh/coco/internal/ast"

func (tc *TypeChecker) isLiteralExpression(expr ast.Expression) bool {
	switch expr.(type) {
	case *ast.IntegerExpression, *ast.FloatExpression, *ast.BooleanExpression, *ast.StringExpression:
		return true
	default:
		return false
	}
}
