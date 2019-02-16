package gel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	g, err := New("")
	assert.Nil(t, err)
	assert.NotNil(t, g)

	g, err = New("(+ 1)")
	assert.Nil(t, err)
	assert.NotNil(t, g)

	g, err = New("(+ 1))")
	assert.NotNil(t, err)
	assert.Nil(t, g)
}
func TestMissing(t *testing.T) {
	g, err := New("")
	assert.Nil(t, err)

	e := NewEnv()
	assert.Zero(t, len(g.Missing(e)))

	g, err = New("(+ 1)")
	assert.Nil(t, err)
	assert.Zero(t, len(g.Missing(e)))

	g, err = New("(+ x)")
	assert.Nil(t, err)
	assert.EqualValues(t, g.Missing(e), []string{"x"})

	g, err = New("(+ x (f y))")
	assert.Nil(t, err)
	assert.EqualValues(t, g.Missing(e), []string{"x", "f", "y"})

	e.AddVar("x", 1)

	g, err = New("(+ x)")
	assert.Nil(t, err)
	assert.Zero(t, len(g.Missing(e)))

	g, err = New("(+ x (f y))")
	assert.Nil(t, err)
	assert.EqualValues(t, g.Missing(e), []string{"f", "y"})
}

func TestMissingVar(t *testing.T) {
	e := NewEnv()

	g, err := New("(var x 1.0) (+ x 3.0)")
	assert.Nil(t, err)
	assert.Zero(t, len(g.Missing(e)))

	r, err := g.Eval(e)
	assert.Nil(t, err)
	assert.Equal(t, 4.0, r)
}

func TestEval(t *testing.T) {
	g, err := New("")
	assert.Nil(t, err)

	e := NewEnv()

	r, err := g.Eval(e)
	assert.Nil(t, err)
	assert.Nil(t, r)

	g, err = New("(+ 1)")
	assert.Nil(t, err)
	r, err = g.Eval(e)
	assert.Nil(t, err)
	assert.EqualValues(t, r, 1)

	e.AddVar("x", 1.0)

	g, err = New("(+ x 2.14)")
	assert.Nil(t, err)
	r, err = g.Eval(e)
	assert.NoError(t, err)
	assert.EqualValues(t, r, 3.14)
}
