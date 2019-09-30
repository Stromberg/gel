package f64s_test

import (
	"testing"

	"github.com/Stromberg/gel/f64s"
	"github.com/stretchr/testify/assert"
)

func TestNrank(t *testing.T) {

	assert.EqualValues(t, []float64(nil), f64s.Nrank([]float64{}))
	assert.EqualValues(t, []float64{1}, f64s.Nrank([]float64{10}))
	assert.EqualValues(t, []float64{0, 0.5, 1.0}, f64s.Nrank([]float64{1, 5, 10}))
	assert.EqualValues(t, []float64{0, 0.25, 0.5, 0.75, 1.0}, f64s.Nrank([]float64{1, 5, 10, 15, 20}))
	assert.EqualValues(t, []float64{0, 1.0, 0.75, 0.5, 0.25}, f64s.Nrank([]float64{1, 20, 15, 10, 5}))
}
