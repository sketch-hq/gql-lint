package ops

import (
	"encoding/json"
	"fmt"

	"github.com/sketch-hq/gql-lint/output"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [flags] fileA fileB",
	Short: "Find deprecated fields present in the first file but not in the second",
	Args:  diffCmdArgsValidation,
	RunE:  diffCmdRun,
}

func init() {
	Program.AddCommand(diffCmd)
}

func diffCmdArgsValidation(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("You must specify two files to diff")
	}
	return nil
}

func diffCmdRun(cmd *cobra.Command, args []string) error {
	fileA, fileB := args[0], args[1]

	result, err := output.CompareFiles(fileA, fileB)
	if err != nil {
		return fmt.Errorf("Unable to diff: %s", err)
	}

	switch outputFormat {
	case stdoutFormat:
		diffStdOut(fileA, fileB, result)
	case jsonFormat:
		err = diffJsonOut(result)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s is not a valid output format. Choose between json and stdout", outputFormat)
	}

	return nil
}

func diffStdOut(_ string, fileB string, out output.Data) {
	if len(out) == 0 {
		return
	}

	for _, f := range out {
		fmt.Printf("%s (%s)\n", f.Field, f.DeprecationReason)
		fmt.Printf("  %s:%d\n", f.File, f.Line)
	}
}

func diffJsonOut(out output.Data) error {
	bytes, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("Failed to encode json: %s", err)
	}

	fmt.Print(string(bytes))
	return nil
}
