package object

import (
	"bytes"
	"fmt"

	"github.com/0xmukesh/coco/internal/ast"
)

type ObjectType string

const (
	INT_OBJECT      = "INTEGER"
	FLOAT_OBJECT    = "FLOAT"
	BOOL_OBJECT     = "BOOLEAN"
	NULL_OBJECT     = "NULL"
	RETURN_OBJECT   = "RETURN"
	FUNCTION_OBJECT = "FUNCTION"
	ERROR_OBJECT    = "ERROR"
)

type Object interface {
	Type() string
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() string {
	return INT_OBJECT
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Float struct {
	Value float64
}

func (f *Float) Type() string {
	return FLOAT_OBJECT
}
func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() string {
	return BOOL_OBJECT
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() string {
	return NULL_OBJECT
}
func (n *Null) Inspect() string {
	return "null"
}

type Return struct {
	Value Object
}

func (r *Return) Type() string {
	return RETURN_OBJECT
}
func (r *Return) Inspect() string {
	return r.Inspect()
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() string {
	return FUNCTION_OBJECT
}
func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("fn(")
	for i, p := range f.Parameters {
		out.WriteString(p.String())
		out.WriteString(",")

		if i < len(f.Parameters)-2 {
			out.WriteString(" ")
		}
	}
	out.WriteString(") ")

	out.WriteString("{ " + f.Body.String() + " }")

	return out.String()
}

type Error struct {
	Message string
}

func (e *Error) Type() string {
	return ERROR_OBJECT
}
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
