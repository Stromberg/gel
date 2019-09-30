package dataserie

import "github.com/Stromberg/gel/utils"

type DataPoint struct {
	X string  `json:"x"`
	Y float64 `json:"y"`
}

type DataSerie struct {
	Name       string      `json:"name"`
	Data       []DataPoint `json:"data"`
	Type       string      `json:"type"`
	Mode       string      `json:"mode"`
	StackGroup string      `json:"stackgroup"`
	Fill       string      `json:"fill"`
}

type DateRange struct {
	Start string
	End   string
}

func Xs(data []DataPoint) []string {
	names := make([]string, len(data))
	for i := range names {
		names[i] = data[i].X
	}

	return names
}

func Ys(data []DataPoint) (res []float64) {
	res = make([]float64, len(data))
	for i, dp := range data {
		res[i] = dp.Y
	}
	return
}

func After(data []DataPoint, after string) []DataPoint {
	l := len(data)
	if l == 0 {
		return nil
	}

	for i := range data {
		if data[i].X >= after {
			return data[i:l]
		}
	}

	return nil
}

func Before(data []DataPoint, before string) []DataPoint {
	l := len(data)
	if l == 0 {
		return nil
	}

	for i := range data {
		if data[i].X == before {
			return data[0 : i+1]
		} else if data[i].X > before {
			return data[0:i]
		}
	}

	return nil
}

func Range(d []DataPoint) DateRange {
	return DateRange{
		Start: d[0].X,
		End:   d[len(d)-1].X,
	}
}

func CommonRange(r1, r2 DateRange) DateRange {
	r := r2

	if r1.Start > r.Start {
		r.Start = r1.Start
	}

	if r1.End < r.End {
		r.End = r1.End
	}

	return r
}

func Sub(r DateRange, d []DataPoint) []DataPoint {
	return Before(After(d, r.Start), r.End)
}

func (serie *DataSerie) After(after string) *DataSerie {
	return &DataSerie{
		Name: serie.Name,
		Data: After(serie.Data, after),
		Type: serie.Type,
		Mode: serie.Mode,
	}
}

func (serie *DataSerie) Before(after string) *DataSerie {
	return &DataSerie{
		Name: serie.Name,
		Data: Before(serie.Data, after),
		Type: serie.Type,
		Mode: serie.Mode,
	}
}

func (serie *DataSerie) Xs() []string {
	return Xs(serie.Data)
}

func (serie *DataSerie) Ys() []float64 {
	return Ys(serie.Data)
}

func (serie *DataSerie) YsAfter(after string) []float64 {
	return Ys(After(serie.Data, after))
}

func (serie *DataSerie) Range() DateRange {
	return Range(serie.Data)
}

func (serie *DataSerie) Sub(r DateRange) *DataSerie {
	return &DataSerie{
		Name: serie.Name,
		Data: Sub(r, serie.Data),
		Type: serie.Type,
		Mode: serie.Mode,
	}
}

func (serie *DataSerie) Copy() *DataSerie {
	res := make([]DataPoint, len(serie.Data))
	copy(res, serie.Data)
	return &DataSerie{
		Name: serie.Name,
		Data: res,
		Type: serie.Type,
		Mode: serie.Mode,
	}
}

func (ds *DataSerie) Shift(n int) *DataSerie {
	if (n > 0 && n > len(ds.Data)) ||
		(n < 0 && -n > len(ds.Data)) {
		return nil
	}

	xs := ds.Xs()
	ys := ds.Ys()

	if n > 0 {
		xs = xs[n:]
		ys = ys[:len(ys)-n]
	} else if n < 0 {
		xs = xs[:len(xs)+n]
		ys = ys[-n:]
	}

	return NewLine(ds.Name, ToPoints(xs, ys))
}

// Ugly hack
func (ds *DataSerie) PadLastUntil(end string, next func(s string) string) *DataSerie {

	xs := ds.Xs()
	ys := ds.Ys()
	ysLast := ys[len(ys)-1]

	for xs[len(xs)-1] != end {
		n := next(xs[len(xs)-1])
		xs = append(xs, n)
		ys = append(ys, ysLast)
	}

	return NewLine(ds.Name, ToPoints(xs, ys))
}

func (s1 *DataSerie) Union(s2 *DataSerie) (res1, res2 *DataSerie) {
	r := CommonRange(s1.Range(), s2.Range())
	return s1.Sub(r), s2.Sub(r)
}

func UnionInPlace(dss ...*DataSerie) {
	cr := dss[0].Range()

	for _, ds := range dss {
		cr = CommonRange(cr, ds.Range())
	}

	for _, ds := range dss {
		r := ds.Sub(cr)
		ds.Data = r.Data
	}
}

func (s *DataSerie) Map(name string, op func(float64) float64) *DataSerie {
	return NewLine(name, ToPoints(s.Xs(), utils.MapVec(op, s.Ys())))
}

func (s *DataSerie) Apply(op func([]float64) []float64) *DataSerie {
	return NewLine(s.Name, ToPoints(s.Xs(), op(s.Ys())))
}

func ToPoints(xs []string, ys []float64) []DataPoint {
	points := make([]DataPoint, len(ys))
	offset := len(xs) - len(ys)
	for i := range points {
		points[i] = DataPoint{xs[i+offset], ys[i]}
	}

	return points
}

func NewLine(name string, points []DataPoint) *DataSerie {
	return &DataSerie{
		Name: name,
		Data: points,
		Type: "scatter",
		Mode: "lines+points",
	}
}
