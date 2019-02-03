package dataext

type funcExtender struct {
	Extender

	id   string
	vars []string
	f    Func
}

// Func is a function that can make a calculation based on values in a Store
type Func func(store Store) (interface{}, error)

// NewFunc creates a new extender for id.
// vars are the variables in the store that it depends on.
// f is the calculation function.
func NewFunc(id string, vars []string, f Func) Extender {
	return &funcExtender{
		id:   id,
		vars: vars,
		f:    f,
	}
}

func (e *funcExtender) ID() string {
	return e.id
}

func (e *funcExtender) Missing(store Store) (res []string) {
	if _, ok := store.Get(e.id); ok {
		return nil
	}

	for _, v := range e.vars {
		if _, ok := store.Get(v); !ok {
			res = append(res, v)
		}
	}

	return res
}

func (e *funcExtender) Extend(store Store) error {
	if _, ok := store.Get(e.id); ok {
		return nil
	}

	data, err := e.f(store)
	if err != nil {
		return err
	}

	store.Set(e.id, data)

	return nil
}
