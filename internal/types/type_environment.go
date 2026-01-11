package types

type TypeEnvironment struct {
	store  map[string]Type
	parent *TypeEnvironment
}

func NewTypeEnvironment() *TypeEnvironment {
	return &TypeEnvironment{
		store:  make(map[string]Type),
		parent: nil,
	}
}

func NewEnclosedTypeEnvironment(parent *TypeEnvironment) *TypeEnvironment {
	env := NewTypeEnvironment()
	env.parent = parent
	return env
}

func (te *TypeEnvironment) Get(name string) (Type, bool) {
	typ, ok := te.store[name]
	if !ok && te.parent != nil {
		typ, ok = te.parent.Get(name)
	}

	return typ, ok
}

func (te *TypeEnvironment) Set(name string, typ Type) {
	te.store[name] = typ
}
