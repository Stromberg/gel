package gel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"path"
	"sort"
	"time"

	"github.com/Stromberg/gel/ast"
	"github.com/Stromberg/gel/dataserie"
	"github.com/Stromberg/gel/module"
	"github.com/Stromberg/gel/utils"
	"github.com/google/uuid"
)

func init() {
	module.RegisterModules(GlobalsModule)
	module.RegisterModules(dataserie.Module)
}

var GlobalsModule = &module.Module{
	Name: "globals",
	Funcs: []*module.Func{
		&module.Func{Name: "eval", F: evalFn},
		&module.Func{Name: "eval-file", F: evalFileFn},
		&module.Func{Name: "load", F: loadFn},
		&module.Func{Name: "load-file", F: loadFileFn},
		&module.Func{Name: "slurp", F: slurpFn},
		&module.Func{Name: "true", F: true},
		&module.Func{Name: "false", F: false},
		&module.Func{Name: "nil", F: nil},
		&module.Func{Name: "nan", F: math.NaN()},
		&module.Func{Name: "error", F: errorFn},
		&module.Func{Name: "==", F: eqFn},
		&module.Func{Name: "<", F: lessThanFn},
		&module.Func{Name: ">", F: greaterThanFn},
		&module.Func{Name: "<=", F: lessThanEqualFn},
		&module.Func{Name: ">=", F: greaterThanEqualFn},
		&module.Func{Name: "!=", F: neFn},
		&module.Func{Name: "+", F: plusFn},
		&module.Func{Name: "-", F: minusFn},
		&module.Func{Name: "*", F: mulFn},
		&module.Func{Name: "/", F: divFn},
		&module.Func{Name: "%", F: modFn},
		&module.Func{Name: "!", F: notFn},
		&module.Func{Name: "not", F: notFn},
		&module.Func{Name: "int", F: intFn},
		&module.Func{Name: "float", F: floatFn},
		&module.Func{Name: "min", F: minFn},
		&module.Func{Name: "max", F: maxFn},
		&module.Func{Name: "or", F: orFn},
		&module.Func{Name: "and", F: andFn},
		&module.Func{Name: "if", F: ifFn},
		&module.Func{Name: "cond", F: condFn},
		&module.Func{Name: "var", F: varFn},
		&module.Func{Name: "set", F: setFn},
		&module.Func{Name: "do", F: doFn},
		&module.Func{Name: "code", F: codeFn},
		&module.Func{Name: "func", F: funcFn},
		&module.Func{Name: "fn", F: funcFn},
		&module.Func{Name: "#", F: macroFn},
		&module.Func{Name: "for", F: forFn},
		&module.Func{Name: "vec", F: vecFn},
		&module.Func{Name: "vec2list", F: vecToListFn},
		&module.Func{Name: "list2vec", F: listToVecFn},
		&module.Func{Name: "list", F: utils.NewList},
		&module.Func{Name: "vec?", F: isVecFn},
		&module.Func{Name: "list?", F: isListFn},
		&module.Func{Name: "dict", F: utils.NewDict},
		&module.Func{Name: "dict?", F: isDictFn},
		&module.Func{Name: "dict-keys", F: dictKeysFn},
		&module.Func{Name: "get", F: utils.GetFn},
		&module.Func{Name: "sub", F: subFn},
		&module.Func{Name: "contains?", F: containsFn},
		&module.Func{Name: "update!", F: updateFn},
		&module.Func{Name: "len", F: lenFn},
		&module.Func{Name: "append", F: appendFn},
		&module.Func{Name: "concat", F: concatFn},
		&module.Func{Name: "merge", F: mergeFn},
		&module.Func{Name: "range", F: rangeFn},
		&module.Func{Name: "vec-range", F: vecRangeFn},
		&module.Func{Name: "repeat", F: repeatFn},
		&module.Func{Name: "reverse", F: reverseFn},
		&module.Func{Name: "vec-repeat", F: vecRepeatFn},
		&module.Func{Name: "map", F: mapFn},
		&module.Func{Name: "map-indexed", F: mapIndexedFn},
		&module.Func{Name: "vec-map", F: vecMapFn},
		&module.Func{Name: "vec-map-indexed", F: vecMapIndexedFn},
		&module.Func{Name: "apply", F: applyFn},
		&module.Func{Name: "vec-apply", F: vecApplyFn},
		&module.Func{Name: "vec-rand", F: vecRandFn},
		&module.Func{Name: "list-rand", F: listRandFn},
		&module.Func{Name: "reduce", F: reduceFn},
		&module.Func{Name: "filter", F: filterFn},
		&module.Func{Name: "count-if", F: countIfFn},
		&module.Func{Name: "flatten", F: flattenFn},
		&module.Func{Name: "skip", F: skipFn},
		&module.Func{Name: "take", F: takeFn},
		&module.Func{Name: "sort-asc", F: sortAscFn},
		&module.Func{Name: "sort-desc", F: sortDescFn},
		&module.Func{Name: "sortindex", F: sortIndexFn},
		&module.Func{Name: "bind", F: bindFn},
		&module.Func{Name: "json", F: jsonFn},
		&module.Func{Name: "uuid", F: uuidFn},
		&module.Func{Name: "rand", F: randFn()},
		&module.Func{Name: "repeatedly", F: repeatedlyFn},
		&module.Func{Name: "time", F: timeFn},
		&module.Func{Name: "printf", F: printfFn},
		&module.Func{Name: "docs", F: docsFn},
	},
	LispFuncs: []*module.LispFunc{
		&module.LispFunc{Name: "identity", F: "(func [x] x)"},
		&module.LispFunc{Name: "empty?", F: "(func [x] (== (len x) 0))"},
		&module.LispFunc{Name: "first", F: "(func [s] (get s 0))"},
		&module.LispFunc{Name: "second", F: "(func [s] (get s 1))"},
		&module.LispFunc{Name: "rest", F: "(func [s] (skip 1 s))"},
		&module.LispFunc{Name: "last", F: "(func [s] (if (empty? s) nil (get s (- (len s) 1))))"},
		&module.LispFunc{Name: "inc", F: "(func [s] (+ s 1))"},
		&module.LispFunc{Name: "dec", F: "(func [s] (- s 1))"},
		&module.LispFunc{Name: "def", F: "var"},
	},
}

func errorFn(args ...interface{}) (value interface{}, err error) {
	if len(args) == 1 {
		if s, ok := args[0].(string); ok {
			return nil, errors.New(s)
		}
	}
	return nil, errors.New("error function takes a single string argument")
}

func printfFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("printf function requires at least a string argument")
	}

	formatRaw, err := scope.Eval(args[0])
	if err != nil {
		return nil, err
	}
	format, ok := formatRaw.(string)
	if !ok {
		return nil, errors.New("printf requires a string argument")
	}

	vargs := make([]interface{}, len(args[1:]))
	for i, arg := range args[1:] {
		varg, err := scope.Eval(arg)
		if err != nil {
			return nil, err
		}
		vargs[i] = varg
	}

	err = scope.Printf(format, vargs...)
	return nil, err
}

var evalFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	code, ok := args[0].(string)
	if !ok {
		return nil, utils.ErrParameterType
	}

	g, err := New(code)
	if err != nil {
		return nil, fmt.Errorf("Error in eval: %v", err)
	}

	return g.Eval(NewEnv())
}, utils.CheckArity(1))

func loadFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, errors.New("load function takes a single string argument")
	}

	r, err := scope.Eval(args[0])
	if err != nil {
		return nil, err
	}
	code, ok := r.(string)
	if !ok {
		return nil, utils.ErrParameterType
	}

	node, err := ParseString(scope.fset, "", code)
	if err != nil {
		return nil, err
	}

	return scope.Eval(node)
}

func loadFileFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, errors.New("load-file function takes a single string argument")
	}

	r, err := scope.Eval(args[0])
	if err != nil {
		return nil, err
	}

	file, ok := r.(string)
	if !ok {
		return nil, utils.ErrParameterType
	}
	realPath := path.Join(module.BasePath, file)
	data, err := ioutil.ReadFile(realPath)
	if err != nil {
		return nil, err
	}
	code := string(data)

	node, err := ParseString(scope.fset, file, code)
	if err != nil {
		return nil, err
	}

	return scope.Eval(node)
}

var slurpFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	file, ok := args[0].(string)
	if !ok {
		return nil, utils.ErrParameterType
	}
	realPath := path.Join(module.BasePath, file)
	data, err := ioutil.ReadFile(realPath)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}, utils.CheckArity(1))

func docsFn(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		fnames := module.AllFunctionNames()
		res := make([]interface{}, len(fnames))
		for i, f := range fnames {
			res[i] = f
		}

		return res, nil
	}

	if len(args) != 1 {
		return nil, errors.New("Expected 0 or 1 argument")
	}

	expr, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Expected string argument")
	}

	names := module.MatchingFuncNames(expr)
	if len(names) == 0 {
		return "Function not found", nil
	}

	if len(names) == 1 {
		return module.FunctionRepr(names[0]), nil
	}

	res := make([]interface{}, len(names))
	for i, f := range names {
		res[i] = f
	}

	return res, nil
}

func uuidFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 0 {
		return nil, errors.New("uuid function takes no arguments")
	}

	res := uuid.New().String()
	return res, nil
}

var notFn = utils.SimpleFunc(func(b bool) bool {
	return !b
}, utils.CheckArity(1))

func randFn() func(args ...interface{}) (value interface{}, err error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	return func(args ...interface{}) (value interface{}, err error) {
		if len(args) != 0 {
			return nil, errors.New("rand function takes no arguments")
		}

		return r.Float64(), nil
	}
}

var evalFileFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	file, ok := args[0].(string)
	if !ok {
		return nil, utils.ErrParameterType
	}
	realPath := path.Join(module.BasePath, file)
	data, err := ioutil.ReadFile(realPath)
	if err != nil {
		return nil, err
	}
	code := string(data)

	g, err := NewWithName(code, file)
	if err != nil {
		return nil, fmt.Errorf("Error in eval: %v", err)
	}

	return g.Eval(NewEnv())
}, utils.CheckArity(1))

func vecFn(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return []float64{}, nil
	}

	if list, ok := args[0].([]float64); ok {
		return list, nil
	}

	if list, ok := args[0].([]interface{}); ok {
		res := make([]float64, len(list))
		for i, e := range list {
			switch e.(type) {
			case float64:
				res[i] = e.(float64)
			case int64:
				res[i] = float64(e.(int64))
			default:
				return nil, fmt.Errorf("Cannot use %v as float64", e)
			}
		}
		return res, nil
	}

	res := make([]float64, len(args))
	for i, e := range args {
		switch e.(type) {
		case float64:
			res[i] = e.(float64)
		case int64:
			res[i] = float64(e.(int64))
		default:
			return nil, fmt.Errorf("Cannot use %v as float64", e)
		}
	}
	return res, nil
}

var dictKeysFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	dict, ok := args[0].(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("dict-keys expects a dictionary")
	}

	keys := make([]interface{}, len(dict))

	i := 0
	for k := range dict {
		keys[i] = k
		i++
	}

	return keys, nil
}, utils.CheckArity(1))

func vecToListFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, utils.ErrWrongNumberPar
	}

	if list, ok := args[0].([]float64); ok {
		res := make([]interface{}, len(list))
		for i, e := range list {
			res[i] = e
		}
		return res, nil
	}

	return nil, utils.ErrParameterType
}

func listToVecFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, utils.ErrWrongNumberPar
	}

	if list, ok := args[0].([]interface{}); ok {
		res := make([]float64, len(list))
		for i, e := range list {
			switch e.(type) {
			case int64:
				res[i] = float64(e.(int64))
			case int:
				res[i] = float64(e.(int))
			case float64:
				res[i] = e.(float64)
			default:
				return nil, utils.ErrParameterType
			}
		}
		return res, nil
	}

	return nil, utils.ErrParameterType
}

var rangeFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	res := []interface{}{}

	switch start := args[0].(type) {
	case int64:
		step := args[2].(int64)
		end := args[1].(int64)
		if step == 0 {
			return nil, errors.New("Invalid argument")
		} else if step > 0 {
			for i := start; i < end; i += step {
				res = append(res, i)
			}
		} else {
			for i := start; i > end; i += step {
				res = append(res, i)
			}
		}

		return res, nil
	case float64:
		step := args[2].(float64)
		end := args[1].(float64)
		if step == 0 {
			return nil, errors.New("Invalid argument")
		} else if step > 0 {
			for i := start; i < end; i += step {
				res = append(res, i)
			}
		} else {
			for i := start; i > end; i += step {
				res = append(res, i)
			}
		}
		return res, nil
	}

	return nil, utils.ErrParameterType
}, utils.CheckArity(3), utils.ParamsToSameBaseType())

func bindFn(args ...interface{}) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New("bind takes 2 or more arguments")
	}

	fn, ok := args[0].(func(...interface{}) (interface{}, error))
	if !ok {
		return nil, errors.New("Expected function as first argument")
	}

	boundArgs := args[1:]

	return func(args ...interface{}) (interface{}, error) {
		allArgs := []interface{}{}
		allArgs = append(allArgs, boundArgs...)
		allArgs = append(allArgs, args...)
		return fn(allArgs...)
	}, nil
}

var repeatFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n, ok := args[0].(int64)
	if !ok {
		return nil, utils.ErrParameterType
	}

	v := args[1]

	res := make([]interface{}, n)
	for i := range res {
		res[i] = v
	}
	return res, nil
}, utils.CheckArity(2))

var repeatedlyFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n, ok := args[0].(int64)
	if !ok {
		return nil, utils.ErrParameterType
	}

	fn, ok := args[1].(func(args ...interface{}) (interface{}, error))
	if !ok {
		return nil, errors.New("repeatedly takes a function as second parameter")
	}

	res := make([]interface{}, n)
	for i := range res {
		v, err := fn()
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}, utils.CheckArity(2))

var vecRepeatFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n, ok := args[0].(int64)
	if !ok {
		return nil, utils.ErrParameterType
	}

	v, ok := args[1].(float64)
	if !ok {
		vi, ok := args[1].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}

		v = float64(vi)
	}

	res := make([]float64, n)
	for i := range res {
		res[i] = v
	}
	return res, nil
}, utils.CheckArity(2))

var vecRandFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n, ok := args[0].(int64)
	if !ok {
		return nil, utils.ErrParameterType
	}

	rand.Seed(time.Now().UnixNano())

	res := make([]float64, n)
	for i := range res {
		res[i] = rand.Float64()
	}
	return res, nil
}, utils.CheckArity(1))

var listRandFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n, ok := args[0].(int64)
	if !ok {
		return nil, utils.ErrParameterType
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))

	res := make([]interface{}, n)
	for i := range res {
		res[i] = r.Float64()
	}
	return res, nil
}, utils.CheckArity(1))

var vecRangeFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	res := []float64{}

	switch start := args[0].(type) {
	case int64:
		step := args[2].(int64)
		end := args[1].(int64)
		if step == 0 {
			return nil, errors.New("Invalid argument")
		} else if step > 0 {
			for i := start; i < end; i += step {
				res = append(res, float64(i))
			}
		} else {
			for i := start; i > end; i += step {
				res = append(res, float64(i))
			}
		}

		return res, nil
	case float64:
		step := args[2].(float64)
		end := args[1].(float64)
		if step == 0 {
			return nil, errors.New("Invalid argument")
		} else if step > 0 {
			for i := start; i < end; i += step {
				res = append(res, i)
			}
		} else {
			for i := start; i > end; i += step {
				res = append(res, i)
			}
		}
		return res, nil
	}

	return nil, utils.ErrParameterType
}, utils.CheckArity(3), utils.ParamsToSameBaseType())

var isVecFn = utils.SimpleFunc(func(args ...interface{}) interface{} {
	_, ok := args[0].([]float64)
	return ok
}, utils.CheckArity(1))

var isListFn = utils.SimpleFunc(func(args ...interface{}) interface{} {
	_, ok := args[0].([]interface{})
	return ok
}, utils.CheckArity(1))

var isDictFn = utils.SimpleFunc(func(args ...interface{}) interface{} {
	_, ok := args[0].(map[interface{}]interface{})
	return ok
}, utils.CheckArity(1))

var subFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case []interface{}:
		v := arg
		i1, ok := args[1].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}
		i2, ok := args[2].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}

		if i1 < 0 {
			if -int(i1) > len(v) {
				return nil, errors.New("Key not found")
			}

			i1 = int64(len(v)) + i1
		}

		if i2 < 0 {
			if -int(i2) > len(v) {
				return nil, errors.New("Key not found")
			}

			i2 = int64(len(v)) + i2 + 1
		}

		if int(i1) >= len(v) || int(i2) > len(v) || i1 >= i2 {
			return nil, errors.New("Key not found")
		}
		return v[i1:i2], nil
	case []float64:
		v := arg
		i1, ok := args[1].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}
		i2, ok := args[2].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}

		if i1 < 0 {
			if -int(i1) > len(v) {
				return nil, errors.New("Key not found")
			}

			i1 = int64(len(v)) + i1
		}

		if i2 < 0 {
			if -int(i2) > len(v) {
				return nil, errors.New("Key not found")
			}

			i2 = int64(len(v)) + i2 + 1
		}

		if int(i1) >= len(v) || int(i2) > len(v) || i1 >= i2 {
			return nil, errors.New("Key not found")
		}
		return v[i1:i2], nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArity(3))

var containsFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case map[interface{}]interface{}:
		_, ok := arg[args[1]]
		return ok, nil
	case []interface{}:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}

		if int(i) >= len(v) {
			return false, nil
		}
		return true, nil
	case []float64:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}

		if int(i) >= len(v) {
			return false, nil
		}
		return true, nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArity(2))

var jsonFn = utils.ErrFunc(func(arg interface{}) (interface{}, error) {
	var fix func(interface{}) interface{}

	fix = func(arg interface{}) interface{} {
		switch rarg := arg.(type) {
		case map[interface{}]interface{}:
			d := map[string]interface{}{}
			for k, v := range rarg {
				s, ok := k.(string)
				if !ok {
					s = fmt.Sprintf("%v", k)
				}
				d[s] = fix(v)
			}
			return d
		case []interface{}:
			d := make([]interface{}, len(rarg))
			for i, v := range rarg {
				d[i] = fix(v)
			}
			return d
		}
		return arg
	}

	b, err := json.Marshal(fix(arg))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}, utils.CheckArity(1))

var updateFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case map[interface{}]interface{}:
		arg[args[1]] = args[2]
		return args[0], nil
	case []interface{}:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}

		if int(i) >= len(v) {
			return nil, errors.New("Out of range")
		}
		v[i] = args[2]
		return args[0], nil
	case []float64:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, utils.ErrParameterType
		}
		f, ok := args[2].(float64)
		if !ok {
			iv, ok := args[2].(int64)
			if !ok {
				return nil, utils.ErrParameterType
			}
			f = float64(iv)
		}

		if int(i) >= len(v) {
			return args[0], errors.New("Out of range")
		}
		v[i] = f
		return args[0], nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArity(3))

var appendFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case []interface{}:
		v := make([]interface{}, len(arg))
		copy(v, arg)
		for _, n := range args[1:] {
			v = append(v, n)
		}
		return v, nil
	case []float64:
		v := make([]float64, len(arg))
		copy(v, arg)

		for _, n := range args[1:] {
			f, ok := n.(float64)
			if !ok {
				iv, ok := n.(int64)
				if !ok {
					return nil, utils.ErrParameterType
				}
				f = float64(iv)
			}
			v = append(v, f)
		}

		return v, nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArityAtLeast(2))

var concatFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case []interface{}:
		if len(args) == 1 {
			return arg, nil
		}

		v := make([]interface{}, len(arg))
		copy(v, arg)
		for _, n := range args[1:] {
			v2, ok := n.([]interface{})
			if !ok {
				return nil, utils.ErrParameterType
			}
			v = append(v, v2...)
		}
		return v, nil
	case []float64:
		if len(args) == 1 {
			return arg, nil
		}

		v := make([]float64, len(arg))
		copy(v, arg)

		for _, n := range args[1:] {
			v2, ok := n.([]float64)
			if !ok {
				return nil, utils.ErrParameterType
			}
			v = append(v, v2...)
		}

		return v, nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArityAtLeast(1))

var mergeFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	res := map[interface{}]interface{}{}
	for _, arg := range args {
		d, ok := arg.(map[interface{}]interface{})
		if !ok {
			return nil, utils.ErrParameterType
		}
		for k, v := range d {
			res[k] = v
		}
	}
	return res, nil

}, utils.CheckArityAtLeast(2))

var lenFn = utils.ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case map[interface{}]interface{}:
		return int64(len(arg)), nil
	case []interface{}:
		return int64(len(arg)), nil
	case []float64:
		return int64(len(arg)), nil
	case string:
		return int64(len(arg)), nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArity(1))

func eqFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New("== takes two values")
	}
	return args[0] == args[1], nil
}

func neFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New("!= takes two values")
	}
	return args[0] != args[1], nil
}

var plusFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return int64(0), nil
	}

	if utils.IsSlice(args[0]) {
		res := make([]float64, len(args[0].([]float64)))
		for i, v := range args[0].([]float64) {
			res[i] = v
		}

		for _, arg := range args[1:] {
			if len(arg.([]float64)) != len(res) {
				return nil, errors.New("Vectors of different length")
			}
			for i, v := range arg.([]float64) {
				res[i] += v
			}
		}

		return res, nil
	}

	var resi int64
	var resf float64
	var f bool
	for _, arg := range args {
		switch arg := arg.(type) {
		case int64:
			resi += arg
			resf += float64(arg)
		case float64:
			resf += arg
			f = true
		default:
			return nil, fmt.Errorf("cannot sum %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return resi, nil
}, utils.ParamsToSameBaseType(), utils.ParamsSlicify())

var modFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	var resi int64
	switch arg := args[0].(type) {
	case int64:
		resi = arg
	default:
		return nil, fmt.Errorf("cannot $ %#v", arg)
	}

	for _, arg := range args[1:] {
		switch arg := arg.(type) {
		case int64:
			resi %= arg
		default:
			return nil, fmt.Errorf("cannot $ %#v", arg)
		}
	}
	return resi, nil
}, utils.CheckArityAtLeast(2), utils.ParamsToSameBaseType(), utils.ParamsSlicify())

var intFn = utils.SimpleFunc(func(args ...interface{}) (value interface{}) {
	return args[0]
}, utils.CheckArity(1), utils.ParamToInt64(0))

var floatFn = utils.SimpleFunc(func(args ...interface{}) (value interface{}) {
	return args[0]
}, utils.CheckArity(1), utils.ParamToFloat64(0))

var minusFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf(`function "-" takes one or more arguments`)
	}

	if utils.IsSlice(args[0]) {
		res := make([]float64, len(args[0].([]float64)))
		if len(args) == 1 {
			for i, v := range args[0].([]float64) {
				res[i] = -v
			}
		} else {
			for i, v := range args[0].([]float64) {
				res[i] = v
			}
		}

		for _, arg := range args[1:] {
			if len(arg.([]float64)) != len(res) {
				return nil, errors.New("Vectors of different length")
			}
			for i, v := range arg.([]float64) {
				res[i] -= v
			}
		}

		return res, nil
	}

	var resi int64
	var resf float64
	var f bool
	for i, arg := range args {
		switch arg := arg.(type) {
		case int64:
			if i == 0 && len(args) > 1 {
				resi = arg
				resf = float64(arg)
			} else {
				resi -= arg
				resf -= float64(arg)
			}
		case float64:
			if i == 0 && len(args) > 1 {
				resf = arg
			} else {
				resf -= arg
			}
			f = true
		default:
			return nil, fmt.Errorf("cannot subtract %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return resi, nil
}, utils.ParamsToSameBaseType(), utils.ParamsSlicify())

var mulFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return int64(1), nil
	}

	if utils.IsSlice(args[0]) {
		res := make([]float64, len(args[0].([]float64)))
		for i, v := range args[0].([]float64) {
			res[i] = v
		}

		for _, arg := range args[1:] {
			if len(arg.([]float64)) != len(res) {
				return nil, errors.New("Vectors of different length")
			}
			for i, v := range arg.([]float64) {
				res[i] *= v
			}
		}

		return res, nil
	}

	var resi = int64(1)
	var resf = float64(1)
	var f bool
	for _, arg := range args {
		switch arg := arg.(type) {
		case int64:
			resi *= arg
			resf *= float64(arg)
		case float64:
			resf *= arg
			f = true
		default:
			return nil, fmt.Errorf("cannot multiply %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return resi, nil
}, utils.ParamsToSameBaseType(), utils.ParamsSlicify())

var minFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	var resf = math.MaxFloat64
	var f bool
	for _, arg := range args {
		switch arg := arg.(type) {
		case int64:
			resf = math.Min(resf, float64(arg))
		case float64:
			resf = math.Min(resf, arg)
			f = true
		default:
			return nil, fmt.Errorf("cannot min %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return int64(resf), nil
}, utils.CheckArityAtLeast(1))

var maxFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	var resf = -math.MaxFloat64
	var f bool
	for _, arg := range args {
		switch arg := arg.(type) {
		case int64:
			resf = math.Max(resf, float64(arg))
		case float64:
			resf = math.Max(resf, arg)
			f = true
		default:
			return nil, fmt.Errorf("cannot max %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return int64(resf), nil
}, utils.CheckArityAtLeast(1))

var divFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New("function \"/\" takes two or more arguments")
	}

	if utils.IsSlice(args[0]) {
		res := make([]float64, len(args[0].([]float64)))
		for i, v := range args[0].([]float64) {
			res[i] = v
		}

		for _, arg := range args[1:] {
			if len(arg.([]float64)) != len(res) {
				return nil, errors.New("Vectors of different length")
			}
			for i, v := range arg.([]float64) {
				res[i] /= v
			}
		}

		return res, nil
	}

	var resi int64
	var resf float64
	var f bool
	for i, arg := range args {
		switch arg := arg.(type) {
		case int64:
			if i == 0 && len(args) > 1 {
				resi = arg
				resf = float64(arg)
			} else {
				resi /= arg
				resf /= float64(arg)
			}
		case float64:
			if i == 0 && len(args) > 1 {
				resf = float64(arg)
			} else {
				resf /= arg
			}
			f = true
		default:
			return nil, fmt.Errorf("cannot divide with %#v", arg)
		}
	}
	if f {
		return resf, nil
	}
	return resi, nil
}, utils.ParamsToSameBaseType(), utils.ParamsSlicify())

var skipFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n := int(args[0].(int64))

	if !utils.IsSlice(args[1]) {
		return nil, utils.ErrParameterType
	}

	switch args[1].(type) {
	case []interface{}:
		s := args[1].([]interface{})
		if len(s) <= n {
			return []interface{}(nil), nil
		}
		return s[n:len(s)], nil
	case []float64:
		s := args[1].([]float64)
		if len(s) <= n {
			return []float64(nil), nil
		}
		return s[n:len(s)], nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArity(2), utils.ParamToInt64(0))

var reverseFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if !utils.IsSlice(args[0]) {
		return nil, utils.ErrParameterType
	}

	switch arg := args[0].(type) {
	case []interface{}:
		l := len(arg)
		if l == 0 {
			return []interface{}(nil), nil
		}
		res := make([]interface{}, l)
		for i := range arg {
			res[i] = arg[l-i-1]
		}
		return res, nil
	case []float64:
		l := len(arg)
		if l == 0 {
			return []float64(nil), nil
		}
		res := make([]float64, l)
		for i := range arg {
			res[i] = arg[l-i-1]
		}
		return res, nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArity(1))

var takeFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n := int(args[0].(int64))

	if !utils.IsSlice(args[1]) {
		return nil, utils.ErrParameterType
	}

	switch args[1].(type) {
	case []interface{}:
		s := args[1].([]interface{})
		if len(s) < n {
			n = len(s)
		}

		if n <= 0 {
			return []interface{}(nil), nil
		}

		return s[0:n], nil
	case []float64:
		s := args[1].([]float64)
		if len(s) < n {
			n = len(s)
		}

		if n <= 0 {
			return []float64(nil), nil
		}
		return s[0:n], nil
	}
	return nil, utils.ErrParameterType
}, utils.CheckArity(2), utils.ParamToInt64(0))

func andFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) == 0 {
		return true, nil
	}
	for _, arg := range args {
		value, err = scope.Eval(arg)
		if err != nil {
			return nil, err
		}
		if value == false {
			return false, nil
		}
	}
	return value, err
}

func orFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) == 0 {
		return false, nil
	}
	for _, arg := range args {
		value, err = scope.Eval(arg)
		if err != nil {
			return nil, err
		}
		if value != false {
			return value, nil
		}
	}
	return value, err
}

func timeFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, errors.New("time takes 1 argument")
	}

	start := time.Now()
	value, err = scope.Eval(args[0])
	elapsed := time.Since(start)
	scope.Printf("Elapsed %.2f milliseconds\n", float64(elapsed.Nanoseconds())/1e6)
	return value, err
}

func ifFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New(`function "if" takes two or three arguments`)
	}
	value, err = scope.Eval(args[0])
	if err != nil {
		return nil, err
	}
	if value == false {
		if len(args) == 3 {
			return scope.Eval(args[2])
		}
		return false, nil
	}
	return scope.Eval(args[1])
}

func condFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New(`function "cond" takes two or more arguments`)
	}

	i := 0
	for ; i+1 < len(args); i += 2 {
		value, err = scope.Eval(args[i])
		if err != nil {
			return nil, err
		}
		if value != false {
			return scope.Eval(args[i+1])
		}
	}

	if len(args)%2 == 0 {
		return false, nil
	}

	return scope.Eval(args[len(args)-1])
}

func varFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) == 0 || len(args) > 2 {
		return nil, errors.New("var takes one or two arguments")
	}
	symbol, ok := args[0].(*ast.Symbol)
	if !ok {
		return nil, errors.New("var takes a symbol as first argument")
	}
	if len(args) == 1 {
		value = nil
	} else {
		value, err = scope.Eval(args[1])
		if err != nil {
			return nil, err
		}
	}
	return nil, scope.Create(symbol.Name, value)
}

func setFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`function "set" takes two arguments`)
	}
	symbol, ok := args[0].(*ast.Symbol)
	if !ok {
		return nil, errors.New(`function "set" takes a symbol as first argument`)
	}
	value, err = scope.Eval(args[1])
	if err != nil {
		return nil, err
	}
	return nil, scope.Set(symbol.Name, value)
}

func doFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	scope = scope.Branch()
	for _, arg := range args {
		value, err = scope.Eval(arg)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func codeFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, errors.New("code takes one argument")
	}

	return scope.Code(args[0]), nil
}

func funcFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New(`func takes two or more arguments`)
	}
	i := 0
	var name string
	if symbol, ok := args[0].(*ast.Symbol); ok {
		name = symbol.Name
		i++
	}
	list, ok := args[i].(*ast.ListList)
	if !ok {
		return nil, errors.New(`func takes a list of parameters`)
	}
	params := list.Nodes
	for _, param := range params {
		if _, ok := param.(*ast.Symbol); !ok {
			return nil, errors.New("func's list of parameters must be a list of symbols")
		}
	}
	body := args[i+1:]
	if len(body) == 0 {
		return nil, fmt.Errorf("func takes a body sequence")
	}
	fn := func(args ...interface{}) (value interface{}, err error) {
		if len(args) != len(params) {
			nameInfo := "anonymous function"
			if name != "" {
				nameInfo = fmt.Sprintf("function %q", name)
			}
			switch len(params) {
			case 0:
				return nil, fmt.Errorf("%s takes no arguments", nameInfo)
			case 1:
				return nil, fmt.Errorf("%s takes one argument", nameInfo)
			default:
				return nil, fmt.Errorf("%s takes %d arguments", nameInfo, len(params))
			}
		}
		scope = scope.Branch()
		for i, arg := range args {
			err := scope.Create(params[i].(*ast.Symbol).Name, arg)
			if err != nil {
				panic("must not happen: " + err.Error())
			}
		}
		for _, node := range body {
			value, err = scope.Eval(node)
			if err != nil {
				return nil, err
			}
		}
		return value, nil
	}
	if name != "" {
		if err = scope.Create(name, fn); err != nil {
			return nil, err
		}
	}
	return fn, nil
}

func macroFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, errors.New(`# takes one argument`)
	}

	node := args[0]

	fn := func(args ...interface{}) (value interface{}, err error) {
		for i, arg := range args {
			scope.SetOrCreate(fmt.Sprintf("%%%v", i+1), arg)
		}

		return scope.Eval(node)
	}

	return fn, nil
}

var mapFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	lists := [][]interface{}{}
	for _, arg := range args[1:] {
		list, ok := utils.ToList(arg)
		if !ok {
			return nil, utils.ErrParameterType
		}
		lists = append(lists, list)
	}

	l := len(lists[0])

	res := []interface{}{}
	for i := 0; i < l; i++ {
		fnArgs := make([]interface{}, len(lists))
		for j, list := range lists {
			if len(list) != l {
				return nil, errors.New("Lists must be of same length")
			}
			fnArgs[j] = list[i]
		}
		r, err := utils.Call(fn, fnArgs...)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}, utils.CheckArityAtLeast(2))

var mapIndexedFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	lists := [][]interface{}{}
	for _, arg := range args[1:] {
		list, ok := utils.ToList(arg)
		if !ok {
			return nil, utils.ErrParameterType
		}
		lists = append(lists, list)
	}

	l := len(lists[0])

	res := []interface{}{}
	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		for i := 0; i < l; i++ {
			fnArgs := make([]interface{}, len(lists)+1)
			fnArgs[0] = int64(i)
			for j, list := range lists {
				if len(list) != l {
					return nil, errors.New("Lists must be of same length")
				}
				fnArgs[j+1] = list[i]
			}
			r, err := fn(fnArgs...)
			if err != nil {
				return nil, err
			}
			res = append(res, r)
		}

		return res, nil
	}

	return nil, utils.ErrParameterType
}, utils.CheckArityAtLeast(2))

var sortAscFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]
	list, ok := args[1].([]interface{})
	if !ok {
		return nil, utils.ErrParameterType
	}

	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		res := make([]interface{}, len(list))
		copy(res, list)
		sort.Slice(res, func(i, j int) bool {
			v, _ := fn(res[i], res[j])
			return v.(bool)
		})

		return res, nil
	}

	return nil, utils.ErrParameterType
}, utils.CheckArity(2))

var sortIndexFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]
	list, ok := args[1].([]interface{})
	if !ok {
		return nil, utils.ErrParameterType
	}

	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		return SortIndex(list, func(v1, v2 interface{}) bool {
			v, _ := fn(v1, v2)
			return v.(bool)
		}), nil
	}

	return nil, utils.ErrParameterType
}, utils.CheckArity(2))

var sortDescFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]
	list, ok := args[1].([]interface{})
	if !ok {
		return nil, utils.ErrParameterType
	}

	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		res := make([]interface{}, len(list))
		copy(res, list)
		sort.Slice(res, func(i, j int) bool {
			v, _ := fn(res[j], res[i])
			return v.(bool)
		})

		return res, nil
	}

	return nil, utils.ErrParameterType
}, utils.CheckArity(2))

func reduceFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 2 && len(args) != 3 {
		return nil, errors.New(`reduce takes 2 or three arguments arguments`)
	}

	fn, err := scope.Eval(args[0])
	if err != nil {
		return nil, scope.errorAt(args[0], err)
	}
	listRaw, err := scope.Eval(args[1])
	if err != nil {
		return nil, scope.errorAt(args[1], err)
	}

	list, ok := listRaw.([]interface{})
	if !ok {
		return nil, utils.ErrParameterType
	}

	var init interface{}

	if len(args) == 3 {
		init, err = scope.Eval(args[2])
		if err != nil {
			return nil, scope.errorAt(args[2], err)
		}
	}

	r := init
	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		for i, v := range list {
			if i == 0 && r == nil {
				r = v
			} else {
				r, err = fn(r, v)
				if err != nil {
					return nil, err
				}
			}
		}

		return r, nil
	}

	return nil, utils.ErrParameterType
}

func filterFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`filter takes two arguments`)
	}

	fn := args[0]

	switch list := args[1].(type) {
	case []interface{}:
		res := []interface{}{}
		for _, v := range list {
			r, err := utils.Call(fn, v)
			if err != nil {
				return nil, err
			}
			b, ok := r.(bool)
			if !ok {
				return nil, errors.New("callback must return bool")
			}

			if b {
				res = append(res, v)
			}
		}
		return res, nil
	case []float64:
		res := []float64{}
		for _, v := range list {
			r, err := utils.Call(fn, v)
			if err != nil {
				return nil, err
			}
			b, ok := r.(bool)
			if !ok {
				return nil, errors.New("callback must return bool")
			}

			if b {
				res = append(res, v)
			}
		}

		return res, nil
	}

	return nil, utils.ErrParameterType
}

func countIfFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`count-if takes two arguments`)
	}

	fn := args[0]

	switch list := args[1].(type) {
	case []interface{}:
		res := 0
		for _, v := range list {
			r, err := utils.Call(fn, v)
			if err != nil {
				return nil, err
			}
			b, ok := r.(bool)
			if !ok {
				return nil, errors.New("callback must return bool")
			}

			if b {
				res++
			}
		}

		return int64(res), nil
	case []float64:
		res := 0
		for _, v := range list {
			r, err := utils.Call(fn, v)
			if err != nil {
				return nil, err
			}
			b, ok := r.(bool)
			if !ok {
				return nil, errors.New("callback must return bool")
			}

			if b {
				res++
			}
		}

		return int64(res), nil
	}

	return nil, utils.ErrParameterType
}

var applyFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	fnArgs, ok := args[1].([]interface{})
	if !ok {
		return nil, utils.ErrParameterType
	}

	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		return fn(fnArgs...)
	}

	return nil, utils.ErrParameterType
}, utils.CheckArity(2))

var flattenFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	res := []interface{}{}
	var sub func(args []interface{})
	sub = func(args []interface{}) {
		for _, arg := range args {
			switch arg.(type) {
			case []interface{}:
				sub(arg.([]interface{}))
			default:
				res = append(res, arg)
			}
		}
	}
	sub(args)

	return res, nil
}, utils.CheckArityAtLeast(1))

var vecApplyFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	list, ok := args[1].([]float64)
	if !ok {
		return nil, utils.ErrParameterType
	}

	fnArgs := make([]interface{}, len(list))
	for i, v := range list {
		fnArgs[i] = v
	}

	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		return fn(fnArgs...)
	}

	return nil, utils.ErrParameterType
}, utils.CheckArity(2))

var vecMapFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	lists := [][]float64{}
	for _, arg := range args[1:] {
		list, ok := arg.([]float64)
		if !ok {
			return nil, utils.ErrParameterType
		}
		lists = append(lists, list)
	}

	l := len(lists[0])

	res := []float64{}
	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		for i := 0; i < l; i++ {
			fnArgs := make([]interface{}, len(lists))
			for j, list := range lists {
				if len(list) != l {
					return nil, errors.New("Lists must be of same length")
				}
				fnArgs[j] = list[i]
			}
			r, err := fn(fnArgs...)
			if err != nil {
				return nil, err
			}
			v, ok := r.(float64)
			if !ok {
				return nil, errors.New("Expected function to return float64")
			}
			res = append(res, v)
		}

		return res, nil
	}

	return nil, utils.ErrParameterType
}, utils.CheckArityAtLeast(2))

var vecMapIndexedFn = utils.ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	lists := [][]float64{}
	for _, arg := range args[1:] {
		list, ok := arg.([]float64)
		if !ok {
			return nil, utils.ErrParameterType
		}
		lists = append(lists, list)
	}

	l := len(lists[0])

	res := []float64{}
	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		for i := 0; i < l; i++ {
			fnArgs := make([]interface{}, len(lists)+1)
			fnArgs[0] = int64(i)
			for j, list := range lists {
				if len(list) != l {
					return nil, errors.New("Lists must be of same length")
				}
				fnArgs[j+1] = list[i]
			}
			r, err := fn(fnArgs...)
			if err != nil {
				return nil, err
			}
			v, ok := r.(float64)
			if !ok {
				return nil, errors.New("Expected function to return float64")
			}
			res = append(res, v)
		}

		return res, nil
	}

	return nil, utils.ErrParameterType
}, utils.CheckArityAtLeast(2))

func forFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 4 {
		return nil, errors.New(`for takes four or more arguments`)
	}
	init, test, step, code := args[0], args[1], args[2], args[3:]
	scope = scope.Branch()
	_, err = scope.Eval(init)
	if err != nil {
		return nil, err
	}
	for {
		more, err := scope.Eval(test)
		if err != nil {
			return nil, err
		}
		if more == false {
			return value, nil
		}

		for _, c := range code {
			value, err = scope.Eval(c)
			if err != nil {
				return nil, err
			}
		}

		_, err = scope.Eval(step)
		if err != nil {
			return nil, err
		}
	}
	panic("unreachable")
}

var lessThanEqualFn = utils.SimpleFunc(func(v ...interface{}) bool {
	switch v[0].(type) {
	case float64:
		return v[0].(float64) <= v[1].(float64)
	case string:
		return v[0].(string) <= v[1].(string)
	default:
		return v[0].(int64) <= v[1].(int64)
	}
}, utils.CheckArity(2), utils.ParamsToSameBaseType())

var greaterThanEqualFn = utils.SimpleFunc(func(v ...interface{}) bool {
	switch v[0].(type) {
	case float64:
		return v[0].(float64) >= v[1].(float64)
	case string:
		return v[0].(string) >= v[1].(string)
	default:
		return v[0].(int64) >= v[1].(int64)
	}
}, utils.CheckArity(2), utils.ParamsToSameBaseType())

var lessThanFn = utils.SimpleFunc(func(v ...interface{}) bool {
	switch v[0].(type) {
	case float64:
		return v[0].(float64) < v[1].(float64)
	case string:
		return v[0].(string) < v[1].(string)
	default:
		return v[0].(int64) < v[1].(int64)
	}
}, utils.CheckArity(2), utils.ParamsToSameBaseType())

var greaterThanFn = utils.SimpleFunc(func(v ...interface{}) bool {
	switch v[0].(type) {
	case float64:
		return v[0].(float64) > v[1].(float64)
	case string:
		return v[0].(string) > v[1].(string)
	default:
		return v[0].(int64) > v[1].(int64)
	}
}, utils.CheckArity(2), utils.ParamsToSameBaseType())
