package ops

import (
	"fmt"

	"github.com/sketch-hq/gql-lint/format"
	"github.com/sketch-hq/gql-lint/output"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [flags] fileA fileB",
	Short: "Find deprecated fields present in the first file but not in the second",
	Args:  ExactArgs(2, "You must specify two files to diff"),
	RunE:  diffCmdRun,
}

func init() {
	Program.AddCommand(diffCmd)
}

func diffCmdRun(cmd *cobra.Command, args []string) error {
	fileA, fileB := args[0], args[1]

	result, err := output.CompareFiles(fileA, fileB)
	if err != nil {
		return fmt.Errorf("unable to diff: %s", err)
	}

	r, err := format.DiffFormatter.Format(flags.outputFormat, result)
	if err == nil {
		fmt.Print(r)
	}
	return err
}
