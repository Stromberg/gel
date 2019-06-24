package gel

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

func init() {
	RegisterModules(StdLibModule)
}

var StdLibModule = &Module{
	Name: "stdlib",
	Funcs: []*Func{
		// &Func{Name: "strings.Join", F: SimpleFunc(strings.Join, CheckArity(2))},
		// &Func{Name: "strings.Split", F: SimpleFunc(strings.Split, CheckArity(2))},
		&Func{Name: "strings.Title", F: SimpleFunc(strings.Title, CheckArity(1))},
		&Func{Name: "strings.ToLower", F: SimpleFunc(strings.ToLower, CheckArity(1))},
		&Func{Name: "strings.ToUpper", F: SimpleFunc(strings.ToUpper, CheckArity(1))},
		&Func{Name: "strings.TrimSpace", F: SimpleFunc(strings.TrimSpace, CheckArity(1))},
		&Func{Name: "printf", F: ErrFunc(func(args ...interface{}) (int, error) {
			format := args[0].(string)
			return fmt.Printf(format, args[1:]...)
		}, CheckArityAtLeast(1))},
		&Func{Name: "sprintf", F: SimpleFunc(func(args ...interface{}) string {
			format := args[0].(string)
			return fmt.Sprintf(format, args[1:]...)
		}, CheckArityAtLeast(1))},
		&Func{Name: "math.Pow", F: SimpleFunc(math.Pow, CheckArity(2), ParamToFloat64(0), ParamToFloat64(1))},
		&Func{Name: "math.Sqrt", F: SimpleFunc(math.Sqrt, CheckArity(1), ParamToFloat64(0))},
		&Func{Name: "nan?", F: SimpleFunc(math.IsNaN, CheckArity(1), ParamToFloat64(0))},
		&Func{Name: "pos-inf?", F: SimpleFunc(func(v float64) bool { return math.IsInf(v, 0) }, CheckArity(1), ParamToFloat64(0))},
		&Func{Name: "combinations", F: combinationsFn},
		&Func{Name: "transpose", F: transposeFn},
	},
	LispFuncs: []*LispFunc{
		&LispFunc{Name: "cap", F: "(func (lower upper) (func (x) (max lower (min upper x))))"},
		&LispFunc{Name: "pow", F: "(func (n) (func (x) (math.Pow x n)))"},
		&LispFunc{Name: "with-default", F: "(func (d) (func (x) (if (or (nan? x) (pos-inf? x)) d x)))"},
		&LispFunc{Name: "positive", F: "(func (d) (func (x) (if (or (nan? x) (pos-inf? x) (< x 0)) d x)))"},
		&LispFunc{Name: "str", F: "(func (n) (sprintf \"%v\" n))"},
	},
}

var combinationsFn = SimpleFunc(func(lists ...[]interface{}) interface{} {
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
}, CheckArityAtLeast(1))

var transposeFn = ErrFunc(func(listOfLists []interface{}) (interface{}, error) {
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
}, CheckArity(1))
