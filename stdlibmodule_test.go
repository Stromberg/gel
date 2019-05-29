package gel_test

import (
	"testing"

	"github.com/Stromberg/gel"
	"github.com/stretchr/testify/assert"
)

func TestStdLibModuleStrings(t *testing.T) {
	test := func(expr string, expected interface{}) {
		g, err := gel.New(expr, gel.StdLibModule)
		assert.Nil(t, err)
		assert.NotNil(t, g)
		s, err := g.Eval(gel.NewEnv())
		assert.Nil(t, err)
		assert.Equal(t, expected, s)
	}

	test("(strings.Title \"world\")", "World")
	test("(strings.ToUpper \"world\")", "WORLD")
	test("(strings.ToLower \"worLd\")", "world")
	test("(strings.TrimSpace \"  worLd  \")", "worLd")
	test("(printf \"Grr\\n\")", 4)
	test("(printf \"Grr: %v\\n\" 3.14)", 10)
	test("(sprintf \"Grr\\n\")", "Grr\n")
	test("(sprintf \"Grr: %v\\n\" 3.14)", "Grr: 3.14\n")
}
