package gel

import (
	"fmt"
	"reflect"
)

// SimpleFunc is a wrapper for a Gel function that works on specific data types
// by trying to find an exact argument match among the given funcs
func SimpleFunc(funcs ...interface{}) func(args []interface{}) (interface{}, error) {
	return func(args []interface{}) (interface{}, error) {
		var callFunc interface{}
		var argv []reflect.Value
		for _, f := range funcs {
			res := true
			v := reflect.ValueOf(f)
			t := v.Type()
			if len(args) != t.NumIn() {
				continue
			}
			argv = make([]reflect.Value, len(args))
			for i := range argv {
				argType := reflect.TypeOf(args[i]).Kind()

				if t.In(i).Kind() != reflect.Interface {
					if t.In(i).Kind() != argType {
						res = false
						break
					}
				}

				argv[i] = reflect.ValueOf(args[i])
			}

			if res {
				callFunc = f
				break
			}
		}

		if callFunc == nil {
			return nil, fmt.Errorf("Wrong Arguments")
		}

		v := reflect.ValueOf(callFunc)
		return v.Call(argv)[0].Interface(), nil
	}
}
