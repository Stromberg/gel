package main

import (
	"fmt"
	"os"

	"github.com/Stromberg/gel"
)

func main() {
	if len(os.Args) > 1 {
		g, err := gel.New(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		res, err := g.Eval(gel.NewEnv())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "%v\n", res)
		os.Exit(1)
		return
	}

	g, err := gel.New("")
	err = g.Repl(gel.NewEnv())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
