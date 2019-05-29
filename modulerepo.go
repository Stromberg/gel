package gel

var registeredModules = map[string]*Module{}

// Modules returns all registered modules
func Modules() (res []*Module) {
	for _, v := range registeredModules {
		res = append(res, v)
	}
	return
}

// AllModules returns all registered module names
func AllModuleNames() (res []string) {
	for k := range registeredModules {
		res = append(res, k)
	}
	return
}

// FindModule finds the registered Module with the given name.
func FindModule(name string) *Module {
	m, ok := registeredModules[name]
	if !ok {
		return nil
	}

	return m
}

// RegisterModule registers a new Module.
func RegisterModule(module *Module) {
	registeredModules[module.Name] = module
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
