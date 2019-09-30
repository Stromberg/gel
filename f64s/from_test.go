package f64s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	assert.Empty(t, Repeat(12, 0))
	assert.EqualValues(t, []float64{12}, Repeat(12, 1))
	assert.EqualValues(t, []float64{12, 12, 12, 12}, Repeat(12, 4))
}

func TestRange(t *testing.T) {
	assert.Nil(t, Range(0, 0, 1))
	assert.EqualValues(t, []float64{0}, Range(0, 1, 1))
	assert.EqualValues(t, []float64{0, 1}, Range(0, 1.1, 1))
	assert.EqualValues(t, []float64{1}, Range(1, 0, -1))
	assert.EqualValues(t, []float64{1, 0}, Range(1, -0.1, -1))
	assert.EqualValues(t, []float64{0, 0.11, 0.22, 0.33, 0.44, 0.55, 0.66, 0.77, 0.88, 0.99}, Range(0, 1, 0.11))
	assert.InDeltaSlice(t, []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0}, Range(0.1, 1, 0.1), 0.001)
}

func TestIntRange(t *testing.T) {
	assert.Nil(t, IntRange(0, 0, 1))
	assert.EqualValues(t, []float64{0}, IntRange(0, 1, 1))
	assert.EqualValues(t, []float64{0, 10}, IntRange(0, 11, 10))
	assert.EqualValues(t, []float64{1}, IntRange(1, 0, -1))
	assert.EqualValues(t, []float64{10, 0}, IntRange(10, -1, -10))
	assert.EqualValues(t, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, IntRange(1, 10, 1))
}
