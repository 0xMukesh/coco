package codegen

import (
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
)

type CodegenError struct {
	message string
	node    ast.Node
}

func isCodegenError(err error) bool {
	_, ok := err.(*CodegenError)
	return ok
}

func (e *CodegenError) Error() string {
	return fmt.Sprintf("codegen error at %q: %s", e.node, e.message)
}

func (e *CodegenError) Node() ast.Node {
	return e.node
}
