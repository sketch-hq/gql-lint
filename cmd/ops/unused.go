package ops

import (
	"encoding/json"
	"fmt"

	"github.com/sketch-hq/gql-lint/input"
	"github.com/sketch-hq/gql-lint/output"
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
	unusedCmd.Flags().StringArrayVar(&flags.schemaFiles, schemaFileFlagName, []string{}, "Server's schema as file or url. Can be repeated (required)")
	unusedCmd.MarkFlagRequired(schemaFileFlagName) //nolint:errcheck // will err if flag doesn't exist

	unusedCmd.Flags().StringArrayVar(&flags.ignore, ignoreFlagName, []string{}, "Files to ignore")
}

func unusedCmdRun(cmd *cobra.Command, args []string) error {
	queryFiles, err := input.ExpandGlobs(args, flags.ignore)
	if err != nil {
		return err
	}

	if flags.verbose {
		fmt.Println("debug: Processing the following query files:")
		for _, file := range queryFiles {
			fmt.Println("  -", file)
		}
	}

	out, err := unused.GetUnusedFields(flags.schemaFiles, queryFiles, flags.verbose)
	if err != nil {
		return err
	}

	switch flags.outputFormat {
	case stdoutFormat:
		unusedStdOut(out)

	case markdownFormat:
		unusedMarkdownOut(out)

	case jsonFormat:
		err = unusedJSONOut(out)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("%s is not a valid output format. Choose between json, markdown and stdout", flags.outputFormat)
	}

	return nil
}

func unusedStdOut(out output.Data) {
	out.Walk(func(schema string, f output.Field, i int) {
		if i == 0 {
			fmt.Println("Schema:", schema)
		}
		fmt.Printf("  %s (line %d) is unused and can be removed \n", f.Field, f.Line)
		fmt.Println()
	})
}

func unusedJSONOut(out output.Data) error {
	bytes, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("failed to encode json: %s", err)
	}

	fmt.Print(string(bytes))
	return nil
}

func unusedMarkdownOut(out output.Data) {
	hasUnused := false

	out.Walk(func(schema string, f output.Field, fieldIdx int) {
		if fieldIdx == 0 {
			hasUnused = true
			fmt.Println("**", schema, "**")
		}
		fmt.Printf("- %s (line `%d`)\n", f.Field, f.Line)
	})

	if !hasUnused {
		fmt.Println("Nothing can be removed right now")
	}
}
