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
	"github.com/google/uuid"
)

func init() {
	RegisterModules(GlobalsModule)
}

var GlobalsModule = &Module{
	Name: "globals",
	Funcs: []*Func{
		&Func{Name: "eval", F: evalFn},
		&Func{Name: "eval-file", F: evalFileFn},
		&Func{Name: "slurp", F: slurpFn},
		&Func{Name: "true", F: true},
		&Func{Name: "false", F: false},
		&Func{Name: "nil", F: nil},
		&Func{Name: "nan", F: math.NaN()},
		&Func{Name: "error", F: errorFn},
		&Func{Name: "==", F: eqFn},
		&Func{Name: "<", F: lessThanFn},
		&Func{Name: ">", F: greaterThanFn},
		&Func{Name: "<=", F: lessThanEqualFn},
		&Func{Name: ">=", F: greaterThanEqualFn},
		&Func{Name: "!=", F: neFn},
		&Func{Name: "+", F: plusFn},
		&Func{Name: "-", F: minusFn},
		&Func{Name: "*", F: mulFn},
		&Func{Name: "/", F: divFn},
		&Func{Name: "%", F: modFn},
		&Func{Name: "int", F: intFn},
		&Func{Name: "float", F: floatFn},
		&Func{Name: "min", F: minFn},
		&Func{Name: "max", F: maxFn},
		&Func{Name: "or", F: orFn},
		&Func{Name: "and", F: andFn},
		&Func{Name: "if", F: ifFn},
		&Func{Name: "cond", F: condFn},
		&Func{Name: "var", F: varFn},
		&Func{Name: "set", F: setFn},
		&Func{Name: "do", F: doFn},
		&Func{Name: "code", F: codeFn},
		&Func{Name: "func", F: funcFn},
		&Func{Name: "for", F: forFn},
		&Func{Name: "vec", F: vecFn},
		&Func{Name: "vec2list", F: vecToListFn},
		&Func{Name: "list2vec", F: listToVecFn},
		&Func{Name: "list", F: listFn},
		&Func{Name: "vec?", F: isVecFn},
		&Func{Name: "list?", F: isListFn},
		&Func{Name: "dict", F: dictFn},
		&Func{Name: "dict?", F: isDictFn},
		&Func{Name: "dict-keys", F: dictKeysFn},
		&Func{Name: "get", F: getFn},
		&Func{Name: "sub", F: subFn},
		&Func{Name: "contains?", F: containsFn},
		&Func{Name: "update!", F: updateFn},
		&Func{Name: "len", F: lenFn},
		&Func{Name: "append", F: appendFn},
		&Func{Name: "concat", F: concatFn},
		&Func{Name: "merge", F: mergeFn},
		&Func{Name: "range", F: rangeFn},
		&Func{Name: "vec-range", F: vecRangeFn},
		&Func{Name: "repeat", F: repeatFn},
		&Func{Name: "reverse", F: reverseFn},
		&Func{Name: "vec-repeat", F: vecRepeatFn},
		&Func{Name: "map", F: mapFn},
		&Func{Name: "map-indexed", F: mapIndexedFn},
		&Func{Name: "vec-map", F: vecMapFn},
		&Func{Name: "vec-map-indexed", F: vecMapIndexedFn},
		&Func{Name: "apply", F: applyFn},
		&Func{Name: "vec-apply", F: vecApplyFn},
		&Func{Name: "vec-rand", F: vecRandFn},
		&Func{Name: "reduce", F: reduceFn},
		&Func{Name: "filter", F: filterFn},
		&Func{Name: "flatten", F: flattenFn},
		&Func{Name: "skip", F: skipFn},
		&Func{Name: "take", F: takeFn},
		&Func{Name: "sort-asc", F: sortAscFn},
		&Func{Name: "sort-desc", F: sortDescFn},
		&Func{Name: "sortindex", F: sortIndexFn},
		&Func{Name: "bind", F: bindFn},
		&Func{Name: "json", F: jsonFn},
		&Func{Name: "uuid", F: uuidFn},
	},
	LispFuncs: []*LispFunc{
		&LispFunc{Name: "identity", F: "(func (x) x)"},
		&LispFunc{Name: "empty?", F: "(func (x) (== (len x) 0))"},
		&LispFunc{Name: "first", F: "(func (s) (get s 0))"},
		&LispFunc{Name: "rest", F: "(func (s) (skip 1 s))"},
		&LispFunc{Name: "last", F: "(func (s) (if (empty? s) nil (get s (- (len s) 1))))"},
		&LispFunc{Name: "inc", F: "(func (s) (+ s 1))"},
		&LispFunc{Name: "dec", F: "(func (s) (- s 1))"},
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

var evalFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	code, ok := args[0].(string)
	if !ok {
		return nil, errParameterType
	}

	g, err := New(code)
	if err != nil {
		return nil, fmt.Errorf("Error in eval: %v", err)
	}

	return g.Eval(NewEnv())
}, CheckArity(1))

var slurpFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	file, ok := args[0].(string)
	if !ok {
		return nil, errParameterType
	}
	realPath := path.Join(BasePath, file)
	data, err := ioutil.ReadFile(realPath)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}, CheckArity(1))

func uuidFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 0 {
		return nil, errors.New("uuid function takes no arguments")
	}

	res := uuid.New().String()
	return res, nil
}

var evalFileFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	file, ok := args[0].(string)
	if !ok {
		return nil, errParameterType
	}
	realPath := path.Join(BasePath, file)
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
}, CheckArity(1))

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

var dictFn = SimpleFunc(func(args ...interface{}) interface{} {
	if len(args) == 0 {
		return map[interface{}]interface{}{}
	}

	res := make(map[interface{}]interface{})
	for i := 0; i+1 < len(args); i += 2 {
		res[args[i]] = args[i+1]
	}

	return res
}, CheckArityEven())

var dictKeysFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
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
}, CheckArity(1))

func listFn(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return []interface{}{}, nil
	}

	res := make([]interface{}, len(args))
	for i, arg := range args {
		res[i] = arg
	}

	return res, nil
}

func vecToListFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, errWrongNumberPar
	}

	if list, ok := args[0].([]float64); ok {
		res := make([]interface{}, len(list))
		for i, e := range list {
			res[i] = e
		}
		return res, nil
	}

	return nil, errParameterType
}

func listToVecFn(args ...interface{}) (value interface{}, err error) {
	if len(args) != 1 {
		return nil, errWrongNumberPar
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
				return nil, errParameterType
			}
		}
		return res, nil
	}

	return nil, errParameterType
}

var rangeFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
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

	return nil, errParameterType
}, CheckArity(3), ParamsToSameBaseType())

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

var repeatFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n, ok := args[0].(int64)
	if !ok {
		return nil, errParameterType
	}

	v := args[1]

	res := make([]interface{}, n)
	for i := range res {
		res[i] = v
	}
	return res, nil
}, CheckArity(2))

var vecRepeatFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n, ok := args[0].(int64)
	if !ok {
		return nil, errParameterType
	}

	v, ok := args[1].(float64)
	if !ok {
		vi, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}

		v = float64(vi)
	}

	res := make([]float64, n)
	for i := range res {
		res[i] = v
	}
	return res, nil
}, CheckArity(2))

var vecRandFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n, ok := args[0].(int64)
	if !ok {
		return nil, errParameterType
	}

	rand.Seed(time.Now().UnixNano())

	res := make([]float64, n)
	for i := range res {
		res[i] = rand.Float64()
	}
	return res, nil
}, CheckArity(1))

var vecRangeFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
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

	return nil, errParameterType
}, CheckArity(3), ParamsToSameBaseType())

var isVecFn = SimpleFunc(func(args ...interface{}) interface{} {
	_, ok := args[0].([]float64)
	return ok
}, CheckArity(1))

var isListFn = SimpleFunc(func(args ...interface{}) interface{} {
	_, ok := args[0].([]interface{})
	return ok
}, CheckArity(1))

var isDictFn = SimpleFunc(func(args ...interface{}) interface{} {
	_, ok := args[0].(map[interface{}]interface{})
	return ok
}, CheckArity(1))

var getFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case map[interface{}]interface{}:
		v, ok := arg[args[1]]
		if !ok {
			return nil, errors.New("Key not found")
		}
		return v, nil
	case []interface{}:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}

		if i < 0 {
			if -int(i) > len(v) {
				return nil, errors.New("Key not found")
			}

			i = int64(len(v)) + i
		}

		if int(i) >= len(v) {
			return nil, errors.New("Key not found")
		}
		return v[i], nil
	case []float64:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}

		if i < 0 {
			if -int(i) > len(v) {
				return nil, errors.New("Key not found")
			}

			i = int64(len(v)) + i
		}

		if int(i) >= len(v) {
			return nil, errors.New("Key not found")
		}
		return v[i], nil
	}
	return nil, errParameterType
}, CheckArity(2))

var subFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case []interface{}:
		v := arg
		i1, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}
		i2, ok := args[2].(int64)
		if !ok {
			return nil, errParameterType
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
			return nil, errParameterType
		}
		i2, ok := args[2].(int64)
		if !ok {
			return nil, errParameterType
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
	return nil, errParameterType
}, CheckArity(3))

var containsFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case map[interface{}]interface{}:
		_, ok := arg[args[1]]
		return ok, nil
	case []interface{}:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}

		if int(i) >= len(v) {
			return false, nil
		}
		return true, nil
	case []float64:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}

		if int(i) >= len(v) {
			return false, nil
		}
		return true, nil
	}
	return nil, errParameterType
}, CheckArity(2))

var jsonFn = ErrFunc(func(arg interface{}) (interface{}, error) {
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
		}
		return arg
	}

	b, err := json.Marshal(fix(arg))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}, CheckArity(1))

var updateFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case map[interface{}]interface{}:
		arg[args[1]] = args[2]
		return args[0], nil
	case []interface{}:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
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
			return nil, errParameterType
		}
		f, ok := args[2].(float64)
		if !ok {
			iv, ok := args[2].(int64)
			if !ok {
				return nil, errParameterType
			}
			f = float64(iv)
		}

		if int(i) >= len(v) {
			return args[0], errors.New("Out of range")
		}
		v[i] = f
		return args[0], nil
	}
	return nil, errParameterType
}, CheckArity(3))

var appendFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
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
					return nil, errParameterType
				}
				f = float64(iv)
			}
			v = append(v, f)
		}

		return v, nil
	}
	return nil, errParameterType
}, CheckArityAtLeast(2))

var concatFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case []interface{}:
		v := make([]interface{}, len(arg))
		copy(v, arg)
		for _, n := range args[1:] {
			v2, ok := n.([]interface{})
			if !ok {
				return nil, errParameterType
			}
			v = append(v, v2...)
		}
		return v, nil
	case []float64:
		v := make([]float64, len(arg))
		copy(v, arg)

		for _, n := range args[1:] {
			v2, ok := n.([]float64)
			if !ok {
				return nil, errParameterType
			}
			v = append(v, v2...)
		}

		return v, nil
	}
	return nil, errParameterType
}, CheckArityAtLeast(2))

var mergeFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	res := map[interface{}]interface{}{}
	for _, arg := range args {
		d, ok := arg.(map[interface{}]interface{})
		if !ok {
			return nil, errParameterType
		}
		for k, v := range d {
			res[k] = v
		}
	}
	return res, nil

}, CheckArityAtLeast(2))

var lenFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case map[interface{}]interface{}:
		return int64(len(arg)), nil
	case []interface{}:
		return int64(len(arg)), nil
	case []float64:
		return int64(len(arg)), nil
	}
	return nil, errParameterType
}, CheckArity(1))

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

var plusFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return int64(0), nil
	}

	if IsSlice(args[0]) {
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
}, ParamsToSameBaseType(), ParamsSlicify())

var modFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
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
}, CheckArityAtLeast(2), ParamsToSameBaseType(), ParamsSlicify())

var intFn = SimpleFunc(func(args ...interface{}) (value interface{}) {
	return args[0]
}, CheckArity(1), ParamToInt64(0))

var floatFn = SimpleFunc(func(args ...interface{}) (value interface{}) {
	return args[0]
}, CheckArity(1), ParamToFloat64(0))

var minusFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf(`function "-" takes one or more arguments`)
	}

	if IsSlice(args[0]) {
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
}, ParamsToSameBaseType(), ParamsSlicify())

var mulFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return int64(1), nil
	}

	if IsSlice(args[0]) {
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
}, ParamsToSameBaseType(), ParamsSlicify())

var minFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
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
}, CheckArityAtLeast(1))

var maxFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
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
}, CheckArityAtLeast(1))

var divFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New("function \"/\" takes two or more arguments")
	}

	if IsSlice(args[0]) {
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
}, ParamsToSameBaseType(), ParamsSlicify())

var skipFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n := int(args[0].(int64))

	if !IsSlice(args[1]) {
		return nil, errParameterType
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
	return nil, errParameterType
}, CheckArity(2), ParamToInt64(0))

var reverseFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	if !IsSlice(args[0]) {
		return nil, errParameterType
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
	return nil, errParameterType
}, CheckArity(1))

var takeFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	n := int(args[0].(int64))

	if !IsSlice(args[1]) {
		return nil, errParameterType
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
	return nil, errParameterType
}, CheckArity(2), ParamToInt64(0))

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
	list, ok := args[i].(*ast.List)
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

var mapFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	lists := [][]interface{}{}
	for _, arg := range args[1:] {
		list, ok := ToList(arg)
		if !ok {
			return nil, errParameterType
		}
		lists = append(lists, list)
	}

	l := len(lists[0])

	res := []interface{}{}
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
			res = append(res, r)
		}

		return res, nil
	}

	return nil, errParameterType
}, CheckArityAtLeast(2))

var mapIndexedFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	lists := [][]interface{}{}
	for _, arg := range args[1:] {
		list, ok := ToList(arg)
		if !ok {
			return nil, errParameterType
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

	return nil, errParameterType
}, CheckArityAtLeast(2))

var sortAscFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]
	list, ok := args[1].([]interface{})
	if !ok {
		return nil, errParameterType
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

	return nil, errParameterType
}, CheckArity(2))

var sortIndexFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]
	list, ok := args[1].([]interface{})
	if !ok {
		return nil, errParameterType
	}

	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		return SortIndex(list, func(v1, v2 interface{}) bool {
			v, _ := fn(v1, v2)
			return v.(bool)
		}), nil
	}

	return nil, errParameterType
}, CheckArity(2))

var sortDescFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]
	list, ok := args[1].([]interface{})
	if !ok {
		return nil, errParameterType
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

	return nil, errParameterType
}, CheckArity(2))

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
		return nil, errParameterType
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

	return nil, errParameterType
}

func filterFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`filter takes two arguments`)
	}

	fn, err := scope.Eval(args[0])
	if err != nil {
		return nil, scope.errorAt(args[0], err)
	}
	listRaw, err := scope.Eval(args[1])
	if err != nil {
		return nil, scope.errorAt(args[1], err)
	}

	switch listRaw.(type) {
	case []interface{}:
		list := listRaw.([]interface{})

		res := []interface{}{}
		if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
			for _, v := range list {
				r, err := fn(v)
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
	case []float64:
		list := listRaw.([]float64)

		res := []float64{}
		if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
			for _, v := range list {
				r, err := fn(v)
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
	}

	return nil, errParameterType
}

var applyFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	fnArgs, ok := args[1].([]interface{})
	if !ok {
		return nil, errParameterType
	}

	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		return fn(fnArgs...)
	}

	return nil, errParameterType
}, CheckArity(2))

var flattenFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
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
}, CheckArityAtLeast(1))

var vecApplyFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	list, ok := args[1].([]float64)
	if !ok {
		return nil, errParameterType
	}

	fnArgs := make([]interface{}, len(list))
	for i, v := range list {
		fnArgs[i] = v
	}

	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		return fn(fnArgs...)
	}

	return nil, errParameterType
}, CheckArity(2))

var vecMapFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	lists := [][]float64{}
	for _, arg := range args[1:] {
		list, ok := arg.([]float64)
		if !ok {
			return nil, errParameterType
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

	return nil, errParameterType
}, CheckArityAtLeast(2))

var vecMapIndexedFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	fn := args[0]

	lists := [][]float64{}
	for _, arg := range args[1:] {
		list, ok := arg.([]float64)
		if !ok {
			return nil, errParameterType
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

	return nil, errParameterType
}, CheckArityAtLeast(2))

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

var lessThanEqualFn = SimpleFunc(func(v ...interface{}) bool {
	switch v[0].(type) {
	case float64:
		return v[0].(float64) <= v[1].(float64)
	case string:
		return v[0].(string) <= v[1].(string)
	default:
		return v[0].(int64) <= v[1].(int64)
	}
}, CheckArity(2), ParamsToSameBaseType())

var greaterThanEqualFn = SimpleFunc(func(v ...interface{}) bool {
	switch v[0].(type) {
	case float64:
		return v[0].(float64) >= v[1].(float64)
	case string:
		return v[0].(string) >= v[1].(string)
	default:
		return v[0].(int64) >= v[1].(int64)
	}
}, CheckArity(2), ParamsToSameBaseType())

var lessThanFn = SimpleFunc(func(v ...interface{}) bool {
	switch v[0].(type) {
	case float64:
		return v[0].(float64) < v[1].(float64)
	case string:
		return v[0].(string) < v[1].(string)
	default:
		return v[0].(int64) < v[1].(int64)
	}
}, CheckArity(2), ParamsToSameBaseType())

var greaterThanFn = SimpleFunc(func(v ...interface{}) bool {
	switch v[0].(type) {
	case float64:
		return v[0].(float64) > v[1].(float64)
	case string:
		return v[0].(string) > v[1].(string)
	default:
		return v[0].(int64) > v[1].(int64)
	}
}, CheckArity(2), ParamsToSameBaseType())
