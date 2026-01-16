package cotypes

type TypeCategory int

type Type interface {
	String() string
	Equals(Type) bool
}

const (
	CategoryUnknown TypeCategory = iota
	CategoryNumeric
)

type IntType struct{}

func (i IntType) String() string { return "int" }
func (i IntType) Equals(t Type) bool {
	_, ok := t.(IntType)
	return ok
}

type FloatType struct{}

func (f FloatType) String() string { return "float" }
func (f FloatType) Equals(t Type) bool {
	_, ok := t.(FloatType)
	return ok
}

type BoolType struct{}

func (b BoolType) String() string { return "bool" }
func (b BoolType) Equals(t Type) bool {
	_, ok := t.(BoolType)
	return ok
}

type StringType struct{}

func (s StringType) String() string { return "string" }
func (s StringType) Equals(t Type) bool {
	_, ok := t.(StringType)
	return ok
}

type VoidType struct{}

func (v VoidType) String() string { return "void" }
func (v VoidType) Equals(t Type) bool {
	_, ok := t.(VoidType)
	return ok
}

func GetTypeCategory(T Type) TypeCategory {
	switch T {
	case FloatType{}:
		return CategoryNumeric
	case IntType{}:
		return CategoryNumeric
	default:
		return CategoryUnknown
	}
}
