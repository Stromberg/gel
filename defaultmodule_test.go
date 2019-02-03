package gel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultModule(t *testing.T) {
	testf := func(expr string, expected bool) {
		g, err := New(expr)
		assert.Nil(t, err)

		e := NewEnv()
		e.AddModule(DefaultModule)
		r, err := g.Eval(e)
		assert.Nil(t, err)
		assert.EqualValues(t, r, expected)
	}

	testf("(< 1 2)", true)
	testf("(< 1 1)", false)
	testf("(< 1 1.0)", false)
	testf("(< 1.0 1.0)", false)
	testf("(< 1.0 1)", false)

	testf("(> 1 2)", false)
	testf("(> 1 1)", false)
	testf("(> 1 1.0)", false)
	testf("(> 1.0 1.0)", false)
	testf("(> 1.0 1)", false)

	testf("(<= 1 2)", true)
	testf("(<= 2 1)", false)
	testf("(<= 1 1)", true)
	testf("(<= 1 1.0)", true)
	testf("(<= 1.0 1.0)", true)
	testf("(<= 1.0 1)", true)

	testf("(>= 1 2)", false)
	testf("(>= 2 1)", true)
	testf("(>= 1 1)", true)
	testf("(>= 1 1.0)", true)
	testf("(>= 1.0 1.0)", true)
	testf("(>= 1.0 1)", true)
}
