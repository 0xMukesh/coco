package object

import "fmt"

type ObjectType string

const (
	OBJECT_INT   = "INTEGER"
	OBJECT_FLOAT = "FLOAT"
	OBJECT_BOOL  = "BOOLEAN"
	OBJECT_NULL  = "NULL"
)

type Object interface {
	Type() string
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() string {
	return OBJECT_INT
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Float struct {
	Value float64
}

func (f *Float) Type() string {
	return OBJECT_FLOAT
}
func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() string {
	return OBJECT_BOOL
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

func (n *Null) Type() string {
	return OBJECT_NULL
}
func (n *Null) Inspect() string {
	return "null"
}
