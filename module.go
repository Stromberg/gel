package gel

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

type Script struct {
	Name   string
	Source string
}

// Module is a description of a module containing functions
type Module struct {
	Name        string
	Description string
	Funcs       []*Func
	LispFuncs   []*LispFunc
	Scripts     []*Script
}

func (f *Func) Repr() string {
	return fmt.Sprintf("%v\n%v\n%v", f.Name, f.Signature, f.Description)
}

func (f *LispFunc) Repr() string {
	return fmt.Sprintf("%v\n%v\n%v", f.Name, f.Signature, f.Description)
}
