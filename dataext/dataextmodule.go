package dataext

import "github.com/Stromberg/gel/module"

func init() {
	module.RegisterModules(DataExtModule)
}

var DataExtModule = &module.Module{
	Name: "dataext",
	LispFuncs: []*module.LispFunc{
		&module.LispFunc{Name: "dataext.Fix", F: "(func [d] (func [x] (if (vec? x) (vec-map (with-default d) x) ((with-default d) x))))"},
		&module.LispFunc{Name: "dataext.FixPos", F: "(func [d] (func [x] (if (vec? x) (vec-map (positive d) x) ((positive d) x))))"},
	},
}
