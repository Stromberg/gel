package f64s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	assert.Nil(t, Skip(0)(Range(0, 0, 1)))
	assert.EqualValues(t, []float64{0}, Skip(0)(Range(0, 1, 1)))
	assert.EqualValues(t, []float64{0, 1}, Skip(0)(Range(0, 1.1, 1)))
	assert.EqualValues(t, []float64{1}, Skip(0)(Range(1, 0, -1)))
	assert.EqualValues(t, []float64{1, 0}, Skip(0)(Range(1, -0.1, -1)))
	assert.EqualValues(t, []float64{0, 0.11, 0.22, 0.33, 0.44, 0.55, 0.66, 0.77, 0.88, 0.99}, Skip(0)(Range(0, 1, 0.11)))

	assert.Nil(t, Skip(1)(Range(0, 0, 1)))
	assert.Nil(t, Skip(1)(Range(0, 1, 1)))
	assert.EqualValues(t, []float64{1}, Skip(1)(Range(0, 1.1, 1)))
	assert.EqualValues(t, []float64{0.11, 0.22, 0.33, 0.44, 0.55, 0.66, 0.77, 0.88, 0.99}, Skip(1)(Range(0, 1, 0.11)))

	assert.EqualValues(t, []float64{0.44, 0.55, 0.66, 0.77, 0.88, 0.99}, Skip(4)(Range(0, 1, 0.11)))
}
