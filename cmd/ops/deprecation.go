package ops

import (
	"encoding/json"
	"fmt"

	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/parser"
	"github.com/spf13/cobra"
)

const (
	schemaFileFlag = "schema"
)

var (
	//Flags
	schemaFile string

	//Command
	deprecationsCmd = &cobra.Command{
		Use:   "deprecation [flags] queries_directory|queries_files_list",
		Short: "Find deprecated fields in queries and mutations given a directory or a list of files",
		Long: `
Find deprecated fields in queries and mutations given a directory or a list of files.

The "queries_directory" argument is a directory containing all the queries and mutations. They can be in subdirectories. 
The "queries_files_list" argument is a file containing a list of paths to queries and mutations. The file should contain one query or mutation per line.`,
		Args: deprecationCmdArgsValidation,
		RunE: deprecationsCmdRun,
	}
)

func init() {
	Program.AddCommand(deprecationsCmd)
	deprecationsCmd.Flags().StringVar(&schemaFile, schemaFileFlag, "", "Server's schema file (required)")
	deprecationsCmd.MarkFlagRequired(schemaFileFlag) //nolint:errcheck // will err if flag doesn't exist

	// This is required because the test suite doesn't finish the program and flags are not reset.
	cobra.OnFinalize(func() {
		schemaFile = ""
	})
}

func deprecationCmdArgsValidation(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("You must specify a directory for queries and mutations")
	}
	return nil
}

func deprecationsCmdRun(cmd *cobra.Command, args []string) error {
	schema, err := parser.ParseSchemaFile(schemaFile)
	if err != nil {
		return fmt.Errorf("Unable to parse schema file %s: %s", schemaFile, err)
	}

	queriesSource := args[0]
	queryFields, err := parser.ParseQuerySource(queriesSource, schema)
	if err != nil {
		return fmt.Errorf("Unable to parse files in %s: %s", queriesSource, err)
	}

	switch outputFormat {
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
		return fmt.Errorf("%s is not a valid output format. Choose between json and stdout", outputFormat)
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
