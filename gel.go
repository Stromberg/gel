// Package gel contains the core functionality for the Go Embedded Lisplike/Language.
package gel

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/Stromberg/gel/ast"
	"golang.org/x/crypto/ssh/terminal"
)

// Gel is the language expression handler.
type Gel struct {
	node ast.Node
	fset *ast.FileSet
	code string
}

// New creates a new Gel from a code string
func New(code string) (*Gel, error) {
	fset := NewFileSet()

	node, err := ParseString(fset, "", code)
	if err != nil {
		return nil, err
	}

	return &Gel{node, fset, code}, nil
}

func (g *Gel) Code() string {
	return g.code
}

// Missing returns the symbols that are missing
// in the environment in order to evaluate the expression.
func (g *Gel) Missing(env *Env) []string {
	scope := g.scope(env)
	return missing(scope, g.node)
}

// Eval evaluates the expression in the given environment.
func (g *Gel) Eval(env *Env) (interface{}, error) {
	scope := g.scope(env)
	return scope.Eval(g.node)
}

func (g *Gel) scope(env *Env) *Scope {
	scope := NewScope(g.fset)
	env.fillScope(scope)
	return scope
}

func (g *Gel) Repl(env *Env) error {
	scope := g.scope(env)

	state, err := terminal.MakeRaw(1)
	if err != nil {
		return err
	}
	defer terminal.Restore(1, state)

	t := terminal.NewTerminal(os.Stdout, "> ")
	unclosed := ""
	for {
		line, err := t.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if line == "exit" {
			break
		}
		if unclosed != "" {
			line = unclosed + "\n" + line
		}
		unclosed = ""
		t.SetPrompt("> ")
		node, err := ParseString(g.fset, "", line)
		if err != nil {
			if strings.HasSuffix(err.Error(), "missing )") {
				unclosed = line
				t.SetPrompt(". ")
				continue
			}
			fmt.Println(err)
			continue
		}
		value, err := scope.Eval(node)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if value != nil {
			if reflect.TypeOf(value).Kind() == reflect.Func {
				fmt.Println("#func")
			} else if v, ok := value.([]interface{}); ok {
				if len(v) == 0 {
					fmt.Println("()")
				} else {
					fmt.Print("(list")
					for _, e := range v {
						fmt.Printf(" %#v", e)
					}
					fmt.Println(")")
				}
			} else if v, ok := value.([]float64); ok {
				if len(v) == 0 {
					fmt.Println("()")
				} else {
					fmt.Print("(vec")
					for _, e := range v {
						fmt.Printf(" %#v", e)
					}
					fmt.Println(")")
				}
			} else {
				fmt.Printf("%#v\n", value)
			}
		}
	}
	fmt.Println()
	return nil
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
