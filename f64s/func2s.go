package f64s

// Op2 creates a function that can be pairwise applied to two iterables
func Op2(f func(v1, v2 float64) float64) func(s1, s2 []float64) []float64 {
	return func(s1, s2 []float64) (res []float64) {
		l := len(s1)
		if l > len(s2) {
			l = len(s2)
		}

		res = make([]float64, l)
		for i := 0; i < l; i++ {
			res[i] = f(s1[i], s2[i])
		}

		return
	}
}

// Mul returns the product of two numbers
var Mul = func(v1, v2 float64) float64 { return v1 * v2 }
