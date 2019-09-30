package f64s

import (
	"math"
)

type Iterator func() (value float64, done bool)

type Iterable struct {
	Iterate func() Iterator
}

// Slice creates a float64 slice of the data
func (s Iterable) Slice() []float64 {
	var res []float64

	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		res = append(res, item)
	}

	return res
}

func First(s []float64) float64 {
	if len(s) > 0 {
		return s[0]
	}
	return math.NaN()
}

func Last(s []float64) float64 {
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return math.NaN()
}

func CountIf(test func(float64) bool) func(s []float64) int {
	return func(s []float64) (res int) {
		res = 0
		for _, v := range s {
			if test(v) {
				res++
			}
		}
		return res
	}
}

func Mean(s []float64) float64 {
	return Sum(s) / float64(len(s))
}

func Max(s []float64) float64 {
	if len(s) == 0 {
		return math.NaN()
	}

	max := s[0]
	for _, v := range s {
		max = math.Max(max, v)
	}
	return max
}

func Min(s []float64) float64 {
	if len(s) == 0 {
		return math.NaN()
	}

	min := s[0]
	for _, v := range s {
		min = math.Min(min, v)
	}
	return min
}

func Sum(s []float64) float64 {
	sum := 0.0

	for _, v := range s {
		sum += v
	}
	return sum
}

func Stdev(s []float64) float64 {
	mean := Mean(s)
	sum := 0.0

	for _, v := range s {
		sum += math.Pow(v-mean, 2)
	}

	return math.Sqrt(sum / float64(len(s)))
}

// GeometricMeanDev calculates the geometric mean of the sequence.
func GeometricMeanDev(s []float64) float64 {
	prod := 1.0

	for _, v := range s {
		prod *= (1.0 + v)
	}

	return math.Pow(prod, 1.0/float64(len(s))) - 1
}

func Normalize(s []float64) []float64 {
	res := make([]float64, len(s))
	min := Min(s)
	max := Max(s)

	for i, v := range s {
		res[i] = (v - min) / (max - min)
	}

	return res
}
