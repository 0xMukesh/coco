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
	if e.node != nil {
		return fmt.Sprintf("[line: %d] typechecker error at %s: %s", e.node.Location(), e.node, e.message)
	} else {
		return fmt.Sprintf("[line: %d] typechecker error: %s", e.node.Location(), e.message)
	}
}

func (tc *TypeChecker) propagateOrWrapError(err error, node ast.Node, msg string, args ...any) error {
	_, ok := err.(*TypeCheckerError)
	if ok {
		return err
	}

	tcErr := &TypeCheckerError{
		message: fmt.Sprintf(msg, args...),
		node:    node,
	}

	tc.errors = append(tc.errors, tcErr)
	return tcErr
}
