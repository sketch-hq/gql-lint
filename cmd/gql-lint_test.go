package main

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmdtest"
)

// set this to true if you want to update the test cases in `testdata/`
var update = flag.Bool("update", false, "update test files with results")

func TestCLI(t *testing.T) {
	ts, err := cmdtest.Read("testdata")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("VALID_SCHEMA_FILE", "parser/testdata/schemas/with_deprecations.gql"); err != nil {
		t.Fatal(err)
	}

	ts.Commands["gql-lint"] = cmdtest.InProcessProgram("gql-lint", run)
	ts.Run(t, *update)
}
