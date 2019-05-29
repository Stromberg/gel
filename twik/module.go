package twik

import "fmt"

// Func is a description of a Module function.
type Func struct {
	Name        string
	Signature   string
	Description string
	F           interface{}
}

// LispFunc is a description of a Module function in lisp.
type LispFunc struct {
	Name        string
	Signature   string
	Description string
	F           string
}

// Module is a description of a module containing functions
type Module struct {
	Name        string
	Description string
	Funcs       []*Func
	LispFuncs   []*LispFunc
}

func (m *Module) LoadLisp() string {
	s := ""
	for _, f := range m.LispFuncs {
		s += fmt.Sprintf("(set %s %s)\n", f.Name, f.F)
	}
	return s
}

func (m *Module) Load(scope *Scope) {
	for _, f := range m.Funcs {
		scope.SetOrCreate(f.Name, f.F)
	}
}
