package gel

import (
	"fmt"

	"github.com/ryanuber/go-glob"
)

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

// RegisterModules registers a new Module.
func RegisterModules(modules ...*Module) {
	for _, m := range modules {
		registeredModules = append(registeredModules, m)
	}
}

func AllFunctionNames() []string {
	res := []string{}

	for _, m := range registeredModules {
		for _, f := range m.Funcs {
			res = append(res, f.Name)
		}
		for _, f := range m.LispFuncs {
			res = append(res, f.Name)
		}
	}

	return res
}

func MatchingFuncNames(expr string) []string {
	funcs := AllFunctionNames()

	res := []string{}
	for _, f := range funcs {
		if glob.Glob(expr, f) {
			res = append(res, f)
		}
	}

	return res
}

func FunctionRepr(name string) string {
	for _, m := range registeredModules {
		for _, f := range m.Funcs {
			if f.Name == name {
				return f.Repr()
			}
		}
		for _, f := range m.LispFuncs {
			if f.Name == name {
				return f.Repr()
			}
		}
	}

	return fmt.Sprintf("Function \"%v\" not found", name)
}
