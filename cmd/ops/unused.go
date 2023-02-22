package ops

import (
	"fmt"

	"github.com/sketch-hq/gql-lint/format"
	"github.com/sketch-hq/gql-lint/input"
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
	unusedCmd.Flags().StringSliceVar(&flags.schemaFiles, schemaFileFlagName, []string{}, "Server's schema as file or url (required)")
	unusedCmd.MarkFlagRequired(schemaFileFlagName) //nolint:errcheck // will err if flag doesn't exist
	unusedCmd.Flags().StringSliceVar(&flags.include, includeFlagName, []string{}, "Only include files matching this pattern")
	unusedCmd.Flags().StringSliceVar(&flags.ignore, ignoreFlagName, []string{}, "Files to ignore")
}

func unusedCmdRun(cmd *cobra.Command, args []string) error {
	queryFiles, err := input.ExpandGlobs(args, flags.include, flags.ignore)
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

	r, err := format.UnusedFormatter.Format(flags.outputFormat, out)
	if err == nil {
		fmt.Print(r)
	}
	return err
}
