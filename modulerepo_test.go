package gel

import (
	"testing"

	"github.com/Stromberg/gel/module"
	"github.com/stretchr/testify/assert"
)

func TestAllFunctionNames(t *testing.T) {
	names := module.AllFunctionNames()
	assert.NotEqual(t, 0, len(names))
}

func TestMatchingFuncNames(t *testing.T) {
	names := module.MatchingFuncNames("math.*")
	assert.Equal(t, 5, len(names))
}

func TestFunctionRepr(t *testing.T) {
	repr := module.FunctionRepr("sort-asc")
	assert.NotNil(t, repr)
}
