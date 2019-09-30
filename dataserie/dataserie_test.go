package dataserie

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShift(t *testing.T) {
	ds := NewLine("", ToPoints([]string{"0", "1", "2"}, []float64{12, 13, 14}))
	ds2 := ds.Shift(1)
	assert.Equal(t, "1", ds2.Range().Start)
}

func TestPadLastUntil(t *testing.T) {
	ds := NewLine("", ToPoints([]string{"0", "1", "2"}, []float64{12, 13, 14}))
	next := func(s string) string {
		i, _ := strconv.ParseInt(s, 0, 64)
		return fmt.Sprintf("%v", i+1)
	}
	ds2 := ds.PadLastUntil("5", next)
	assert.Equal(t, []string{"0", "1", "2", "3", "4", "5"}, ds2.Xs())
	assert.Equal(t, []float64{12, 13, 14, 14, 14, 14}, ds2.Ys())
}
