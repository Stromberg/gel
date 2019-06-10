package gel

var registeredModules = []*Module{}

var BasePath = "."

// Modules returns all registered modules
func Modules() (res []*Module) {
	return registeredModules
}

// FindModule finds the registered Module with the given name.
func FindModule(name string) *Module {
	for _, m := range registeredModules {
		if m.Name == name {
			return m
		}
	}

	return nil
}

// RegisterModule registers a new Module.
func RegisterModules(modules ...*Module) {
	for _, m := range modules {
		registeredModules = append(registeredModules, m)
	}
}

// FindFunction finds the Func with the given name and the Module it is in.
func FindFunction(name string) (*Func, *Module) {
	for _, m := range registeredModules {
		for _, f := range m.Funcs {
			if f.Name == name {
				return f, m
			}
		}
	}

	return nil, nil
}
