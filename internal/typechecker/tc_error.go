package typechecker

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
)

type TypeCheckerError struct {
	message string
	node    ast.Node
}

func (e *TypeCheckerError) Error() string {
	return fmt.Sprintf("typechecker error at %q: %s", e.node, e.message)
}

func isTypecheckerError(err error) bool {
	_, ok := err.(*TypeCheckerError)
	return ok
}

func (e *TypeCheckerError) Node() ast.Node {
	return e.node
}

func (tc *TypeChecker) addErrorAtNode(node ast.Node, msg string, args ...any) error {
	err := &TypeCheckerError{
		message: fmt.Sprintf(msg, args...),
		node:    node,
	}

	tc.errors = append(tc.errors, err)
	return err
}

func (tc *TypeChecker) addError(msg string, args ...any) error {
	err := fmt.Errorf(msg, args...)
	tc.errors = append(tc.errors, err)
	return err
}

func (tc *TypeChecker) propagateOrWrapError(err error, node ast.Node, msg string, args ...any) error {
	if isTypecheckerError(err) {
		return err
	}

	return tc.addErrorAtNode(node, msg, args...)
}
