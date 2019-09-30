package f64s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTake(t *testing.T) {
	assert.Nil(t, Take(0)(Range(0, 0, 1)))
	assert.Nil(t, Take(0)(Range(0, 1, 1)))
	assert.Nil(t, Take(0)(Range(0, 1.1, 1)))

	assert.EqualValues(t, []float64{0}, Take(1)(Range(0, 1, 0.11)))
	assert.EqualValues(t, []float64{0, 0.11}, Take(2)(Range(0, 1, 0.11)))
	assert.EqualValues(t, []float64{0, 0.11, 0.22}, Take(3)(Range(0, 1, 0.11)))
}
