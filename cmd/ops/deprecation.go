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
		Use:   "deprecation [flags] queries_directory",
		Short: "Find deprecated fields in queries and mutations",
		Args:  deprecationCmdArgsValidation,
		RunE:  deprecationsCmdRun,
	}
)

func init() {
	Program.AddCommand(deprecationsCmd)
	deprecationsCmd.Flags().StringVar(&schemaFile, schemaFileFlag, "", "Server's schema file (required)")
	deprecationsCmd.MarkFlagRequired(schemaFileFlag) //nolint:errcheck // will err if flag doesn't exist
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

	queriesDir := args[0]
	queryFields, err := parser.ParseQueryDir(queriesDir, schema)
	if err != nil {
		return fmt.Errorf("Unable to parse files in %s: %s", queriesDir, err)
	}

	switch outputFormat {
	case stdoutFormat:
		deprecationStdOut(queryFields)

	case jsonFormat:
		err = deprecationJsonOut(queryFields)
		if err != nil {
			return err
		}

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
