package f64s

func Take(count int) func(s []float64) []float64 {
	return func(s []float64) []float64 {
		l := count
		if l > len(s) {
			l = len(s)
		}

		if l == 0 {
			return nil
		}

		return s[0:l]
	}
}
