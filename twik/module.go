package twik

// Func is a description of a Module function.
type Func struct {
	Name        string
	Signature   string
	Description string
	F           interface{}
}

// Module is a description of a module containing functions
type Module struct {
	Name        string
	Description string
	Funcs       []*Func
}

func (m *Module) Load(scope *Scope) {
	for _, f := range m.Funcs {
		scope.SetOrCreate(f.Name, f.F)
	}
}
