package gel

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
