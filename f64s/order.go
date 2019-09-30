package f64s

import "sort"

func Sort(s []float64) (index []int) {
	slice := newSlice(sort.Float64Slice(s))
	sort.Sort(slice)
	return slice.idx
}

type floatSlice struct {
	values []float64
	idx    []int
}

func (s floatSlice) Len() int {
	return len(s.values)
}

func (s floatSlice) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func (s floatSlice) Less(i, j int) bool {
	return s.values[i] < s.values[j]
}

func newSlice(n []float64) *floatSlice {
	s := &floatSlice{values: n, idx: make([]int, len(n))}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

// Nrank returns the normalized rank of the elements in a list in ascending order.
func Nrank(values []float64) []float64 {
	cvalues := make([]float64, len(values))
	copy(cvalues, values)

	if len(cvalues) == 0 {
		return nil
	}

	if len(cvalues) == 1 {
		return []float64{1.0}
	}

	s := newSlice(cvalues)
	sort.Sort(s)

	max := float64(len(s.idx)) - 1
	scale := 1.0 / max

	res := make([]float64, len(s.idx))
	for i, id := range s.idx {
		res[id] = float64(i) * scale
	}

	return res
}
