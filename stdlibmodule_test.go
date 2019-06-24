package gel_test

import (
	"testing"

	"github.com/Stromberg/gel"
	"github.com/stretchr/testify/assert"
)

func TestStdLibModuleStrings(t *testing.T) {
	test := func(expr string, expected interface{}) {
		g, err := gel.New(expr)
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
	test("(str 3.14)", "3.14")
}

func TestStdLibModuleMath(t *testing.T) {
	test := func(expr string, expected interface{}) {
		g, err := gel.New(expr)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		s, err := g.Eval(gel.NewEnv())
		assert.NoError(t, err)
		assert.Equal(t, expected, s)
	}

	test("((cap 3 4) 3.5)", 3.5)
	test("((cap 3 4) 10.5)", float64(4))
	test("((cap 3 4) 1)", int64(3))

	test("(math.Pow 3 2)", float64(9))
	test("(math.Pow 2 3)", float64(8))

	test("(math.Sqrt 4)", float64(2))

	test("((pow 3) 2)", float64(8))
	test("((pow 2) 3)", float64(9))

	test("(nan? 3)", false)
	test("(nan? nan)", true)
	test("(pos-inf? 3)", false)
	test("((with-default 0.0) 3.0)", 3.0)
	test("((with-default 0.0) nan)", 0.0)
	test("((positive 1.0) 3.0)", 3.0)
	test("((positive 1.0) -1.0)", 1.0)
	test("((positive 1.0) nan)", 1.0)
}

func TestStdLibModuleCombinations(t *testing.T) {
	test := func(expr string, expected interface{}) {
		g, err := gel.New(expr)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		s, err := g.Eval(gel.NewEnv())
		assert.NoError(t, err)
		assert.Equal(t, expected, s)
	}

	test(
		"(combinations (list))",
		[]interface{}{},
	)
	test(
		"(combinations (list 1.0))",
		[]interface{}{
			[]interface{}{1.0},
		},
	)
	test(
		"(combinations (list 1.0 2.0))",
		[]interface{}{
			[]interface{}{1.0},
			[]interface{}{2.0},
		},
	)
	test(
		"(combinations (list 1.0 2.0) (list 3.0))",
		[]interface{}{
			[]interface{}{1.0, 3.0},
			[]interface{}{2.0, 3.0},
		},
	)
	test(
		"(combinations (list 1.0 2.0) (list 3.0 4.0) (list 5.0 6.0))",
		[]interface{}{
			[]interface{}{1.0, 3.0, 5.0},
			[]interface{}{1.0, 3.0, 6.0},
			[]interface{}{1.0, 4.0, 5.0},
			[]interface{}{1.0, 4.0, 6.0},
			[]interface{}{2.0, 3.0, 5.0},
			[]interface{}{2.0, 3.0, 6.0},
			[]interface{}{2.0, 4.0, 5.0},
			[]interface{}{2.0, 4.0, 6.0},
		},
	)
}

func TestStdLibModuleTranspose(t *testing.T) {
	test := func(expr string, expected interface{}) {
		g, err := gel.New(expr)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		s, err := g.Eval(gel.NewEnv())
		assert.NoError(t, err)
		assert.Equal(t, expected, s)
	}

	test(
		"(transpose (list))",
		[]interface{}{},
	)
	test(
		"(transpose (list (list 1.0)))",
		[]interface{}{
			[]interface{}{1.0},
		},
	)
	test(
		"(transpose (list (list 1.0 2.0) (list 3.0 4.0)))",
		[]interface{}{
			[]interface{}{1.0, 3.0},
			[]interface{}{2.0, 4.0},
		},
	)
	test(
		"(transpose (list (list 1.0 2.0) (list 3.0 4.0) (list 5.0 6.0)))",
		[]interface{}{
			[]interface{}{1.0, 3.0, 5.0},
			[]interface{}{2.0, 4.0, 6.0},
		},
	)
}
