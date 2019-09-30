package f64s

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdev(t *testing.T) {
	assert.True(t, math.IsNaN(Stdev(Range(0, 0, 1))))
	assert.InDelta(t, 0.0, Stdev(Range(0.1, 1, 1)), 0.0001)
	assert.InDelta(t, 0.25, Stdev(Range(0.1, 1, 0.5)), 0.0001)
	assert.InDelta(t, 0.1556, Stdev(Range(0.11, 0.6, 0.11)), 0.0001)
}

func TestMean(t *testing.T) {
	assert.True(t, math.IsNaN(Mean(Range(0, 0, 1))))
	assert.InDelta(t, 0.1, Mean(Range(0.1, 1, 1)), 0.0001)
	assert.InDelta(t, 0.35, Mean(Range(0.1, 1, 0.5)), 0.0001)
	assert.InDelta(t, 0.33, Mean(Range(0.11, 0.6, 0.11)), 0.0001)
}
