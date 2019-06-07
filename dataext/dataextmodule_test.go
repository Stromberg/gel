package dataext_test

import (
	"testing"

	"github.com/Stromberg/gel"
	"github.com/stretchr/testify/assert"
)

func TestDataExtModule(t *testing.T) {
	test := func(expr string, expected interface{}) {
		g, err := gel.New(expr)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		s, err := g.Eval(gel.NewEnv())
		assert.NoError(t, err)
		assert.Equal(t, expected, s)
	}

	test("((dataext.Fix 3.14) 3.5)", 3.5)
	test("((dataext.Fix 3.14) -3.5)", -3.5)
	test("((dataext.Fix 3.14) nan)", 3.14)
	test("((dataext.Fix 3.14) (vec 3.5 4.5))", []float64{3.5, 4.5})
	test("((dataext.Fix 3.14) (vec 3.5 -4.5 4.5))", []float64{3.5, -4.5, 4.5})
	test("((dataext.Fix 3.14) (vec 3.5 nan))", []float64{3.5, 3.14})

	test("((dataext.FixPos 3.14) 3.5)", 3.5)
	test("((dataext.FixPos 3.14) -3.5)", 3.14)
	test("((dataext.FixPos 3.14) nan)", 3.14)
	test("((dataext.FixPos 3.14) (vec 3.5 4.5))", []float64{3.5, 4.5})
	test("((dataext.FixPos 3.14) (vec 3.5 -4.5 4.5))", []float64{3.5, 3.14, 4.5})
	test("((dataext.FixPos 3.14) (vec 3.5 nan))", []float64{3.5, 3.14})
}
