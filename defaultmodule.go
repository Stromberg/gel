package gel

// DefaultModule contains some basic functions that
// are needed in many Gel expressions
var DefaultModule = &Module{
	Name:        "",
	Description: "",
	Funcs: []*Func{
		&Func{
			Name:        "<=",
			Signature:   "",
			Description: "",
			F: SimpleFunc(func(v1, v2 interface{}) bool {
				if f1, ok := v1.(float64); ok {
					if f2, ok := v2.(float64); ok {
						return f1 <= f2
					}

					if f2, ok := v2.(int64); ok {
						return f1 <= float64(f2)
					}
				}

				if f2, ok := v2.(float64); ok {
					return float64(v1.(int64)) <= f2
				}

				return v1.(int64) <= v2.(int64)
			}),
		},
		&Func{
			Name:        ">=",
			Signature:   "",
			Description: "",
			F: SimpleFunc(func(v1, v2 interface{}) bool {
				if f1, ok := v1.(float64); ok {
					if f2, ok := v2.(float64); ok {
						return f1 >= f2
					}

					if f2, ok := v2.(int64); ok {
						return f1 >= float64(f2)
					}
				}

				if f2, ok := v2.(float64); ok {
					return float64(v1.(int64)) >= f2
				}

				return v1.(int64) >= v2.(int64)
			}),
		},
		&Func{
			Name:        "<",
			Signature:   "",
			Description: "",
			F: SimpleFunc(func(v1, v2 interface{}) bool {
				if f1, ok := v1.(float64); ok {
					if f2, ok := v2.(float64); ok {
						return f1 < f2
					}

					if f2, ok := v2.(int64); ok {
						return f1 < float64(f2)
					}
				}

				if f2, ok := v2.(float64); ok {
					return float64(v1.(int64)) < f2
				}

				return v1.(int64) < v2.(int64)
			}),
		},
		&Func{
			Name:        ">",
			Signature:   "",
			Description: "",
			F: SimpleFunc(func(v1, v2 interface{}) bool {
				if f1, ok := v1.(float64); ok {
					if f2, ok := v2.(float64); ok {
						return f1 > f2
					}

					if f2, ok := v2.(int64); ok {
						return f1 > float64(f2)
					}
				}

				if f2, ok := v2.(float64); ok {
					return float64(v1.(int64)) > f2
				}

				return v1.(int64) > v2.(int64)
			}),
		},
	},
}
