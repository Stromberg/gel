package gel

import (
	"fmt"
	"strings"
)

func init() {
	RegisterModule(StdLibModule)
}

var StdLibModule = &Module{
	Name: "stdlib",
	Funcs: []*Func{
		// &Func{Name: "strings.Join", F: SimpleFunc(strings.Join, CheckArity(2))},
		// &Func{Name: "strings.Split", F: SimpleFunc(strings.Split, CheckArity(2))},
		&Func{Name: "strings.Title", F: SimpleFunc(strings.Title, CheckArity(1))},
		&Func{Name: "strings.ToLower", F: SimpleFunc(strings.ToLower, CheckArity(1))},
		&Func{Name: "strings.ToUpper", F: SimpleFunc(strings.ToUpper, CheckArity(1))},
		&Func{Name: "strings.TrimSpace", F: SimpleFunc(strings.TrimSpace, CheckArity(1))},
		&Func{Name: "printf", F: ErrFunc(func(args ...interface{}) (int, error) {
			format := args[0].(string)
			return fmt.Printf(format, args[1:]...)
		}, CheckArityAtLeast(1))},
		&Func{Name: "sprintf", F: SimpleFunc(func(args ...interface{}) string {
			format := args[0].(string)
			return fmt.Sprintf(format, args[1:]...)
		}, CheckArityAtLeast(1))},
	},
}
