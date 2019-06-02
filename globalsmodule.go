package gel

import (
	"errors"
	"fmt"

	"github.com/Stromberg/gel/ast"
)

var GlobalsModule = &Module{
	Name: "globals",
	Funcs: []*Func{
		&Func{Name: "true", F: true},
		&Func{Name: "false", F: false},
		&Func{Name: "nil", F: nil},
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
		&Func{Name: "or", F: orFn},
		&Func{Name: "and", F: andFn},
		&Func{Name: "if", F: ifFn},
		&Func{Name: "cond", F: condFn},
		&Func{Name: "var", F: varFn},
		&Func{Name: "set", F: setFn},
		&Func{Name: "do", F: doFn},
		&Func{Name: "func", F: funcFn},
		&Func{Name: "for", F: forFn},
		&Func{Name: "vec", F: vecFn},
		&Func{Name: "list", F: listFn},
		&Func{Name: "vec?", F: isVecFn},
		&Func{Name: "list?", F: isListFn},
		&Func{Name: "dict", F: dictFn},
		&Func{Name: "dict?", F: isDictFn},
		&Func{Name: "dict-keys", F: dictKeysFn},
		&Func{Name: "get", F: getFn},
		&Func{Name: "contains?", F: containsFn},
		&Func{Name: "update", F: updateFn},
		&Func{Name: "len", F: lenFn},
		&Func{Name: "append", F: appendFn},
		&Func{Name: "range", F: rangeFn},
		&Func{Name: "vec-range", F: vecRangeFn},
		&Func{Name: "repeat", F: repeatFn},
		&Func{Name: "vec-repeat", F: vecRepeatFn},
		&Func{Name: "map", F: mapFn},
		&Func{Name: "vec-map", F: vecMapFn},
		&Func{Name: "apply", F: applyFn},
		&Func{Name: "vec-apply", F: vecApplyFn},
		&Func{Name: "reduce", F: reduceFn},
	},
	LispFuncs: []*LispFunc{
		&LispFunc{Name: "identity", F: "(func (x) x)"},
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

	if list, ok := args[0].([]interface{}); ok {
		return list, nil
	}

	if list, ok := args[0].([]float64); ok {
		res := make([]interface{}, len(list))
		for i, e := range list {
			res[i] = e
		}
		return res, nil
	}

	return args, nil
}

var rangeFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	res := []interface{}{}

	switch args[0].(type) {
	case int64:
		for i := args[0].(int64); i < args[1].(int64); i += args[2].(int64) {
			res = append(res, i)
		}
		return res, nil
	case float64:
		for i := args[0].(float64); i < args[1].(float64); i += args[2].(float64) {
			res = append(res, i)
		}
		return res, nil
	}

	return nil, errParameterType
}, CheckArity(3), ParamsToSameBaseType())

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

var vecRangeFn = ErrFunc(func(args ...interface{}) (value interface{}, err error) {
	res := []float64{}

	switch args[0].(type) {
	case int64:
		for i := args[0].(int64); i < args[1].(int64); i += args[2].(int64) {
			res = append(res, float64(i))
		}
		return res, nil
	case float64:
		for i := args[0].(float64); i < args[1].(float64); i += args[2].(float64) {
			res = append(res, i)
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
	switch args[0].(type) {
	case map[interface{}]interface{}:
		v, ok := args[0].(map[interface{}]interface{})[args[1]]
		if !ok {
			return nil, errors.New("Key not found")
		}
		return v, nil
	case []interface{}:
		v := args[0].([]interface{})
		i, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}

		if int(i) >= len(v) {
			return nil, errors.New("Key not found")
		}
		return v[i], nil
	case []float64:
		v := args[0].([]float64)
		i, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}

		if int(i) >= len(v) {
			return nil, errors.New("Key not found")
		}
		return v[i], nil
	}
	return nil, errParameterType
}, CheckArity(2))

var containsFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch args[0].(type) {
	case map[interface{}]interface{}:
		_, ok := args[0].(map[interface{}]interface{})[args[1]]
		return ok, nil
	case []interface{}:
		v := args[0].([]interface{})
		i, ok := args[1].(int64)
		if !ok {
			return nil, errParameterType
		}

		if int(i) >= len(v) {
			return false, nil
		}
		return true, nil
	case []float64:
		v := args[0].([]float64)
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

var updateFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch args[0].(type) {
	case map[interface{}]interface{}:
		args[0].(map[interface{}]interface{})[args[1]] = args[2]
		return args[0], nil
	case []interface{}:
		v := args[0].([]interface{})
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
		v := args[0].([]float64)
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
	switch args[0].(type) {
	case []interface{}:
		v := args[0].([]interface{})
		for _, n := range args[1:] {
			v = append(v, n)
		}
		return v, nil
	case []float64:
		v := args[0].([]float64)

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

var lenFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch args[0].(type) {
	case map[interface{}]interface{}:
		return int64(len(args[0].(map[interface{}]interface{}))), nil
	case []interface{}:
		return int64(len(args[0].([]interface{}))), nil
	case []float64:
		return int64(len(args[0].([]float64))), nil
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

func funcFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New(`func takes three or more arguments`)
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

func mapFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`map takes two arguments`)
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

	res := []interface{}{}
	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		for _, v := range list {
			r, err := fn(v)
			if err != nil {
				return nil, err
			}
			res = append(res, r)
		}

		return res, nil
	}

	return nil, errParameterType
}

func reduceFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 3 {
		return nil, errors.New(`reduce takes three arguments`)
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

	init, err := scope.Eval(args[2])
	if err != nil {
		return nil, scope.errorAt(args[2], err)
	}

	r := init
	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		for _, v := range list {
			r, err = fn(r, v)
			if err != nil {
				return nil, err
			}
		}

		return r, nil
	}

	return nil, errParameterType
}

func applyFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New(`apply takes two or more arguments`)
	}

	fn, err := scope.Eval(args[0])
	if err != nil {
		return nil, scope.errorAt(args[0], err)
	}

	lists := [][]interface{}{}
	for _, arg := range args[1:] {
		listRaw, err := scope.Eval(arg)
		if err != nil {
			return nil, scope.errorAt(arg, err)
		}
		list, ok := listRaw.([]interface{})
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
}

func vecApplyFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 2 {
		return nil, errors.New(`vec-apply takes two or more arguments`)
	}

	fn, err := scope.Eval(args[0])
	if err != nil {
		return nil, scope.errorAt(args[0], err)
	}

	lists := [][]float64{}
	for _, arg := range args[1:] {
		listRaw, err := scope.Eval(arg)
		if err != nil {
			return nil, scope.errorAt(arg, err)
		}
		list, ok := listRaw.([]float64)
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
			f, ok := r.(float64)
			if !ok {
				vi, ok := r.(int64)
				if !ok {
					return nil, errParameterType
				}
				f = float64(vi)
			}
			res = append(res, f)
		}

		return res, nil
	}

	return nil, errParameterType
}

func vecMapFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) != 2 {
		return nil, errors.New(`vec-map takes two arguments`)
	}

	fn, err := scope.Eval(args[0])
	if err != nil {
		return nil, scope.errorAt(args[0], err)
	}
	listRaw, err := scope.Eval(args[1])
	if err != nil {
		return nil, scope.errorAt(args[1], err)
	}

	list, ok := listRaw.([]float64)
	if !ok {
		return nil, errParameterType
	}

	res := []float64{}
	if fn, ok := fn.(func(...interface{}) (interface{}, error)); ok {
		for _, v := range list {
			r, err := fn(v)
			if err != nil {
				return nil, err
			}
			f, ok := r.(float64)
			if !ok {
				vi, ok := r.(int64)
				if !ok {
					return nil, errParameterType
				}
				f = float64(vi)
			}
			res = append(res, f)
		}

		return res, nil
	}

	return nil, errParameterType
}

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
