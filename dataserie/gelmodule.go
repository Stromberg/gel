package dataserie

import (
	"errors"
	"fmt"
	"time"

	"github.com/Stromberg/gel/module"
	"github.com/Stromberg/gel/utils"
)

var Module = &module.Module{
	Name: "dataserie",
	Funcs: []*module.Func{
		&module.Func{
			Name:        "dataserie?",
			Signature:   "(dataserie? ds)",
			Description: "Checks if dataserie",
			F: utils.SimpleFunc(func(ds interface{}) interface{} {
				_, ok := ds.(*DataSerie)
				return ok
			}, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "ds.Name",
			Signature:   "(dataserie.Name ds)",
			Description: "Name of ds",
			F: utils.SimpleFunc(func(ds *DataSerie) interface{} {
				return ds.Name
			}, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "ds.Name!",
			Signature:   "(ds.Name! name ds)",
			Description: "Set Name of ds",
			F: utils.SimpleFunc(func(name string, ds *DataSerie) interface{} {
				ds.Name = name
				return ds
			}, utils.CheckArity(2)),
		},
		&module.Func{
			Name:        "ds.Xs",
			Signature:   "(ds.Xs ds)",
			Description: "X values of ds",
			F: utils.SimpleFunc(func(ds *DataSerie) interface{} {
				xs := ds.Xs()
				res := make([]interface{}, len(xs))
				for i, x := range xs {
					res[i] = x
				}
				return res
			}, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "ds.Ys",
			Signature:   "(ds.Ys ds)",
			Description: "Y values of ds",
			F: utils.SimpleFunc(func(ds *DataSerie) interface{} {
				return ds.Ys()
			}, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "ds.LastY",
			Signature:   "(ds.LastY ds)",
			Description: "Last Y value of ds",
			F: utils.SimpleFunc(func(ds *DataSerie) interface{} {
				if len(ds.Data) == 0 {
					return nil
				}
				return ds.Data[len(ds.Data)-1].Y
			}, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "ds.FirstY",
			Signature:   "(ds.FirstY ds)",
			Description: "First Y value of ds",
			F: utils.SimpleFunc(func(ds *DataSerie) interface{} {
				if len(ds.Data) == 0 {
					return nil
				}
				return ds.Data[0].Y
			}, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "ds.New",
			Signature:   "(ds.New name xs ys)",
			Description: "Creates a new dataserie from name, xs and ys values",
			F: utils.ErrFunc(func(name string, xs []interface{}, ys []float64) (interface{}, error) {
				if len(xs) != len(ys) {
					return nil, fmt.Errorf("xs and ys must be same length: %v != %v", len(xs), len(ys))
				}

				xss := make([]string, len(xs))
				for i, x := range xs {
					v, ok := x.(string)
					if !ok {
						return nil, fmt.Errorf("Expected string, got %v", x)
					}
					xss[i] = v
				}

				return NewLine(name, ToPoints(xss, ys)), nil
			}, utils.CheckArity(3)),
		},
		&module.Func{
			Name:        "ds.After",
			Signature:   "(ds.After after ds)",
			Description: "DataSerie ds filtered to data after",
			F: utils.SimpleFunc(func(after string, ds *DataSerie) interface{} {
				return ds.After(after)
			}, utils.CheckArity(2)),
		},
		&module.Func{
			Name:        "ds.Before",
			Signature:   "(ds.Before before ds)",
			Description: "DataSerie ds filtered to data before ds",
			F: utils.SimpleFunc(func(before string, ds *DataSerie) interface{} {
				return ds.Before(before)
			}, utils.CheckArity(2)),
		},
		&module.Func{
			Name:        "ds.Map",
			Signature:   "(ds.Map name fn ds)",
			Description: "Applies fn to every element ds",
			F: utils.ErrFunc(func(fn func(...interface{}) (interface{}, error), ds *DataSerie) (interface{}, error) {
				xs := ds.Xs()
				ys := ds.Ys()

				resy := make([]float64, len(ys))
				for i, y := range ys {
					r, err := fn(y)
					if err != nil {
						return nil, err
					}

					v, ok := r.(float64)
					if !ok {
						return nil, fmt.Errorf("Expected float64 got %v", v)
					}
					resy[i] = v
				}
				return NewLine(ds.Name, ToPoints(xs, resy)), nil
			}, utils.CheckArity(2)),
		},
		&module.Func{
			Name:        "ds.Apply",
			Signature:   "(ds.Apply fn ds)",
			Description: "Applies fn to every element ds",
			F: utils.ErrFunc(func(fn func(...interface{}) (interface{}, error), ds *DataSerie) (interface{}, error) {
				xs := ds.Xs()
				ys := ds.Ys()

				res, err := fn(ys)
				if err != nil {
					return nil, err
				}

				resy, ok := res.([]float64)
				if !ok {
					return nil, errors.New("Expected []float64 result")
				}

				return NewLine(ds.Name, ToPoints(xs, resy)), nil
			}, utils.CheckArity(2)),
		},
		&module.Func{
			Name:        "ds.PairMap",
			Signature:   "(ds.PairMap name fn ds1 ds2)",
			Description: "Applies fn to every pair element of ds1 and ds2",
			F: utils.ErrFunc(func(name string, fn func(...interface{}) (interface{}, error), ds1, ds2 *DataSerie) (interface{}, error) {
				ds1, ds2 = ds1.Union(ds2)

				xs := ds1.Xs()
				ys1 := ds1.Ys()
				ys2 := ds2.Ys()

				resy := make([]float64, len(ys1))
				for i := range ys1 {
					r, err := fn(ys1[i], ys2[i])
					if err != nil {
						return nil, err
					}

					v, ok := r.(float64)
					if !ok {
						return nil, fmt.Errorf("Expected float64 got %v", v)
					}
					resy[i] = v
				}
				return NewLine(name, ToPoints(xs, resy)), nil
			}, utils.CheckArity(4)),
		},
		&module.Func{
			Name:        "ds.Lag",
			Signature:   "(dataserie.Lag ds n)",
			Description: "Shifts the data n steps forward in time",
			F: utils.ErrFunc(func(ds *DataSerie, n int) (interface{}, error) {
				if n < 0 {
					return nil, errors.New("Expected positive n")
				}
				if n > 0 && n > len(ds.Data) {
					return nil, errors.New("Cannot shift more than there is data")
				}
				return ds.Shift(n), nil
			}, utils.CheckArity(2), utils.ParamToInt(1)),
		},
		&module.Func{
			Name:        "ds.Lead",
			Signature:   "(ds.Lead ds n)",
			Description: "Shifts the data n steps backwards in time",
			F: utils.ErrFunc(func(ds *DataSerie, n int) (interface{}, error) {
				if n < 0 {
					return nil, errors.New("Expected positive n")
				}
				if n > 0 && n > len(ds.Data) {
					return nil, errors.New("Cannot shift more than there is data")
				}

				return ds.Shift(-n), nil
			}, utils.CheckArity(2), utils.ParamToInt(1)),
		},
		&module.Func{
			Name:        "ds.LagFill",
			Signature:   "(dataserie.LagFill ds n)",
			Description: "Shifts the data n steps forward in time by adding months at the end.",
			F: utils.ErrFunc(func(ds *DataSerie, n int) (interface{}, error) {
				next := func(s string) string {
					t, _ := time.Parse("2006-01-02", s)
					t = t.AddDate(0, 1, 0)
					return t.Format("2006-01-02")
				}
				return ds.LagFill(n, next), nil
			}, utils.CheckArity(2), utils.ParamToInt(1)),
		},
		&module.Func{
			Name:        "ds.PadLastUntil",
			Signature:   "(dataserie.PadLastUntil ds until)",
			Description: "Pads the data serie with new monthly data repeating the last value in Ys",
			F: utils.ErrFunc(func(ds *DataSerie, until string) (interface{}, error) {
				next := func(s string) string {
					t, _ := time.Parse("2006-01-02", s)
					t = t.AddDate(0, 1, 0)
					return t.Format("2006-01-02")
				}
				return ds.PadLastUntil(until, next), nil
			}, utils.CheckArity(2)),
		},
		&module.Func{
			Name:        "ds.Union!",
			Signature:   "(ds.Union! ds1 ds2 ...)",
			Description: "Ensures that all dataseries are over the same range",
			F: utils.SimpleFunc(func(dss ...*DataSerie) interface{} {
				UnionInPlace(dss...)
				return nil
			}, utils.CheckArityAtLeast(2)),
		},
	},
}
