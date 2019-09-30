package f64s_test

import (
	"testing"

	"github.com/Stromberg/gel"
	"github.com/Stromberg/gel/f64s"
	"github.com/Stromberg/gel/utils"
	"github.com/stretchr/testify/assert"
)

func TestF64sModule(t *testing.T) {
	test := func(expr string, expected interface{}) {
		g, err := gel.New(expr)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		s, err := g.Eval(gel.NewEnv())
		assert.NoError(t, err)
		if utils.IsSlice(s) {
			assert.InDeltaSlice(t, expected, s, 0.001)
		} else {
			assert.InDelta(t, expected, s, 0.0001)
		}
	}

	test("(f64s/RelChange (vec 1 1.5))", []float64{0, 0.5})
	test("(f64s/AbsChange (vec 1.5 2.5))", []float64{0, 1.0})

	test("((f64s/RelChangeN 1) (vec 1 1.5))", []float64{0.5})
	test("((f64s/RelChangeN 2) (vec 1 1.5))", []float64{})
	test("((f64s/RelChangeN 2) (vec 1 1.5 2.0))", []float64{1.0})
	test("((f64s/RelChangeN 2) (vec 1 1.5 2.0 4.5))", []float64{1.0, 2.0})
	test("((f64s/AbsChangeN 2) (vec 1 1.5 2.0 4.5))", []float64{1.0, 3.0})

	//

	test("(f64s/AccumDev (f64s/RelChange (vec-range 0.1 1 0.5)))", []float64{0, 5.0})

	test("((f64s/Momentum 3) (vec-range 0.11 0.6 0.11))", []float64{3, 1.5})

	test("((f64s/Sma 3) (vec-range 1 10 1))", f64s.IntRange(2, 9, 1))

	test("((f64s/Pma 0.1) (vec-range 1 15 1))", []float64{12, 12.25, 12.75})

	test("((f64s/StdevN 2) (vec-range 0.1 1 0.5))", []float64{0.25})

	test("(f64s/CompositeMomentum (vec-range 0.1 1.6 0.1))", []float64{3.8518, 2.2057, 1.6042})

	test("(f64s/ShortMomentum (vec-range 0.1 1.0 0.1))", []float64{18.141, 6.84, 4.0625, 2.858})

	test("((f64s/MaxN 3) (vec-range 1 10 2))", f64s.IntRange(1, 10, 2))

	test("(f64s/Stdev (vec-range 0.11 0.6 0.11))", 0.1556)

	test("(f64s/Mean (vec-range 0.11 0.6 0.11))", 0.33)

	test("(f64s/Sum (vec 0 1 1.5 2.0))", 4.5)

	test("(f64s/GeometricMeanDev (vec 0.2 0.3 0.4))", 0.2974)

	test("(f64s/Nrank (vec 1 20 15 10 5))", []float64{0, 1.0, 0.75, 0.5, 0.25})
}
