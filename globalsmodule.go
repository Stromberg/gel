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
		&Func{Name: "range", F: rangeFn},
		&Func{Name: "vec", F: sliceF64sFn},
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

func sliceF64sFn(args ...interface{}) (value interface{}, err error) {
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

func rangeFn(scope *Scope, args []ast.Node) (value interface{}, err error) {
	if len(args) < 3 {
		return nil, errors.New(`range takes three or more arguments`)
	}
	var iname, ename string
	if symbol, ok := args[0].(*ast.Symbol); ok {
		iname = symbol.Name
	} else if list, ok := args[0].(*ast.List); ok && len(list.Nodes) == 2 {
		symbol1, ok1 := list.Nodes[0].(*ast.Symbol)
		symbol2, ok2 := list.Nodes[1].(*ast.Symbol)
		if ok1 && ok2 {
			iname = symbol1.Name
			ename = symbol2.Name
		}
	}
	if iname == "" {
		return nil, errors.New(`range takes var name or (i elem) var name pair as first argument`)
	}
	scope = scope.Branch()
	value, err = scope.Eval(args[1])
	if err != nil {
		return nil, err
	}
	code := args[2:]
	if n, ok := value.(int64); ok {
		scope.Create(iname, 0)
		for i := int64(0); i < n; i++ {
			scope.Set(iname, i)
			for _, c := range code {
				value, err = scope.Eval(c)
				if err != nil {
					return nil, err
				}
			}
		}
		return value, nil
	}
	if list, ok := value.([]interface{}); ok {
		scope.Create(iname, 0)
		scope.Create(ename, nil)
		for i, e := range list {
			scope.Set(iname, i)
			scope.Set(ename, e)
			for _, c := range code {
				value, err = scope.Eval(c)
				if err != nil {
					return nil, err
				}
			}
		}
		return value, nil
	}
	return nil, errors.New(`range takes an integer or a list as second argument`)
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
