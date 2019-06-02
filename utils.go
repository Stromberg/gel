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

func IsAnySlice(vs ...interface{}) bool {
	for _, v := range vs {
		if IsSlice(v) {
			return true
		}
	}

	return false
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
