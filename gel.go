// Package gel contains the core functionality for the Go Embedded Lisplike/Language.
package gel

import (
	"fmt"

	"github.com/Stromberg/gel/ast"
)

// Gel is the language expression handler.
type Gel struct {
	node    ast.Node
	fset    *ast.FileSet
	code    string
	modules []*Module
}

// New creates a new Gel from a code string
func New(code string, modules []string) (*Gel, error) {
	fset := NewFileSet()

	ms := []*Module{}
	for _, name := range modules {
		m := FindModule(name)
		if m == nil {
			return nil, fmt.Errorf("No module named %s", name)
		}
	}

	node, err := ParseString(fset, "", code)
	if err != nil {
		return nil, err
	}

	return &Gel{node, fset, code, ms}, nil
}

func (g *Gel) Code() string {
	return g.code
}

// Missing returns the symbols that are missing
// in the environment in order to evaluate the expression.
func (g *Gel) Missing(env *Env) []string {
	scope := NewScope(g.fset, g.modules...)

	env.fillScope(scope)

	return missing(scope, g.node)
}

// Eval evaluates the expression in the given environment.
func (g *Gel) Eval(env *Env) (interface{}, error) {
	scope := NewScope(g.fset, g.modules...)

	env.fillScope(scope)
	return scope.Eval(g.node)
}

func missing(s *Scope, node ast.Node) []string {
	var res []string

	switch node := node.(type) {
	case *ast.Symbol:
		_, err := s.Get(node.Name)
		if err != nil {
			return append(res, node.Name)
		}
	case *ast.List:
		if n, ok := node.Nodes[0].(*ast.Symbol); ok {
			if n.Name == "var" {
				if n2, ok := node.Nodes[1].(*ast.Symbol); ok {
					_, err := s.Get(n2.Name)
					if err != nil {
						s.Create(n2.Name, 0.0)
					}
					for i := 2; i < len(node.Nodes); i++ {
						res = append(res, missing(s, node.Nodes[i])...)
					}
					return res
				}
			}
		}
		for _, n := range node.Nodes {
			res = append(res, missing(s, n)...)
		}
	case *ast.Root:
		for _, n := range node.Nodes {
			res = append(res, missing(s, n)...)
		}
	}
	return res
}
