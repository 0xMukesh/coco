package env

type Environent[T any] struct {
	store  map[string]T
	parent *Environent[T]
}

func NewEnvironment[T any]() *Environent[T] {
	return &Environent[T]{
		store:  make(map[string]T),
		parent: nil,
	}
}

func NewEnvironmentWithParent[T any](parent *Environent[T]) *Environent[T] {
	env := NewEnvironment[T]()
	env.parent = parent
	return env
}

func (te *Environent[T]) Get(name string) (T, bool) {
	v, ok := te.store[name]
	if !ok && te.parent != nil {
		v, ok = te.parent.Get(name)
	}

	return v, ok
}

func (te *Environent[T]) Has(name string) bool {
	_, ok := te.store[name]
	return ok
}

func (te *Environent[T]) Set(name string, v T) {
	te.store[name] = v
}

func (te *Environent[T]) Parent() *Environent[T] {
	return te.parent
}
