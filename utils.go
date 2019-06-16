package gel

import (
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
