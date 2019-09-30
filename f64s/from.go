package f64s

func Slice(src []float64) Iterable {
	return Iterable{
		Iterate: func() Iterator {
			index := 0
			len := len(src)

			return func() (item float64, ok bool) {
				ok = index < len
				if ok {
					item = src[index]
					index++
				}

				return
			}
		},
	}
}

func Repeat(value float64, len int) (res []float64) {
	res = make([]float64, len)
	for i := range res {
		res[i] = value
	}
	return
}

func Range(start, end, step float64) (res []float64) {
	if start < end {
		for v := start; v < end; v += step {
			res = append(res, v)
		}
	} else {
		for v := start; v > end; v += step {
			res = append(res, v)
		}
	}

	return res
}

func InfiniteRange(start, step float64) Iterable {
	return Iterable{
		Iterate: func() Iterator {
			v := start - step
			return func() (item float64, ok bool) {
				v = v + step
				ok = true
				item = v
				return
			}
		},
	}
}

func IntRange(start, end, step int) (res []float64) {
	if start < end {
		for v := start; v < end; v += step {
			res = append(res, float64(v))
		}
	} else {
		for v := start; v > end; v += step {
			res = append(res, float64(v))
		}
	}

	return res
}

func Ints(s []int) (res []float64) {
	res = make([]float64, len(s))
	for i, v := range s {
		res[i] = float64(v)
	}

	return res
}

// Empty is an iterable without values
var Empty = Iterable{
	Iterate: func() Iterator {
		return func() (item float64, ok bool) {
			ok = false
			return
		}
	},
}
