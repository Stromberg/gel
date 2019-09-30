package f64s

func Skip(n int) func(s []float64) []float64 {
	return func(s []float64) []float64 {
		if len(s) <= n {
			return nil
		}

		return s[n:len(s)]
	}
}
