package gel

import (
	"errors"
)

var (
	errWrongNumberPar = errors.New("Wrong number of parameters")
	errParameterType  = errors.New("Error in parameter type")
)

type Adapter func(values ...interface{}) ([]interface{}, error)

//CheckArity return an function to check if the arity of function call is n
func CheckArity(n int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		if len(values) != n {
			return []interface{}{}, errWrongNumberPar
		}
		return values, nil
	}
}

//CheckArityAtLeast return an function to check if the arity of function call is at least n
func CheckArityAtLeast(n int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		if len(values) < n {
			return []interface{}{}, errWrongNumberPar
		}
		return values, nil
	}
}

//CheckArityEven return an function to check if the arity of function call is even
func CheckArityEven() Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		if len(values)%2 != 0 {
			return []interface{}{}, errWrongNumberPar
		}
		return values, nil
	}
}

//CheckArityOdd return an function to check if the arity of function call is even
func CheckArityOdd() Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		if len(values)%2 != 1 {
			return []interface{}{}, errWrongNumberPar
		}
		return values, nil
	}
}

func ParamsToSameType() Adapter {
	anyStrings := func(values ...interface{}) bool {
		for _, v := range values {
			_, ok := v.(string)
			if ok {
				return true
			}
		}
		return false
	}

	allStrings := func(values ...interface{}) bool {
		for _, v := range values {
			_, ok := v.(string)
			if !ok {
				return false
			}
		}
		return true
	}

	anyFloats := func(values ...interface{}) bool {
		for _, v := range values {
			_, ok := v.(float64)
			if ok {
				return true
			}
		}
		return false
	}

	return func(values ...interface{}) ([]interface{}, error) {
		if anyStrings(values...) {
			if !allStrings(values) {
				return []interface{}{}, errParameterType
			}
		} else if anyFloats(values...) {
			for i := range values {
				switch values[i].(type) {
				case int64:
					v := values[i].(int64)
					values[i] = float64(v)
				}
			}
		}

		return values, nil
	}
}

//ParamToFloat64 return an Adapter to convert the p nth param to float64 type
func ParamToFloat64(p int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		switch values[p].(type) {
		case int64:
			v := values[p].(int64)
			values[p] = float64(v)
			return values, nil
		case float64:
			return values, nil
		default:
			return []interface{}{}, errParameterType
		}
	}
}

//ParamToInt64 return an Adapter to convert the p nth param to int64 type
func ParamToInt64(p int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		switch values[p].(type) {
		case int64:
			return values, nil
		case float64:
			v := values[p].(float64)
			values[p] = int64(v)
			return values, nil
		default:
			return []interface{}{}, errParameterType
		}
	}
}

//ParamToInt return an Adapter to convert the p nth param to int64 type
func ParamToInt(p int) Adapter {
	return func(values ...interface{}) ([]interface{}, error) {
		switch values[p].(type) {
		case int:
			return values, nil
		case int64:
			v := values[p].(float64)
			values[p] = int(v)
			return values, nil
		case float64:
			v := values[p].(float64)
			values[p] = int64(v)
			return values, nil
		default:
			return []interface{}{}, errParameterType
		}
	}
}
