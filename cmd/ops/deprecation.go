package ops

import (
	"encoding/json"
	"fmt"
	"strings"

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
		return fmt.Errorf("Error: %s", err)
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
			return fmt.Errorf("Unable to parse files: %s", err)
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

	switch flags.outputFormat {
	case stdoutFormat:
		deprecationStdOut(out)

	case jsonFormat:
		err := deprecationJsonOut(out)
		if err != nil {
			return err
		}
	case xcodeFormat:
		deprecationXcodeOut(out)
	case annotateFormat:
		deprecationAnnotateOut(out)
	default:
		return fmt.Errorf("%s is not a valid output format. Choose between json and stdout", flags.outputFormat)
	}

	return nil
}

func deprecationStdOut(out output.Data) {
	out.Walk(func(schema string, f output.Field, i int) {
		if i == 0 {
			fmt.Println("Schema:", schema)
		}
		fmt.Printf("  %s is deprecated\n", f.Field)
		fmt.Printf("    File:   %s:%d\n", f.File, f.Line)
		reason := strings.ReplaceAll(f.DeprecationReason, "\n", " ")
		fmt.Printf("    Reason: %s\n", reason)
		fmt.Println()
	})
}

func deprecationJsonOut(out output.Data) error {
	bytes, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("Failed to encode json: %s\n", err)
	}

	fmt.Print(string(bytes))
	return nil
}

func deprecationXcodeOut(out output.Data) {
	out.Walk(func(_ string, f output.Field, _ int) {
		fmt.Printf("%s:%d: warning: ", f.File, f.Line)
		fmt.Printf("%s is deprecated ", f.Field)
		reason := strings.ReplaceAll(f.DeprecationReason, "\n", " ")
		fmt.Printf("- Reason: %s", reason)
		fmt.Println()
	})
}

func deprecationAnnotateOut(out output.Data) {
	var replacer = strings.NewReplacer(
		"\n", "\\n",
		"\"", "\\\"",
	)

	fmt.Print("[")
	out.Walk(func(_ string, f output.Field, idx int) {
		if idx > 0 {
			fmt.Print(",\n")
		}
		fmt.Printf("{ \"file\": \"%s\", \"line\": %d, ", f.File, f.Line)
		fmt.Printf("\"title\": \"%s is deprecated\", ", f.Field)
		reason := replacer.Replace(f.DeprecationReason)
		fmt.Printf("\"message\": \"%s\", ", reason)
		fmt.Printf("\"annotation_level\": \"warning\" }")
	})
	fmt.Print("]")
}
