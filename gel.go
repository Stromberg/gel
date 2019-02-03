// Package gel contains the core functionality for the Go Embedded Lisplike/Language.
package gel

import (
	twik "gopkg.in/twik.v1"
	"gopkg.in/twik.v1/ast"
)

// Gel is the language expression handler.
type Gel struct {
	node ast.Node
	fset *ast.FileSet
}

// New creates a new Gel from a code string
func New(code string) (*Gel, error) {
	fset := twik.NewFileSet()
	node, err := twik.ParseString(fset, "", code)
	if err != nil {
		return nil, err
	}

	return &Gel{node, fset}, nil
}

// Missing returns the symbols that are missing
// in the environment in order to evaluate the expression.
func (g *Gel) Missing(env *Env) []string {
	scope := twik.NewScope(g.fset)
	env.fillScope(scope)

	return missing(scope, g.node)
}

// Eval evaluates the expression in the given environment.
func (g *Gel) Eval(env *Env) (interface{}, error) {
	scope := twik.NewScope(g.fset)
	env.fillScope(scope)
	return scope.Eval(g.node)
}

func missing(s *twik.Scope, node ast.Node) []string {
	var res []string

	switch node := node.(type) {
	case *ast.Symbol:
		_, err := s.Get(node.Name)
		if err != nil {
			res = append(res, node.Name)
		}
	case *ast.List:
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
