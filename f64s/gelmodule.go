package f64s

import (
	"github.com/Stromberg/gel/module"
	"github.com/Stromberg/gel/utils"
)

var F64sModule = &module.Module{
	Name: "f64s",
	Funcs: []*module.Func{
		&module.Func{
			Name:        "f64s/RelChange",
			Description: "(f64s/RelChange v) calculates the relative change between values, 0 is appended to the beginning of the list for symmetry",
			F:           utils.SimpleFunc(RelChange, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/AbsChange",
			Description: "(f64s/AbsChange v) calculates the absolute change between values, 0 is appended to the beginning of the list for symmetry",
			F:           utils.SimpleFunc(AbsChange, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/RelChangeN",
			Description: "((f64s/RelChangeN n) v) calculates the relative change between values n apart",
			F: utils.SimpleFunc(func(n int) interface{} {
				return utils.SimpleFunc(RelChangeN(n), utils.CheckArity(1))
			}, utils.CheckArity(1), utils.ParamToInt(0)),
		},
		&module.Func{
			Name:        "f64s/AbsChangeN",
			Description: "((f64s/AbsChangeN n) v) calculates the absolute change between values n apart",
			F: utils.SimpleFunc(func(n int) interface{} {
				return utils.SimpleFunc(AbsChangeN(n), utils.CheckArity(1))
			}, utils.CheckArity(1), utils.ParamToInt(0)),
		},
		&module.Func{
			Name:        "f64s/AccumDev",
			Description: "(f64s/AccumDev v) calculates the accumulated development over a series of devs. Note that this function assumes that input is relative change.",
			F:           utils.SimpleFunc(AccumDev, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/Momentum",
			Description: "((f64s/Momentum n) v) calculates the relative change over a period",
			F: utils.SimpleFunc(func(n int) interface{} {
				return utils.SimpleFunc(Momentum(n), utils.CheckArity(1))
			}, utils.CheckArity(1), utils.ParamToInt(0)),
		},
		&module.Func{
			Name:        "f64s/Sma",
			Description: "((f64s/Sma n) v) calculates the simple moving average over a period",
			F: utils.SimpleFunc(func(n int) interface{} {
				return utils.SimpleFunc(Sma(n), utils.CheckArity(1))
			}, utils.CheckArity(1), utils.ParamToInt(0)),
		},
		&module.Func{
			Name:        "f64s/Pma",
			Description: "((f64s/Pma t) v) calculates the periodic moving average with threshold t",
			F: utils.SimpleFunc(func(threshold float64) interface{} {
				return utils.SimpleFunc(Pma(threshold), utils.CheckArity(1))
			}, utils.CheckArity(1), utils.ParamToFloat64(0)),
		},
		&module.Func{
			Name:        "f64s/StdevN",
			Description: "((f64s/StdevN n) v) calculates the standard deviation over period n",
			F: utils.SimpleFunc(func(n int) interface{} {
				return utils.SimpleFunc(StdevN(n), utils.CheckArity(1))
			}, utils.CheckArity(1), utils.ParamToInt(0)),
		},
		&module.Func{
			Name:        "f64s/CompositeMomentum",
			Description: "(f64s/CompositeMomentum v) calculates the composite momentum.",
			F:           utils.SimpleFunc(CompositeMomentum, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/ShortMomentum",
			Description: "(f64s/ShortMomentum v) calculates the composite momentum weighted to the short end.",
			F:           utils.SimpleFunc(ShortMomentum, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/MaxN",
			Description: "((f64s/MaxN n) v) calculates the maximum value over period n",
			F: utils.SimpleFunc(func(n int) interface{} {
				return utils.SimpleFunc(MaxN(n), utils.CheckArity(1))
			}, utils.CheckArity(1), utils.ParamToInt(0)),
		},
		&module.Func{
			Name:        "f64s/Stdev",
			Description: "(f64s.Stdev v) calculates the standard deviation of the vec.",
			F:           utils.SimpleFunc(Stdev, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/Mean",
			Description: "(f64s/Mean v) calculates the mean of the vec.",
			F:           utils.SimpleFunc(Mean, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/Sum",
			Description: "(f64s/Sum v) calculates the sum of the vec.",
			F:           utils.SimpleFunc(Sum, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/GeometricMeanDev",
			Description: "(f64s/GeometricMeanDev v) calculates the geometric mean of the vec.",
			F:           utils.SimpleFunc(GeometricMeanDev, utils.CheckArity(1)),
		},
		&module.Func{
			Name:        "f64s/Nrank",
			Description: "(f64s/Nrank v) calculates the standard deviation of the vec.",
			F:           utils.SimpleFunc(Nrank, utils.CheckArity(1)),
		},
	},
}
