package gel

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/Stromberg/gel/module"
	"github.com/Stromberg/gel/utils"
)

func init() {
	module.RegisterModules(StdLibModule)
}

var StdLibModule = &module.Module{
	Name: "stdlib",
	Funcs: []*module.Func{
		// &module.Func{Name: "strings.Join", F: utils.SimpleFunc(strings.Join, utils.CheckArity(2))},
		// &module.Func{Name: "strings.Split", F: utils.SimpleFunc(strings.Split, utils.CheckArity(2))},
		&module.Func{Name: "strings.Title", F: utils.SimpleFunc(strings.Title, utils.CheckArity(1)),
			Signature:   "(strings.Title cs)",
			Description: "Title cased string.",
		},
		&module.Func{Name: "strings.ToLower", F: utils.SimpleFunc(strings.ToLower, utils.CheckArity(1)),
			Signature:   "(strings.ToLower cs)",
			Description: "Lower cased string.",
		},
		&module.Func{Name: "strings.ToUpper", F: utils.SimpleFunc(strings.ToUpper, utils.CheckArity(1)),
			Signature:   "(strings.ToUpper cs)",
			Description: "Upper cased string.",
		},
		&module.Func{Name: "strings.TrimSpace", F: utils.SimpleFunc(strings.TrimSpace, utils.CheckArity(1)),
			Signature:   "(strings.TrimSpace cs)",
			Description: "Trim spaces from beginning and end of string.",
		},
		&module.Func{Name: "sprintf", F: utils.SimpleFunc(func(args ...interface{}) string {
			format := args[0].(string)
			return fmt.Sprintf(format, args[1:]...)
		}, utils.CheckArityAtLeast(1)),
			Signature:   "(sprintf fmt arg...)",
			Description: "Formatted string.",
		},
		&module.Func{Name: "math.Pow", F: utils.SimpleFunc(math.Pow, utils.CheckArity(2), utils.ParamToFloat64(0), utils.ParamToFloat64(1)),
			Signature:   "(math.Pow v p)",
			Description: "v^p.",
		},
		&module.Func{Name: "math.Sqrt", F: utils.SimpleFunc(math.Sqrt, utils.CheckArity(1), utils.ParamToFloat64(0)),
			Signature: "(math.Sqrt v)",
		},
		&module.Func{Name: "math.Ceil", F: utils.SimpleFunc(math.Ceil, utils.CheckArity(1), utils.ParamToFloat64(0)),
			Signature: "(math.Ceil v)",
		},
		&module.Func{Name: "math.Log", F: utils.SimpleFunc(math.Log, utils.CheckArity(1), utils.ParamToFloat64(0)),
			Signature: "(math.Log v)",
		},
		&module.Func{Name: "nan?", F: utils.SimpleFunc(math.IsNaN, utils.CheckArity(1), utils.ParamToFloat64(0))},
		&module.Func{Name: "pos-inf?", F: utils.SimpleFunc(func(v float64) bool { return math.IsInf(v, 0) }, utils.CheckArity(1), utils.ParamToFloat64(0))},
		&module.Func{Name: "combinations", F: combinationsFn,
			Signature:   "(combinations l...)",
			Description: "Takes lists as input and produces a list of lists of all combinations of those lists. Example: (combinations [1.0 2.0] [3.0]) => [[1.0 3.0] [2.0 3.0]]",
		},
		&module.Func{Name: "transpose", F: transposeFn,
			Signature:   "(transpose l)",
			Description: "Takes list of lists as input and produces a list of lists with all rows and columns transposed. \nExample: (transposed [1.0 2.0] [3.0 4.0]) => [[1.0 3.0] [2.0 4.0]]",
		},
		&module.Func{
			Name:      "in-range?",
			Signature: "((in-range? min max) v)",
			F:         inRangeFn},
	},
	LispFuncs: []*module.LispFunc{
		// &module.LispFunc{Name: "cap", F: "(func (lower upper) (func (x) (max lower (min upper x))))"},
		&module.LispFunc{Name: "pow", F: "(func [n] (func [x] (math.Pow x n)))",
			Signature:   "(pow v p)",
			Description: "v^p.",
		},
		&module.LispFunc{Name: "with-default", F: "(func [d] (func [x] (if (or (nan? x) (pos-inf? x)) d x)))",
			Signature:   "((with-default 3) v)",
			Description: "Returns a function that takes a value v that returns a default value if v is not a valid value",
		},
		&module.LispFunc{Name: "positive", F: "(func [d] (func [x] (if (or (nan? x) (pos-inf? x) (< x 0)) d x)))",
			Signature:   "((positive 3) v)",
			Description: "Returns a function that takes a value v that returns a default value if v is not a valid or positive value",
		},
		&module.LispFunc{Name: "str", F: "(func [n] (sprintf \"%v\" n))",
			Signature:   "(str v)",
			Description: "Converts v to a string representation",
		},
	},
	Scripts: []*module.Script{ // Mainly for test
		&module.Script{Name: "", Source: `
			(var cap (# (func [x] (max %1 (min %2 x)))))`,
		},
	},
}

var combinationsFn = utils.SimpleFunc(func(lists ...[]interface{}) interface{} {
	res := []interface{}{}

	cpy := func(src []interface{}, v interface{}) []interface{} {
		dst := make([]interface{}, len(src)+1)
		copy(dst, src)
		dst[len(src)] = v
		return dst
	}

	var impl func(base []interface{}, li int)
	impl = func(base []interface{}, li int) {
		if li == len(lists) {
			res = append(res, base)
			return
		}

		l := lists[li]
		for _, v := range l {
			impl(cpy(base, v), li+1)
		}
	}
	impl(nil, 0)
	return res
}, utils.CheckArityAtLeast(1))

var transposeFn = utils.ErrFunc(func(listOfLists []interface{}) (interface{}, error) {
	numLists := len(listOfLists)
	if numLists == 0 {
		return listOfLists, nil
	}
	list, ok := listOfLists[0].([]interface{})
	if !ok {
		return nil, errors.New("Expected list of lists")
	}
	listLen := len(list)

	res := make([]interface{}, listLen)
	for i := range res {
		res[i] = make([]interface{}, numLists)
		for j, v := range listOfLists {
			list, ok := v.([]interface{})
			if !ok {
				return nil, errors.New("Expected list of lists")
			}
			if len(list) != listLen {
				return nil, errors.New("All lists must be the same length")
			}
			res[i].([]interface{})[j] = list[i]
		}
	}

	return res, nil
}, utils.CheckArity(1))

var inRangeFn = utils.SimpleFunc(func(min, max float64) interface{} {
	return utils.SimpleFunc(func(v float64) interface{} {
		return v >= min && v <= max
	}, utils.CheckArity(1), utils.ParamToFloat64(0))
}, utils.CheckArity(2), utils.ParamToFloat64(0), utils.ParamToFloat64(1))
