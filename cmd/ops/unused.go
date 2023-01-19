package ops

import (
	"encoding/json"
	"fmt"

	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/unused"
	"github.com/spf13/cobra"
)

var (
	//Command
	unusedCmd = &cobra.Command{
		Use:   "unused [flags] queries_directories",
		Short: "Find unused deprecated fields",
		Args:  unusedCmdArgsValidation,
		RunE:  unusedCmdRun,
	}
)

func init() {
	Program.AddCommand(unusedCmd)
	unusedCmd.Flags().StringVar(&schemaFile, schemaFileFlag, "", "Server's schema file (required)")
	unusedCmd.MarkFlagRequired(schemaFileFlag) //nolint:errcheck // will err if flag doesn't exist
}

func unusedCmdArgsValidation(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must specify at least one directory for queries and mutations")
	}
	return nil
}

func unusedCmdRun(cmd *cobra.Command, queriesDirs []string) error {
	unusedFields, err := unused.GetUnusedFields(schemaFile, queriesDirs)

	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	switch outputFormat {
	case stdoutFormat:
		unusedStdOut(unusedFields)

	case jsonFormat:
		err = unusedJSONOut(unusedFields)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("%s is not a valid output format. Choose between json and stdout", outputFormat)
	}

	return nil
}

func unusedStdOut(fields []unused.UnusedField) {
	if len(fields) == 0 {
		fmt.Println("Nothing can be removed right now")
		return
	}

	for _, q := range fields {
		fmt.Printf("`%s` is unused and can be removed\n", q.Name)
		fmt.Println()
	}
}

func unusedJSONOut(fields []unused.UnusedField) error {
	out := make([]output.UnusedField, len(fields))

	for i, f := range fields {
		out[i] = output.UnusedField{Field: f.Name}
	}
	bytes, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("failed to encode json: %s\n", err)
	}

	fmt.Print(string(bytes))
	return nil
}
