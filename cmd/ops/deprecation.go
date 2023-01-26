package ops

import (
	"encoding/json"
	"fmt"

	"github.com/sketch-hq/gql-lint/input"
	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/parser"
	"github.com/sketch-hq/gql-lint/schema"
	"github.com/spf13/cobra"
)

var deprecationsCmd = &cobra.Command{
	Use:   "deprecation [flags] queries",
	Short: "Find deprecated fields in queries and mutations given a list of files",
	Long: `
Find deprecated fields in queries and mutations given a list of files

The "queries" argument is a file glob matching one or more graphql query or mutation files.`,
	Args: MinimumNArgs(1, "you must specify at least one file with queries or mutations"),
	RunE: deprecationsCmdRun,
}

func init() {
	Program.AddCommand(deprecationsCmd)
	deprecationsCmd.Flags().StringVar(&flags.schemaFile, schemaFileFlagName, "", "Server's schema as file or url (required)")
	deprecationsCmd.MarkFlagRequired(schemaFileFlagName) //nolint:errcheck // will err if flag doesn't exist

	deprecationsCmd.Flags().StringArrayVar(&flags.ignore, ignoreFlagName, []string{}, "Files to ignore (as file blob)")
}

func deprecationsCmdRun(cmd *cobra.Command, args []string) error {
	schema, err := schema.Load(flags.schemaFile)
	if err != nil {
		return err
	}

	queryFiles, err := input.ExpandGlobs(args, flags.ignore)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	queryFields, err := parser.ParseQuerySource(queryFiles, schema)
	if err != nil {
		return fmt.Errorf("Unable to parse files: %s", err)
	}

	switch flags.outputFormat {
	case stdoutFormat:
		deprecationStdOut(queryFields)

	case jsonFormat:
		err = deprecationJsonOut(queryFields)
		if err != nil {
			return err
		}
	case xcodeFormat:
		deprecationXcodeOut(queryFields)
	default:
		return fmt.Errorf("%s is not a valid output format. Choose between json and stdout", flags.outputFormat)
	}

	return nil
}

func deprecationStdOut(queryFields parser.QueryFieldList) {
	for _, q := range queryFields {
		fmt.Printf("%s is deprecated\n", q.Path)
		fmt.Printf("  File:   %s:%d\n", q.File, q.Line)
		fmt.Printf("  Reason: %s\n", q.DeprecationReason)
		fmt.Println()
	}
}

func deprecationJsonOut(queryFields parser.QueryFieldList) error {
	out := output.Data{}

	for _, q := range queryFields {
		f := output.Field{
			Field:             q.Path,
			File:              q.File,
			Line:              q.Line,
			DeprecationReason: q.DeprecationReason,
		}
		out = append(out, f)
	}
	bytes, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("Failed to encode json: %s\n", err)
	}

	fmt.Print(string(bytes))
	return nil
}

func deprecationXcodeOut(queryFields parser.QueryFieldList) {
	for _, q := range queryFields {
		fmt.Printf("%s:%d: warning: ", q.File, q.Line)
		fmt.Printf("%s is deprecated ", q.Path)
		fmt.Printf("- Reason: %s\n", q.DeprecationReason)
		fmt.Println()
	}
}
