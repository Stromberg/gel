package dataext

import (
	"fmt"

	"github.com/Stromberg/gel"
)

type gelExtender struct {
	Extender

	id             string
	gel            *gel.Gel
	env            *gel.Env
	err            error
	isFloat64Slice bool
	expr           string
}

// NewGel creates a new extender for id.
// exprStr is the gel expression to calculate.
// baseEnv is the environment to use to evaluate the expression.
func NewGel(id, exprStr string, baseEnv *gel.Env) Extender {
	gel, err := gel.New(exprStr, nil)
	if err != nil {
		return nil
	}
	return &gelExtender{
		id:             id,
		gel:            gel,
		env:            baseEnv.Copy(),
		isFloat64Slice: false,
		expr:           exprStr,
	}
}

// NewGelOnFloat64Slice creates a new extender for id.
// The extender is used to calculate a float64 slice.
// All the its dependencies in the Store must also be float64 slices.
// exprStr is the gel expression to calculate.
// baseEnv is the environment to use to evaluate the expression.
func NewGelOnFloat64Slice(id, exprStr string, baseEnv *gel.Env) *gelExtender {
	gel, err := gel.New(exprStr, nil)
	if err != nil {
		return nil
	}
	return &gelExtender{
		id:             id,
		gel:            gel,
		env:            baseEnv.Copy(),
		isFloat64Slice: true,
		expr:           exprStr,
	}
}

func (e *gelExtender) ID() string {
	return e.id
}

func (e *gelExtender) String() string {
	return e.expr
}

func (e *gelExtender) Missing(store Store) (res []string) {
	if _, ok := store.Get(e.id); ok {
		return nil
	}

	needed := e.gel.Missing(e.env)
	for _, n := range needed {
		if _, ok := store.Get(n); !ok {
			res = append(res, n)
		}
	}

	return res
}

func (e *gelExtender) Extend(store Store) error {
	if _, ok := store.Get(e.id); ok {
		return nil
	}

	if e.isFloat64Slice {
		return e.extendOnFloat64Slice(store)
	}

	return e.extend(store)
}

func (e *gelExtender) extend(store Store) error {
	env := e.env.Copy()

	needed := e.gel.Missing(env)
	for _, n := range needed {
		v, ok := store.Get(n)
		if !ok {
			return fmt.Errorf("Missing var %s", n)
		}
		env.AddVar(n, v)
	}

	data, err := e.gel.Eval(env)
	if err != nil {
		return err
	}

	store.Set(e.id, data)

	return nil
}

func (e *gelExtender) extendOnFloat64Slice(store Store) error {
	needed := e.gel.Missing(e.env)
	vars := make(map[string][]float64)
	l := 0
	for _, n := range needed {
		v, ok := store.Get(n)
		if !ok {
			return fmt.Errorf("Missing var %s", n)
		}
		vars[n] = v.([]float64)
		if l != 0 && len(vars[n]) != l {
			return fmt.Errorf("Different length of arrays in expression")
		}
		l = len(vars[n])
	}

	values := make([]float64, l)
	for i := range values {
		env := e.env.Copy()

		for _, n := range needed {
			env.AddVar(n, vars[n][i])
		}

		value, err := e.gel.Eval(env)
		if err != nil {
			return err
		}
		values[i] = value.(float64)
	}

	store.Set(e.id, values)

	return nil
}
