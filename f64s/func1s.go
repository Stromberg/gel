package f64s

import "math"

// Map creates a pipeable function that maps an iterable to another
func Map(f func(v1 float64) float64) func(s []float64) []float64 {
	return func(s []float64) (res []float64) {
		res = make([]float64, len(s))

		for i, v := range s {
			res[i] = f(v)
		}
		return
	}
}

// Scale multiplies a value
func Scale(s float64) func([]float64) []float64 {
	return Map(func(v float64) float64 {
		return s * v
	})
}

// Cap limits values to be in a range
func Cap(lower float64, upper float64) func(float64) float64 {
	return func(v float64) float64 {
		if v < lower {
			return lower
		}

		if v > upper {
			return upper
		}
		return v
	}
}

// Add adds a constant value
func Add(s float64) func([]float64) []float64 {
	return Map(func(v float64) float64 {
		return s + v
	})
}

// Pow calculates the n exponential of a value
func Pow(n float64) func([]float64) []float64 {
	return Map(func(v float64) float64 {
		return math.Pow(v, n)
	})
}

// Default returns default value if value is not valid
func Default(defaultValue float64) func(float64) float64 {
	return func(v float64) float64 {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return defaultValue
		}
		return v
	}
}

// Positive returns default value if value is not valid or is negativ
func Positive(defaultValue float64) func(float64) float64 {
	return func(v float64) float64 {
		if math.IsNaN(v) || math.IsInf(v, 0) || v < 0.0 {
			return defaultValue
		}
		return v
	}
}
