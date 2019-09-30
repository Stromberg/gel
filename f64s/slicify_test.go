package f64s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlicify(t *testing.T) {
	f1 := func(v1 float64) float64 { return v1 }
	sf1 := Slicify(f1)
	assert.EqualValues(t, []float64{0, 1, 2}, sf1(Slice([]float64{0, 1, 2})).Slice())

	f2 := func(v1 int) float64 { return float64(v1) }
	sf2 := Slicify(f2)
	assert.EqualValues(t, []float64{0, 1, 2}, sf2(Slice([]float64{0, 1, 2})).Slice())

	f3 := func(v1 int64) float64 { return float64(v1) }
	sf3 := Slicify(f3)
	assert.EqualValues(t, []float64{0, 1, 2}, sf3(Slice([]float64{0, 1, 2})).Slice())

	f4 := func(v1, v2 float64) float64 { return v1 * v2 }
	sf4 := Slicify(f4)
	assert.EqualValues(t, []float64{46, 60, 238}, sf4(Slice([]float64{2, 5, 7}), Slice([]float64{23, 12, 34})).Slice())
}
