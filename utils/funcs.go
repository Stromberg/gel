package utils

import (
	"errors"
	"reflect"
)

// SimpleFunc builds a Gel function from a function that does not return an error
func SimpleFunc(v interface{}, adapters ...Adapter) interface{} {
	return func(values ...interface{}) (interface{}, error) {
		args := values
		var err error
		for _, adapter := range adapters {
			args, err = adapter(args...)
			if err != nil {
				return nil, err
			}
		}

		vargs := []reflect.Value{}
		for _, arg := range args {
			vargs = append(vargs, reflect.ValueOf(arg))
		}

		result := reflect.ValueOf(v).Call(vargs)
		return result[0].Interface(), nil
	}
}

// ErrFunc builds a Gel function from a function that returns an error
func ErrFunc(v interface{}, adapters ...Adapter) interface{} {
	return func(values ...interface{}) (interface{}, error) {
		args := values
		var err error
		for _, adapter := range adapters {
			args, err = adapter(args...)
			if err != nil {
				return nil, err
			}
		}

		vargs := []reflect.Value{}
		for _, arg := range args {
			vargs = append(vargs, reflect.ValueOf(arg))
		}

		result := reflect.ValueOf(v).Call(vargs)
		err = nil
		if result[1].Interface() != nil {
			err = result[1].Interface().(error)
		}
		return result[0].Interface(), err
	}
}

var GetFn = ErrFunc(func(args ...interface{}) (interface{}, error) {
	switch arg := args[0].(type) {
	case map[interface{}]interface{}:
		v, ok := arg[args[1]]
		if !ok {
			return false, nil
		}
		return v, nil
	case []interface{}:
		v := arg
		i, ok := args[1].(int64)
		if !ok {
			return nil, ErrParameterType
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
			return nil, ErrParameterType
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
	return nil, ErrParameterType
}, CheckArity(2))
