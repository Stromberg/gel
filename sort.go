package gel

import "sort"

func Sort(s []interface{}, less func(v1, v2 interface{}) bool) (index []int) {
	slice := newSlice(s, less)
	sort.Sort(slice)
	return slice.idx
}

func SortIndex(s []interface{}, less func(v1, v2 interface{}) bool) (index []interface{}) {
	cpy := make([]interface{}, len(s))
	copy(cpy, s)
	slice := newSlice(cpy, less)
	sort.Sort(slice)

	res := make([]interface{}, len(slice.idx))
	for i, idx := range slice.idx {
		res[i] = int64(idx)
	}
	return res
}

type slice struct {
	values []interface{}
	idx    []int
	less   func(v1, v2 interface{}) bool
}

func (s slice) Len() int {
	return len(s.values)
}

func (s slice) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func (s slice) Less(i, j int) bool {
	return s.less(s.values[i], s.values[j])
}

func newSlice(n []interface{}, less func(v1, v2 interface{}) bool) *slice {
	s := &slice{values: n, idx: make([]int, len(n)), less: less}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}
