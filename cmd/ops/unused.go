package ops

import (
	"encoding/json"
	"fmt"

	"github.com/sketch-hq/gql-lint/input"
	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/schema"
	"github.com/sketch-hq/gql-lint/unused"
	"github.com/spf13/cobra"
)

var (
	//Command
	unusedCmd = &cobra.Command{
		Use:   "unused [flags] queries",
		Short: "Find unused deprecated fields",
		Long: `
Find unused deprecated fields

The "queries" argument is a file glob matching one or more graphql query or mutation files.`,
		Args: MinimumNArgs(1, "you must specify at least one file with queries or mutations"),
		RunE: unusedCmdRun,
	}
)

func init() {
	Program.AddCommand(unusedCmd)
	unusedCmd.Flags().StringVar(&flags.schemaFile, schemaFileFlagName, "", "Server's schema as file or url (required)")
	unusedCmd.MarkFlagRequired(schemaFileFlagName) //nolint:errcheck // will err if flag doesn't exist
}

func unusedCmdRun(cmd *cobra.Command, args []string) error {
	schema, err := schema.Load(flags.schemaFile)
	if err != nil {
		return err
	}

	queryFiles, err := input.QueryFiles(args)
	if err != nil {
		return err
	}

	unusedFields, err := unused.GetUnusedFields(schema, queryFiles)
	if err != nil {
		return err
	}

	switch flags.outputFormat {
	case stdoutFormat:
		unusedStdOut(unusedFields)

	case jsonFormat:
		err = unusedJSONOut(unusedFields)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("%s is not a valid output format. Choose between json and stdout", flags.outputFormat)
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
