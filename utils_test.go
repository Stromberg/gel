package gel_test

import (
	"testing"

	"github.com/Stromberg/gel"
	"github.com/stretchr/testify/assert"
)

func TestIsAnyStrings(t *testing.T) {
	assert.True(t, gel.IsAnyStrings("12"))
	assert.True(t, gel.IsAnyStrings(nil, "12"))
	assert.True(t, gel.IsAnyStrings(12, "12"))
	assert.False(t, gel.IsAnyStrings(12))
	assert.False(t, gel.IsAnyStrings(nil, 12))
}

func TestIsAllStrings(t *testing.T) {
	assert.True(t, gel.IsAllStrings("12"))
	assert.True(t, gel.IsAllStrings("nil", "12"))
	assert.False(t, gel.IsAllStrings(12, "12"))
	assert.False(t, gel.IsAllStrings(12))
	assert.False(t, gel.IsAllStrings(nil, 12))
}

func TestIsAnyFloats(t *testing.T) {
	assert.True(t, gel.IsAnyFloats(12.0))
	assert.False(t, gel.IsAnyFloats("12"))
	assert.False(t, gel.IsAnyFloats("nil", "12"))
	assert.True(t, gel.IsAnyFloats(12.0, "12"))
	assert.False(t, gel.IsAnyFloats(12, "12"))
	assert.False(t, gel.IsAnyFloats(12))
	assert.False(t, gel.IsAnyFloats(nil, 12))
	assert.True(t, gel.IsAnyFloats(nil, 12.0))
}

func TestIsAnySlice(t *testing.T) {
	assert.True(t, gel.IsAnySlice([]interface{}{}))
	assert.True(t, gel.IsAnySlice([]float64{}))
	assert.True(t, gel.IsAnySlice([]interface{}{"D"}))
	assert.True(t, gel.IsAnySlice([]float64{12.0}))
	assert.False(t, gel.IsAnySlice(12.0))
	assert.False(t, gel.IsAnySlice("nil", "12"))
	assert.True(t, gel.IsAnySlice([]float64{}, "12"))
	assert.True(t, gel.IsAnySlice("12", []float64{}))
}

func TestMakeAllEitherSliceOrValue(t *testing.T) {
	v, err := gel.MakeAllEitherSliceOrValue([]float64{})
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{[]float64{}}, v)

	v, err = gel.MakeAllEitherSliceOrValue([]float64{12.0})
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{[]float64{12.0}}, v)

	v, err = gel.MakeAllEitherSliceOrValue([]float64{12.0}, 12.0)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{[]float64{12.0}, []float64{12.0}}, v)

	v, err = gel.MakeAllEitherSliceOrValue([]float64{12.0}, int64(12))
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{[]float64{12.0}, []float64{12.0}}, v)
}
