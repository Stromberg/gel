package gel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllFunctionNames(t *testing.T) {
	names := AllFunctionNames()
	assert.NotEqual(t, 0, len(names))
}

func TestMatchingFuncNames(t *testing.T) {
	names := MatchingFuncNames("math.*")
	assert.Equal(t, 4, len(names))
}

func TestFunctionRepr(t *testing.T) {
	repr := FunctionRepr("sort-asc")
	assert.NotNil(t, repr)
}
