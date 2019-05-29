package gel

import (
	"github.com/Stromberg/gel/twik"
)

// Env contains the variables, functions and modules that
// should be used for Gel expression evaluation
type Env struct {
	vars map[string]interface{}
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
		vars: vars,
	}
}

// AddVar adds a variable or function to the Env.
func (e *Env) AddVar(name string, value interface{}) {
	e.vars[name] = value
}

func (e *Env) fillScope(scope *twik.Scope) {
	for k, v := range e.vars {
		scope.SetOrCreate(k, v)
	}
}
