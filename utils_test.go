package gel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleFunc(t *testing.T) {
	f := SimpleFunc(func(v1, v2 float64) bool { return v1 <= v2 })

	v, err := f([]interface{}{23.0, 12.0})
	assert.NoError(t, err)
	assert.False(t, v.(bool))

	v, err = f([]interface{}{12.0, 13.0})
	assert.NoError(t, err)
	assert.True(t, v.(bool))
}
