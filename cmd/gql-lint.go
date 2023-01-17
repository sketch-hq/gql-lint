package main

import (
	"os"

	"github.com/sketch-hq/gql-lint/cmd/ops"
)

func main() {
	os.Exit(run())
}

// run is a wrapper to allow using go-cmdtest to test the CLI inProcess
func run() int {
	err := ops.Program.Execute()
	if err != nil {
		//Error is already printed by cobra
		return 1
	} else {
		return 0
	}
}
