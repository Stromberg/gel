package gel

import (
	"errors"
	"reflect"
)

func IsAnyStrings(values ...interface{}) bool {
	for _, v := range values {
		_, ok := v.(string)
		if ok {
			return true
		}
	}
	return false
}

func IsAllStrings(values ...interface{}) bool {
	for _, v := range values {
		_, ok := v.(string)
		if !ok {
			return false
		}
	}
	return true
}

func IsAnyFloats(values ...interface{}) bool {
	for _, v := range values {
		_, ok := v.(float64)
		if ok {
			return true
		}
	}
	return false
}

func IsSlice(v interface{}) bool {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		return true
	default:
		return false
	}
}

func CopySlice(v interface{}) (interface{}, error) {
	switch v.(type) {
	case []interface{}:
		res := make([]interface{}, len(v.([]interface{})))
		copy(res, v.([]interface{}))
		return res, nil
	case []float64:
		res := make([]float64, len(v.([]float64)))
		copy(res, v.([]float64))
		return res, nil
	default:
		return nil, errParameterType
	}
}

func IsAnySlice(vs ...interface{}) bool {
	for _, v := range vs {
		if IsSlice(v) {
			return true
		}
	}

	return false
}

func VecToList(vec []float64) []interface{} {
	res := make([]interface{}, len(vec))
	for i, v := range vec {
		res[i] = v
	}
	return res
}

func ToList(data interface{}) (res []interface{}, ok bool) {
	switch arg := data.(type) {
	case []interface{}:
		return arg, true
	case []float64:
		return VecToList(arg), true
	}

	return nil, false
}

func NewDict(args ...interface{}) (value interface{}, err error) {
	if len(args)%2 != 0 {
		return nil, errors.New("dict requires an even number of arguments")
	}

	if len(args) == 0 {
		return map[interface{}]interface{}{}, nil
	}

	res := make(map[interface{}]interface{})
	for i := 0; i+1 < len(args); i += 2 {
		res[args[i]] = args[i+1]
	}

	return res, nil
}

func NewList(args ...interface{}) (value interface{}, err error) {
	if len(args) == 0 {
		return []interface{}{}, nil
	}

	res := make([]interface{}, len(args))
	for i, arg := range args {
		res[i] = arg
	}

	return res, nil
}

func MakeAllEitherSliceOrValue(vs ...interface{}) ([]interface{}, error) {
	if !IsAnySlice(vs...) {
		return vs, nil
	}

	l := 1
	for _, v := range vs {
		if IsSlice(v) {
			s := reflect.ValueOf(v)
			l = s.Len()
			break
		}
	}

	res := make([]interface{}, len(vs))
	for i, v := range vs {
		if IsSlice(v) {
			res[i] = v
		} else {
			r := make([]float64, l)
			switch v.(type) {
			case float64:
				for j := range r {
					r[j] = v.(float64)
				}
			case int64:
				for j := range r {
					r[j] = float64(v.(int64))
				}
			}

			res[i] = r
		}
	}

	return res, nil
}
