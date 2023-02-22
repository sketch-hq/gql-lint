package ops

import (
	"fmt"

	"github.com/sketch-hq/gql-lint/format"
	"github.com/sketch-hq/gql-lint/input"
	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/sources"
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
	deprecationsCmd.Flags().StringSliceVar(&flags.schemaFiles, schemaFileFlagName, []string{}, "Server's schema as file or url (required)")
	deprecationsCmd.MarkFlagRequired(schemaFileFlagName) //nolint:errcheck // will err if flag doesn't exist
	deprecationsCmd.Flags().StringSliceVar(&flags.include, includeFlagName, []string{}, "Only include files matching this pattern")
	deprecationsCmd.Flags().StringSliceVar(&flags.ignore, ignoreFlagName, []string{}, "Files to ignore")
}

func deprecationsCmdRun(cmd *cobra.Command, args []string) error {
	queryFiles, err := input.ExpandGlobs(args, flags.include, flags.ignore)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	if flags.verbose {
		fmt.Println("debug: Processing the following query files:")
		for _, file := range queryFiles {
			fmt.Println("  -", file)
		}
	}

	out := output.Data{}
	for _, schemaFile := range flags.schemaFiles {
		schema, err := sources.LoadSchema(schemaFile)
		if err != nil {
			return err
		}

		if flags.verbose {
			fmt.Println("debug: Succesfully loaded schema from", schemaFile)
		}

		queryFields, err := sources.LoadQueries(schema, queryFiles)
		if err != nil {
			return fmt.Errorf("unable to parse files: %s", err)
		}

		for _, q := range queryFields {
			f := output.Field{
				Field:             q.Path,
				File:              q.File,
				Line:              q.Line,
				DeprecationReason: q.DeprecationReason,
			}
			out.AppendField(schemaFile, f)
		}
	}

	r, err := format.DeprecationFormatter.Format(flags.outputFormat, out)
	if err == nil {
		fmt.Print(r)
	}
	return err
}
