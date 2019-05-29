package dataext

import (
	"fmt"
	"testing"

	"github.com/Stromberg/gel"

	"github.com/stretchr/testify/assert"
)

type store struct {
	Store

	data map[string]interface{}
}

func newStore() *store {
	return &store{data: make(map[string]interface{})}
}

func (s *store) Get(id string) (value interface{}, ok bool) {
	value, ok = s.data[id]
	return
}

func (s *store) Set(id string, value interface{}) {
	s.data[id] = value
}

func mulSlice(s1, s2 []float64) []float64 {
	res := make([]float64, len(s1))
	for i := range res {
		res[i] = s1[i] * s2[i]
	}
	return res
}

func mulSliceWrap(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Need 2 arguments")
	}

	if s1, ok := args[0].([]float64); ok {
		if s2, ok := args[1].([]float64); ok {
			res := mulSlice(s1, s2)
			return res, nil
		}
	}

	return nil, fmt.Errorf("Wrong argument types")
}

///

func TestExtendExpressionSimple(t *testing.T) {
	store := newStore()
	store.Set("v1", []float64{1, 2, 3, 4})
	store.Set("v2", []float64{13, 12, 11, 10})

	env := gel.NewEnv()
	env.AddVar("*", mulSliceWrap)

	e := NewGel("v3", "(* v1 v2)", env)
	assert.NotNil(t, e)

	missing := e.Missing(store)
	assert.Equal(t, 0, len(missing))

	err := Extend(store, e)
	assert.NoError(t, err)
	res, ok := store.Get("v3")
	assert.True(t, ok)
	assert.Equal(t, []float64{13, 24, 33, 40}, res)
}

func TestExtendExpression(t *testing.T) {
	store := newStore()
	store.Set("v1", []float64{1, 2, 3, 4})
	store.Set("v2", []float64{13, 12, 11, 10})

	env := gel.NewEnv()
	env.AddVar("*", mulSliceWrap)

	e1 := NewGel("v5", "(* (* v1 v1) v4)", env)
	e2 := NewGel("v3", "(* v1 v2)", env)
	e3 := NewGel("v4", "(* (* v1 v2) v3)", env)

	err := Extend(store, e1, e2, e3)
	assert.NoError(t, err)

	res, ok := store.Get("v3")
	assert.True(t, ok)
	assert.Equal(t, []float64{13, 24, 33, 40}, res)

	res, ok = store.Get("v4")
	assert.True(t, ok)
	assert.Equal(t, []float64{169, 576, 1089, 1600}, res)

	res, ok = store.Get("v5")
	assert.True(t, ok)
	assert.Equal(t, []float64{169, 2304, 9801, 25600}, res)
}

func TestExtendExpressionVar(t *testing.T) {
	store := newStore()
	store.Set("v1", []float64{1, 2, 3, 4})
	store.Set("v2", []float64{13, 12, 11, 10})

	env := gel.NewEnv()
	env.AddVar("*", mulSliceWrap)

	e := NewGel("v3", "(var x 1.0) (* v1 v2)", env)
	assert.NotNil(t, e)

	missing := e.Missing(store)
	assert.Equal(t, 0, len(missing))

	err := Extend(store, e)
	assert.NoError(t, err)
	res, ok := store.Get("v3")
	assert.True(t, ok)
	assert.Equal(t, []float64{13, 24, 33, 40}, res)
}

func TestExtendPrimitiveExpression(t *testing.T) {
	store := newStore()
	store.Set("v1", []float64{1, 2, 3, 4})
	store.Set("v2", []float64{13, 12, 11, 10})

	env := gel.NewEnv()

	e1 := NewGelOnFloat64Slice("v5", "(* v1 v1 v4)", env)
	e2 := NewGelOnFloat64Slice("v3", "(* v1 v2)", env)
	e3 := NewGelOnFloat64Slice("v4", "(* v1 v2 v3)", env)

	err := Extend(store, e1, e2, e3)
	assert.NoError(t, err)

	res, ok := store.Get("v3")
	assert.True(t, ok)
	assert.Equal(t, []float64{13, 24, 33, 40}, res)

	res, ok = store.Get("v4")
	assert.True(t, ok)
	assert.Equal(t, []float64{169, 576, 1089, 1600}, res)

	res, ok = store.Get("v5")
	assert.True(t, ok)
	assert.Equal(t, []float64{169, 2304, 9801, 25600}, res)
}

func TestExtendFuncSimple(t *testing.T) {
	store := newStore()
	store.Set("v1", []float64{1, 2, 3, 4})
	store.Set("v2", []float64{13, 12, 11, 10})

	e := NewFunc("v3", []string{"v1", "v2"}, func(store Store) (interface{}, error) {
		v1, _ := store.Get("v1")
		v2, _ := store.Get("v2")
		return mulSlice(v1.([]float64), v2.([]float64)), nil
	})

	missing := e.Missing(store)
	assert.Equal(t, 0, len(missing))

	err := Extend(store, e)
	assert.NoError(t, err)

	res, ok := store.Get("v3")
	assert.True(t, ok)
	assert.Equal(t, []float64{13, 24, 33, 40}, res)
}
