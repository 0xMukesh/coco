package cotypes

// import (
// 	"fmt"

// 	"github.com/0xmukesh/coco/internal/ast"
// 	"github.com/0xmukesh/coco/internal/tokens"
// )

// type TypeChecker struct {
// 	env    *TypeEnvironment
// 	errors []string
// }

// func NewTypeChecker(env *TypeEnvironment) *TypeChecker {
// 	return &TypeChecker{
// 		env:    env,
// 		errors: []string{},
// 	}
// }

// func (tc *TypeChecker) addError(format string, args ...any) {
// 	tc.errors = append(tc.errors, fmt.Sprintf(format, args...))
// }

// func (tc *TypeChecker) checkBinaryExpression(expr *ast.BinaryExpression) Type {
// 	leftType := tc.checkExpression(expr.Left)
// 	rightType := tc.checkExpression(expr.Right)
// 	leftTypeCategory := GetTypeCategory(leftType)
// 	rightTypeCategory := GetTypeCategory(rightType)

// 	op := expr.Operator.Type

// 	// numeric arithmetic
// 	if op == tokens.PLUS || op == tokens.MINUS || op == tokens.STAR || op == tokens.SLASH {
// 		if leftTypeCategory != CategsryNumeric || rightTypeCategory != CategoryNumeric {
// 			tc.addError("cannot perform %s operator on %s and %s", op, leftType, rightType)
// 			return leftType
// 		}

// 		if leftType.Equals(FloatType{}) || rightType.Equals(FloatType{}) {

// 		}
// 	}

// 	// string concatenation
// 	if op == tokens.PLUS && leftType.Equals(StringType{}) && rightType.Equals(StringType{}) {
// 		return StringType{}
// 	}

// 	// logical operators
// 	if op == tokens.EQUALS || op == tokens.NOT_EQUALS {
// 		if !leftType.Equals(rightType) {
// 			tc.addError("")
// 		}
// 	}

// 	tc.addError("unknown binary operator %s", op)
// 	return leftType
// }

// func (tc *TypeChecker) checkExpression(expr ast.Expression) Type {
// 	switch e := expr.(type) {
// 	case *ast.IntegerExpression:
// 		return IntType{}
// 	case *ast.FloatExpression:
// 		return FloatType{}
// 	case *ast.StringExpression:
// 		return StringType{}
// 	case *ast.BinaryExpression:
// 		return tc.checkBinaryExpression(e)
// 	default:
// 		tc.addError("unknown expression. type %T", expr)
// 		return nil
// 	}
// }

// func (tc *TypeChecker) checkStatement(stmt ast.Statement) {
// 	switch s := stmt.(type) {
// 	case *ast.ExpressionStatement:
// 		tc.checkExpression(s.Expr)
// 	default:
// 		tc.addError("unknown statement. type %T", stmt)
// 	}
// }

// func (tc *TypeChecker) CheckProgram(program *ast.Program) *ast.Program {
// 	typedProgram := &ast.Program{}
// 	typedProgram.Statements = []ast.Statement{}

// 	for _, stmt := range program.Statements {
// 		tc.checkStatement(stmt)
// 	}

// 	return typedProgram
// }

// func (tc *TypeChecker) Errors() []string {
// 	return tc.errors
// }

// func (tc *TypeChecker) HasErrors() bool {
// 	return len(tc.errors) > 0
// }
