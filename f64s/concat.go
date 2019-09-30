package f64s

func Concat(slices ...[]float64) (res []float64) {
	for _, s := range slices {
		res = append(res, s...)
	}

	return
}
