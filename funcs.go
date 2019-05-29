package gel

import (
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
		return result[0].Interface(), result[1].Interface().(error)
	}
}
