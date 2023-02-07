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
		return fmt.Errorf("Unable to diff: %s", err)
	}

	switch flags.outputFormat {
	case stdoutFormat:
		diffStdOut(fileA, fileB, result)
	case jsonFormat:
		err = diffJsonOut(result)
		if err != nil {
			return err
		}
	case xcodeFormat:
		diffXcodeOut(result)
	default:
		return fmt.Errorf("%s is not a valid output format. Choose between json and stdout", flags.outputFormat)
	}

	return nil
}

func diffStdOut(_ string, _ string, out output.Data) {
	out.Walk(func(schema string, f output.Field, i int) {
		if i == 0 {
			fmt.Println("Schema: ", schema)
		}
		fmt.Printf("%s (%s)\n", f.Field, f.DeprecationReason)
		fmt.Printf("  %s:%d\n", f.File, f.Line)
	})
}

func diffJsonOut(out output.Data) error {
	bytes, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("Failed to encode json: %s", err)
	}

	fmt.Print(string(bytes))
	return nil
}

func diffXcodeOut(out output.Data) {
	out.Walk(func(_ string, f output.Field, _ int) {
		fmt.Printf("%s:%d: warning: ", f.File, f.Line)
		fmt.Printf("%s is deprecated ", f.Field)
		fmt.Printf("- Reason: %s", f.DeprecationReason)
		fmt.Println()
	})
}
