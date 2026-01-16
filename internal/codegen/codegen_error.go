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

func (cg *Codegen) addErrorAtNode(node ast.Node, msg string, args ...any) error {
	err := &CodegenError{
		message: fmt.Sprintf(msg, args...),
		node:    node,
	}

	cg.errors = append(cg.errors, err)
	return err
}

func (cg *Codegen) addError(msg string, args ...any) error {
	err := fmt.Errorf(msg, args...)
	cg.errors = append(cg.errors, err)
	return err
}

func (cg *Codegen) propagateOrWrapError(err error, node ast.Node, msg string, args ...any) error {
	if isCodegenError(err) {
		return err
	}

	return cg.addErrorAtNode(node, msg, args...)
}

func (e *CodegenError) Error() string {
	return fmt.Sprintf("codegen error at %q: %s", e.node, e.message)
}

func (e *CodegenError) Node() ast.Node {
	return e.node
}
