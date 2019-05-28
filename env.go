package gel

import (
	"github.com/Stromberg/gel/twik"
)

// Env contains the variables, functions and modules that
// should be used for Gel expression evaluation
type Env struct {
	vars    map[string]interface{}
	modules []*Module
}

// NewEnv creates a new Env.
func NewEnv() *Env {
	return &Env{vars: make(map[string]interface{}, 0)}
}

// Copy creates a copy of the current Env.
func (e *Env) Copy() *Env {
	vars := make(map[string]interface{})
	for k, v := range e.vars {
		vars[k] = v
	}

	return &Env{
		vars:    vars,
		modules: append([]*Module(nil), e.modules...),
	}
}

// AddVar adds a variable or function to the Env.
func (e *Env) AddVar(name string, value interface{}) {
	e.vars[name] = value
}

// AddModule adds a module to the Env.
func (e *Env) AddModule(module ...*Module) {
	e.modules = append(e.modules, module...)
}

func (e *Env) fillScope(scope *twik.Scope) {
	for _, m := range e.modules {
		for _, f := range m.Funcs {
			err := scope.Set(f.Name, f.F)
			if err != nil {
				scope.Create(f.Name, f.F)
			}
		}
	}

	for k, v := range e.vars {
		err := scope.Set(k, v)
		if err != nil {
			scope.Create(k, v)
		}
	}
}
