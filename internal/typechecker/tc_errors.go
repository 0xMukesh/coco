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
		return fmt.Sprintf("typechecker error at %q: %s", e.node, e.message)
	} else {
		return fmt.Sprintf("typechecker error: %s", e.message)
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
