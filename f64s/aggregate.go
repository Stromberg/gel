package f64s

func Aggregate(
	seed func(float64) (value float64, state interface{}),
	f func(value float64, previousState interface{}) (newValue float64, state interface{}, done bool)) func([]float64) []float64 {
	return func(s []float64) (res []float64) {
		prev := 0.0
		var state interface{}
		done := false
		for i, v := range s {
			if i == 0 {
				prev, state = seed(v)

			} else {
				prev, state, done = f(v, state)
			}

			if !done {
				res = append(res, prev)
			} else {
				return
			}
		}

		return
	}
}

func Aggregate2(
	seed func(float64) (value float64),
	f func(value float64, previous float64, step int) (newValue float64, done bool)) func([]float64) []float64 {
	return func(s []float64) (res []float64) {
		prev := 0.0
		done := false

		for i, v := range s {
			if i == 0 {
				prev = seed(v)
			} else {
				prev, done = f(v, prev, i)
			}

			if !done {
				res = append(res, prev)
			} else {
				return
			}
		}

		return
	}
}
