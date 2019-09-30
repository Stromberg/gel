package f64s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelChange(t *testing.T) {
	assert.Empty(t, RelChange(Range(0, 0, 1)))
	assert.EqualValues(t, []float64{0}, RelChange(Range(0.1, 1, 1)))
	assert.InDeltaSlice(t, []float64{0, 5}, RelChange(Range(0.1, 1, 0.5)), 0.0001)
	assert.InDeltaSlice(t, []float64{0, 1, 0.5, 0.3333, 0.25}, RelChange(Range(0.11, 0.6, 0.11)), 0.0001)
}

func TestAbsChange(t *testing.T) {
	assert.Empty(t, AbsChange(Range(0, 0, 1)))
	assert.EqualValues(t, []float64{0}, AbsChange(Range(0.1, 1, 1)))
	assert.InDeltaSlice(t, []float64{0, 0.5}, AbsChange(Range(0.1, 1, 0.5)), 0.0001)
	assert.InDeltaSlice(t, []float64{0, 0.11, 0.11, 0.11, 0.11}, AbsChange(Range(0.11, 0.6, 0.11)), 0.0001)
}

func TestAccumDev(t *testing.T) {
	assert.Empty(t, AccumDev(Range(0, 0, 1)))
	assert.EqualValues(t, []float64{0}, Pipe(RelChange, AccumDev)(Range(0.1, 1, 1)))
	assert.InDeltaSlice(t, []float64{0, 5}, Pipe(RelChange, AccumDev)(Range(0.1, 1, 0.5)), 0.0001)
	assert.InDeltaSlice(t, []float64{0, 1, 2, 3, 4}, Pipe(RelChange, AccumDev)(Range(0.11, 0.6, 0.11)), 0.0001)
}

func TestMomentum(t *testing.T) {
	assert.Empty(t, Momentum(1)(Range(0, 0, 1)))
	assert.Empty(t, Momentum(1)(Range(0.1, 1, 1)))
	assert.InDeltaSlice(t, []float64{5}, Momentum(1)(Range(0.1, 1, 0.5)), 0.0001)

	assert.InDeltaSlice(t, []float64{1, 0.5, 0.3333, 0.25}, Momentum(1)(Range(0.11, 0.6, 0.11)), 0.0001)
	assert.InDeltaSlice(t, []float64{2, 1, 0.6667}, Momentum(2)(Range(0.11, 0.6, 0.11)), 0.0001)
	assert.InDeltaSlice(t, []float64{3, 1.5}, Momentum(3)(Range(0.11, 0.6, 0.11)), 0.0001)
}

func TestCompositeMomentum(t *testing.T) {
	assert.Empty(t, CompositeMomentum(Range(0, 0, 1)))
	assert.InDeltaSlice(t, []float64{3.8518, 2.2057, 1.6042}, CompositeMomentum(Range(0.1, 1.6, 0.1)), 0.0001)
}

func TestShortMomentum(t *testing.T) {
	assert.Empty(t, ShortMomentum(Range(0, 0, 1)))
	assert.InDeltaSlice(t, []float64{18.141, 6.84, 4.0625, 2.858}, ShortMomentum(Range(0.1, 1.0, 0.1)), 0.001)
}

func TestSma(t *testing.T) {
	assert.Empty(t, Sma(1)(Range(0, 0, 1)))
	assert.InDeltaSlice(t, []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}, Pipe(Scale(0.1), Sma(1))(IntRange(1, 10, 1)), 0.0001)
	assert.InDeltaSlice(t, []float64{0.15, 0.25, 0.35, 0.45, 0.55, 0.65, 0.75, 0.85}, Pipe(Scale(0.1), Sma(2))(IntRange(1, 10, 1)), 0.0001)
	assert.InDeltaSlice(t, []float64{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}, Pipe(Scale(0.1), Sma(3))(IntRange(1, 10, 1)), 0.0001)
}

func TestAdjustedMean(t *testing.T) {
	assert.Empty(t, AdjustedMean(1)(Range(0, 0, 1)))
	assert.InDeltaSlice(t, []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}, Pipe(Scale(0.1), AdjustedMean(1))(IntRange(1, 10, 1)), 0.0001)
	assert.InDeltaSlice(t, []float64{0.1667, 0.2667, 0.3667, 0.4667, 0.5667, 0.6667, 0.7667, 0.8667}, Pipe(Scale(0.1), AdjustedMean(2))(IntRange(1, 10, 1)), 0.0001)
	assert.InDeltaSlice(t, []float64{0.2333, 0.3333, 0.4333, 0.5333, 0.6333, 0.7333, 0.8333}, Pipe(Scale(0.1), AdjustedMean(3))(IntRange(1, 10, 1)), 0.0001)
}

func TestStdevN(t *testing.T) {
	assert.Empty(t, StdevN(1)(Range(0, 0, 1)))
	assert.InDeltaSlice(t, []float64{0.0}, StdevN(1)(Range(0.1, 1, 1)), 0.0001)
	assert.InDeltaSlice(t, []float64{0, 0.0}, StdevN(1)(Range(0.1, 1, 0.5)), 0.0001)
	assert.InDeltaSlice(t, []float64{0.25}, StdevN(2)(Range(0.1, 1, 0.5)), 0.0001)
	assert.InDeltaSlice(t, []float64{0.0, 0.0, 0.0, 0.0, 0.0}, StdevN(1)(Range(0.11, 0.6, 0.11)), 0.0001)
	assert.InDeltaSlice(t, []float64{0.055, 0.055, 0.055, 0.055}, StdevN(2)(Range(0.11, 0.6, 0.11)), 0.0001)
	assert.InDeltaSlice(t, []float64{0.1556}, StdevN(5)(Range(0.11, 0.6, 0.11)), 0.0001)
}
