// Package gel implements a tiny embeddable language for Go.
//
// For details, see the blog post:
//
//     http://blog.labix.org/2013/07/16/twik-a-tiny-language-for-go
//
package gel

import (
	"fmt"

	"github.com/Stromberg/gel/ast"
)

// Scope is an environment where twik logic may be evaluated in.
type Scope struct {
	parent *Scope
	fset   *ast.FileSet
	vars   map[string]interface{}
}

// Error holds an error and the source position where the error was found.
type Error struct {
	Err     error
	PosInfo *ast.PosInfo
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %v", e.PosInfo, e.Err)
}

// NewScope returns a new scope for evaluating logic that was parsed into fset.
func NewScope(fset *ast.FileSet) (*Scope, error) {
	modules := Modules()

	vars := make(map[string]interface{})

	for _, m := range modules {
		for _, f := range m.Funcs {
			vars[f.Name] = f.F
		}
	}

	scope := &Scope{fset: fset, vars: vars}

	for _, m := range modules {
		for _, f := range m.LispFuncs {
			expr := fmt.Sprintf("(var %s %s)", f.Name, f.F)
			node, err := ParseString(fset, fmt.Sprintf("%v:%v", m.Name, f.Name), expr)
			if err != nil {
				return nil, err
			}

			_, err = scope.Eval(node)
			if err != nil {
				return nil, err
			}
		}

		for _, s := range m.Scripts {
			node, err := ParseString(fset, fmt.Sprintf("%v:%v", m.Name, s.Name), s.Source)
			if err != nil {
				return nil, err
			}

			_, err = scope.Eval(node)
			if err != nil {
				return nil, err
			}
		}
	}

	return scope, nil
}

// Create defines a new symbol with the given value in the s scope.
// It is an error to redefine an existent symbol.
func (s *Scope) Create(symbol string, value interface{}) error {
	if _, ok := s.vars[symbol]; ok {
		return fmt.Errorf("symbol already defined in current scope: %s", symbol)
	}
	if s.vars == nil {
		s.vars = make(map[string]interface{})
	}
	s.vars[symbol] = value
	return nil
}

func (s *Scope) SetOrCreate(symbol string, value interface{}) {
	if err := s.Set(symbol, value); err != nil {
		_ = s.Create(symbol, value)
	}
}

// Set sets symbol to the given value in the shallowest scope it is defined in.
// It is an error to set an undefined symbol.
func (s *Scope) Set(symbol string, value interface{}) error {
	for s != nil {
		if _, ok := s.vars[symbol]; ok {
			s.vars[symbol] = value
			return nil
		}
		s = s.parent
	}
	return fmt.Errorf("cannot set undefined symbol: %s", symbol)
}

// Get returns the value of symbol in the shallowest scope it is defined in.
// It is an error to get the value of an undefined symbol.
func (s *Scope) Get(symbol string) (value interface{}, err error) {
	for s != nil {
		if value, ok := s.vars[symbol]; ok {
			return value, nil
		}
		s = s.parent
	}
	return nil, fmt.Errorf("undefined symbol: %s", symbol)
}

// Branch returns a new scope that has s as a parent.
func (s *Scope) Branch() *Scope {
	return &Scope{parent: s, fset: s.fset}
}

var emptyList = make([]interface{}, 0)
var emptyDict = make(map[interface{}]interface{}, 0)

func (s *Scope) errorAt(node ast.Node, err error) error {
	if _, ok := err.(*Error); ok {
		return err
	}
	return &Error{err, s.fset.PosInfo(node.Pos())}
}

func (s *Scope) Code(node ast.Node) string {
	return s.fset.Code(node)
}

// Eval evaluates node in the s scope and returns the resulting value.
func (s *Scope) Eval(node ast.Node) (value interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = s.errorAt(node, fmt.Errorf("%v", r))
		}
	}()

	switch node := node.(type) {
	case *ast.Symbol:
		value, err := s.Get(node.Name)
		if err != nil {
			return nil, s.errorAt(node, err)
		}
		return value, nil
	case *ast.Int:
		return node.Value, nil
	case *ast.Float:
		return node.Value, nil
	case *ast.String:
		return node.Value, nil
	case *ast.List:
		if len(node.Nodes) == 0 {
			return emptyList, nil
		}
		fn, err := s.Eval(node.Nodes[0])
		if err != nil {
			return nil, s.errorAt(node.Nodes[0], err)
		}
		value, err := s.call(fn, node.Nodes[1:])
		if err != nil {
			return nil, s.errorAt(node.Nodes[0], err)
		}
		return value, nil
	case *ast.ListList:
		if len(node.Nodes) == 0 {
			return emptyList, nil
		}
		vargs := make([]interface{}, len(node.Nodes))
		for i, arg := range node.Nodes {
			value, err := s.Eval(arg)
			if err != nil {
				return nil, err
			}
			vargs[i] = value
		}
		return NewList(vargs...)
	case *ast.DictList:
		if len(node.Nodes) == 0 {
			return emptyDict, nil
		}
		vargs := make([]interface{}, len(node.Nodes))
		for i, arg := range node.Nodes {
			value, err := s.Eval(arg)
			if err != nil {
				return nil, err
			}
			vargs[i] = value
		}
		return NewDict(vargs...)
	case *ast.Root:
		for _, node := range node.Nodes {
			value, err = s.Eval(node)
			if err != nil {
				return nil, s.errorAt(node, err)
			}
		}
		return value, nil
	}
	return nil, fmt.Errorf("support for %#v not yet implemeted", node)
}

func (s *Scope) call(fn interface{}, args []ast.Node) (value interface{}, err error) {
	if fn, ok := fn.(func(*Scope, []ast.Node) (interface{}, error)); ok {
		return fn(s, args)
	}

	vargs := make([]interface{}, len(args))
	for i, arg := range args {
		value, err := s.Eval(arg)
		if err != nil {
			return nil, err
		}
		vargs[i] = value
	}

	return Call(fn, vargs...)
}
