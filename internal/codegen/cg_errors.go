package codegen

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
)

type CodegenError struct {
	message string
	node    ast.Node
}

func (e *CodegenError) Error() string {
	if e.node != nil {
		return fmt.Sprintf("codegen error at %q: %s", e.node, e.message)
	} else {
		return fmt.Sprintf("codegen error: %s", e.message)
	}
}

func (cg *Codegen) propagateOrWrapError(err error, node ast.Node, msg string, args ...any) error {
	_, ok := err.(*CodegenError)
	if ok {
		return err
	}

	cgErr := &CodegenError{
		message: fmt.Sprintf(msg, args...),
		node:    node,
	}

	cg.errors = append(cg.errors, cgErr)
	return cgErr
}
