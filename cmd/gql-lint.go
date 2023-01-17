package main

import (
	"fmt"
	"os"

	"github.com/sketch-hq/gql-lint/cmd/ops"
)

func main() {
	os.Exit(run())
}

func run() int {
	if err := ops.Program.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}
