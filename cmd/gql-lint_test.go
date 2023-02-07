package main

import (
	"flag"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/google/go-cmdtest"
)

// Run tests with `UPDATE_CMD_TESTS=true make test` to re-generate the cmd test files
var updateEnv = os.Getenv("UPDATE_CMD_TESTS")
var update = flag.Bool("update", updateEnv == "true", "update test files with results")

func TestCLI(t *testing.T) {
	ts, err := cmdtest.Read("testdata")
	if err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	ts.Setup = func(rootDir string) error {
		// copy over the testdata files as we want the paths relative to the
		// temp root dir. If not paths will be change depending on the
		// computer and location of the project.
		cp := exec.Command("cp", "-rf", path.Join(wd, "testdata"), rootDir)
		err := cp.Run()
		return err
	}

	ts.Commands["gql-lint"] = cmdtest.InProcessProgram("gql-lint", run)
	ts.Run(t, *update)
}
