package twik

var registeredModules = []*Module{}

// Modules returns all registered modules
func Modules() []*Module {
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
func RegisterModule(module *Module) {
	registeredModules = append(registeredModules, module)
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
