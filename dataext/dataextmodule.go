package dataext

import "github.com/Stromberg/gel"

func init() {
	gel.RegisterModules(DataExtModule)
}

var DataExtModule = &gel.Module{
	Name: "dataext",
	LispFuncs: []*gel.LispFunc{
		&gel.LispFunc{Name: "dataext.Fix", F: "(func [d] (func [x] (if (vec? x) (vec-map (with-default d) x) ((with-default d) x))))"},
		&gel.LispFunc{Name: "dataext.FixPos", F: "(func [d] (func [x] (if (vec? x) (vec-map (positive d) x) ((positive d) x))))"},
	},
}
