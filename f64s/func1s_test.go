package f64s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunc1s(t *testing.T) {
	assert.Empty(t, Add(10)(Range(0, 0, 1)))
	assert.EqualValues(t, []float64{10}, Add(10)(Range(0, 1, 1)))
	assert.EqualValues(t, []float64{10, 11}, Add(10)(Range(0, 1.1, 1)))
	assert.EqualValues(t, []float64{10, 10.11, 10.22, 10.33, 10.44, 10.55}, Add(10)(Range(0, 0.6, 0.11)))

	assert.InDeltaSlice(t, []float64{0, 0.33, 0.66, 0.99, 1.32, 1.65}, Scale(3)(Range(0, 0.6, 0.11)), 0.001)
	assert.EqualValues(t, []float64{0.25, 0.25, 0.25, 0.33, 0.44, 0.48}, Map(Cap(0.25, 0.48))(Range(0, 0.6, 0.11)))
	assert.InDeltaSlice(t, []float64{0, 0.0121, 0.0484, 0.1089, 0.1936, 0.3025}, Pow(2)(Range(0, 0.6, 0.11)), 0.001)

	assert.InDeltaSlice(t, []float64{4, 4.4521, 4.9284, 5.4289, 5.9536, 6.5025}, Pipe(Add(2), Pow(2))(Range(0, 0.6, 0.11)), 0.001)
}
