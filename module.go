package gel

// Func is a description of a Module function.
type Func struct {
	Name        string
	Signature   string
	Description string
	F           func(args []interface{}) (interface{}, error)
}

// Module is a description of a module containing functions
type Module struct {
	Name        string
	Description string
	Funcs       []*Func
}

// ModuleRepo is container for modules
type ModuleRepo struct {
	modules []*Module
}

// NewModuleRepo creates a new ModuleRepo.
func NewModuleRepo(modules ...*Module) *ModuleRepo {
	return &ModuleRepo{modules}
}

// Add adds a new module to the ModuleRepo.
func (r *ModuleRepo) Add(module *Module) {
	r.modules = append(r.modules, module)
}

// Module finds the Module with the given name.
func (r *ModuleRepo) Module(name string) *Module {
	for _, m := range r.modules {
		if m.Name == name {
			return m
		}
	}

	return nil
}

// Function finds the Func with the given name and the Module it is in.
func (r *ModuleRepo) Function(name string) (*Func, *Module) {
	for _, m := range r.modules {
		for _, f := range m.Funcs {
			if f.Name == name {
				return f, m
			}
		}
	}

	return nil, nil
}
