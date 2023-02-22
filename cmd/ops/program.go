package ops

import (
	"github.com/sketch-hq/gql-lint/format"
	"github.com/spf13/cobra"
)

var Program = &cobra.Command{
	Use:   "gql-lint",
	Short: "gql-lint is a tool to lint GraphQL queries and mutations",
}

func init() {
	Program.CompletionOptions.DisableDefaultCmd = true
	Program.PersistentFlags().StringVar(
		&flags.outputFormat,
		outputFormatFlagName,
		format.StdoutFormat,
		"Output format. Choose between stdout, json, xcode.",
	)
	Program.PersistentFlags().BoolVarP(
		&flags.verbose,
		verboseFlagName,
		"v",
		false,
		"Verbose mode. Will print debug messages",
	)

	// This is required because the test suite doesn't finish the program and flags are not reset
	cobra.OnFinalize(func() {
		setFlagsToDefault()
	})
}
