package typechecker

import (
	"fmt"
	"go/ast"
)

type TypeCheckerError struct {
	message string
	node    ast.Node
}

func isTypecheckerError(err error) bool {
	_, ok := err.(*TypeCheckerError)
	return ok
}

func (e *TypeCheckerError) Error() string {
	return fmt.Sprintf("typechecker error at %q: %s", e.node, e.message)
}
