package f64s

func Pipe(funcs ...func([]float64) []float64) func([]float64) []float64 {
	return func(s []float64) []float64 {
		for _, f := range funcs {
			s = f(s)
		}
		return s
	}
}
