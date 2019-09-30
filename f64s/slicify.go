package f64s

import "reflect"

// Slicify takes a function that takes float64, int64 and int arguments
// and returns a function that accepts iterable arguments
func Slicify(f interface{}) func(args ...Iterable) Iterable {
	vf := reflect.ValueOf(f)
	t := vf.Type()
	types := make([]reflect.Type, t.NumIn())
	for i := range types {
		types[i] = t.In(i)
	}

	argv := make([]reflect.Value, len(types))

	toValue := func(v float64, t reflect.Type) reflect.Value {
		switch t.Kind() {
		case reflect.Float64:
			return reflect.ValueOf(v)
		case reflect.Int64:
			return reflect.ValueOf(int64(v))
		case reflect.Int:
			return reflect.ValueOf(int(v))
		}
		return reflect.ValueOf(0)
	}

	return func(args ...Iterable) Iterable {
		if len(args) != len(types) {
			return Empty
		}

		nexts := make([]Iterator, len(args))
		for i, arg := range args {
			nexts[i] = arg.Iterate()
		}

		return Iterable{
			Iterate: func() Iterator {
				return func() (item float64, ok bool) {
					for i, next := range nexts {
						v, ok1 := next()
						if !ok1 {
							ok = false
							return
						}
						argv[i] = toValue(v, types[i])
					}
					item = vf.Call(argv)[0].Float()
					ok = true
					return
				}
			},
		}
	}
}
