package main

import (
	"fmt"
	"os"

	"github.com/Stromberg/gel"
)

func main() {
	g, err := gel.New("")
	err = g.Repl(gel.NewEnv())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
