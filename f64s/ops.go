package f64s

import (
	"math"

	"github.com/mxmCherry/movavg"
)

// RelChange calculates the relative change between two value
// The first value gets a change of 0 for symmetry.
func RelChange(s []float64) (res []float64) {
	res = make([]float64, len(s))

	previousValue := 0.0
	for i, v := range s {
		if i == 0 {
			previousValue = v
			res[i] = 0.0
		} else {
			res[i] = (v - previousValue) / previousValue
			previousValue = v
		}
	}
	return
}

// RelChangeN calculates the relative change between two values n apart
func RelChangeN(n int) func(s []float64) (res []float64) {
	return func(s []float64) (res []float64) {
		if n > len(s) {
			return []float64{}
		}
		d := Skip(n)(s)
		res = make([]float64, len(d))

		for i := range d {
			res[i] = (d[i] - s[i]) / s[i]
		}
		return
	}
}

// AbsChange calculates the absolute change between two value
// The first value gets a change of 0 for symmetry.
func AbsChange(s []float64) (res []float64) {
	res = make([]float64, len(s))

	previousValue := 0.0
	for i, v := range s {
		if i == 0 {
			previousValue = v
			res[i] = 0.0
		} else {
			res[i] = (v - previousValue)
			previousValue = v
		}
	}
	return
}

// AbsChangeN calculates the absolute change between two values n apart
func AbsChangeN(n int) func(s []float64) (res []float64) {
	return func(s []float64) (res []float64) {
		if n > len(s) {
			return []float64{}
		}
		d := Skip(n)(s)
		res = make([]float64, len(d))

		for i := range d {
			res[i] = (d[i] - s[i])
		}
		return
	}
}

// AccumDev calculates the accumulated development over a series of devs.
// Note that this function assumes that input is relative change.
func AccumDev(s []float64) (res []float64) {
	res = make([]float64, len(s))

	for i, v := range s {
		if i == 0 {
			res[i] = 0.0
		} else {
			res[i] = (1+res[i-1])*(1+v) - 1
		}
	}
	return
}

// Momentum calculates the relative change over a period
func Momentum(n int) func(s []float64) []float64 {
	return func(s []float64) (res []float64) {
		f := Skip(n)(s)
		res = make([]float64, len(f))
		for i := range res {
			res[i] = f[i]/s[i] - 1
		}

		return
	}
}

// Sma calculates a simple moving average.
func Sma(window int) func(s []float64) []float64 {
	return func(s []float64) (res []float64) {
		adder := movavg.NewSMA(window)

		for i, v := range s {
			adder.Add(v)

			if i >= window-1 {
				res = append(res, adder.Avg())
			}
		}

		return
	}
}

// Periodic Moving Average.
// threshould should be 0.1 for equities, 0.035 for for bonds and 0.2 for metals.
func Pma(threshold float64) func(s []float64) []float64 {
	n := 12

	samplingRate := func(s []float64) (res []int) {
		vol := Pipe(
			RelChange,
			StdevN(12),
			Scale(math.Sqrt(12)),
		)(s)

		res = make([]int, len(vol))

		for i, v := range vol {
			if v < threshold {
				res[i] = 2
			} else {
				res[i] = 1
			}
		}

		return
	}

	return func(s []float64) (res []float64) {
		if len(s) < n {
			return []float64{}
		}

		sr := samplingRate(s)

		last := 0
		adder := movavg.NewSMA(6)
		adder2 := movavg.NewSMA(2)

		for i := range sr {
			last++
			if i == 0 || last >= sr[i] {
				v := s[i+n-1]
				adder.Add(v)
				last = 0
			}
			adder2.Add(adder.Avg())

			res = append(res, adder2.Avg())
		}

		return
	}
}

func StdevN(n int) func(s []float64) []float64 {
	return func(s []float64) []float64 {
		res := make([]float64, len(s)-n+1)
		for i := range res {
			res[i] = Stdev(s[i : i+n])
		}

		return res
	}
}

// AdjustedMean calculates a sliding window mean weighted towards more recent values.
func AdjustedMean(n int) func(s []float64) []float64 {
	return func(s []float64) []float64 {
		res := make([]float64, len(s)-n+1)
		for i := range res {
			v := 0.0
			for j := 0; j < n; j++ {
				v += s[i+j] * float64(j+1)
			}
			res[i] = v / float64(n*(n+1)/2)
		}

		return res
	}
}

// CompositeMomentum calculates a mean momentum for periods 3, 6, 9 and 12.
func CompositeMomentum(s []float64) (res []float64) {
	v12 := Momentum(12)(s)
	v9 := Pipe(Skip(3), Momentum(9))(s)
	v6 := Pipe(Skip(6), Momentum(6))(s)
	v3 := Pipe(Skip(9), Momentum(3))(s)

	res = make([]float64, len(v12))

	for i := range res {
		res[i] = (v12[i] + v9[i] + v6[i] + v3[i]) / 4
	}

	return
}

// ShortMomentum calculates a mean momentum for periods 1, 3 and 6.
func ShortMomentum(s []float64) (res []float64) {
	v6 := Momentum(6)(s)
	v3 := Pipe(Skip(3), Momentum(3))(s)

	res = make([]float64, len(v3))

	for i := range res {
		// Scaling to be similar in scale to CompositeMomentum
		res[i] = math.Pow(1+(v6[i]+v3[i])/2.0, 2) - 1
	}

	return
}

// Max calculates a moving max
func MaxN(window int) func(s []float64) []float64 {
	return func(s []float64) []float64 {
		max := make([]float64, 0)
		i := 0
		for _, v := range s {
			if i < window {
				for j := 0; j < i; j++ {
					v = math.Max(v, max[j])
				}
			} else {
				for j := i - window; j < i; j++ {
					v = math.Max(v, max[j])
				}
			}
			max = append(max, v)
		}

		return max
	}
}
