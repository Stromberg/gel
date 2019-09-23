package utils_test

import (
	"testing"

	"github.com/Stromberg/gel/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsAnyStrings(t *testing.T) {
	assert.True(t, utils.IsAnyStrings("12"))
	assert.True(t, utils.IsAnyStrings(nil, "12"))
	assert.True(t, utils.IsAnyStrings(12, "12"))
	assert.False(t, utils.IsAnyStrings(12))
	assert.False(t, utils.IsAnyStrings(nil, 12))
}

func TestIsAllStrings(t *testing.T) {
	assert.True(t, utils.IsAllStrings("12"))
	assert.True(t, utils.IsAllStrings("nil", "12"))
	assert.False(t, utils.IsAllStrings(12, "12"))
	assert.False(t, utils.IsAllStrings(12))
	assert.False(t, utils.IsAllStrings(nil, 12))
}

func TestIsAnyFloats(t *testing.T) {
	assert.True(t, utils.IsAnyFloats(12.0))
	assert.False(t, utils.IsAnyFloats("12"))
	assert.False(t, utils.IsAnyFloats("nil", "12"))
	assert.True(t, utils.IsAnyFloats(12.0, "12"))
	assert.False(t, utils.IsAnyFloats(12, "12"))
	assert.False(t, utils.IsAnyFloats(12))
	assert.False(t, utils.IsAnyFloats(nil, 12))
	assert.True(t, utils.IsAnyFloats(nil, 12.0))
}

func TestIsAnySlice(t *testing.T) {
	assert.True(t, utils.IsAnySlice([]interface{}{}))
	assert.True(t, utils.IsAnySlice([]float64{}))
	assert.True(t, utils.IsAnySlice([]interface{}{"D"}))
	assert.True(t, utils.IsAnySlice([]float64{12.0}))
	assert.False(t, utils.IsAnySlice(12.0))
	assert.False(t, utils.IsAnySlice("nil", "12"))
	assert.True(t, utils.IsAnySlice([]float64{}, "12"))
	assert.True(t, utils.IsAnySlice("12", []float64{}))
}

func TestMakeAllEitherSliceOrValue(t *testing.T) {
	v, err := utils.MakeAllEitherSliceOrValue([]float64{})
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{[]float64{}}, v)

	v, err = utils.MakeAllEitherSliceOrValue([]float64{12.0})
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{[]float64{12.0}}, v)

	v, err = utils.MakeAllEitherSliceOrValue([]float64{12.0}, 12.0)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{[]float64{12.0}, []float64{12.0}}, v)

	v, err = utils.MakeAllEitherSliceOrValue([]float64{12.0}, int64(12))
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{[]float64{12.0}, []float64{12.0}}, v)
}
