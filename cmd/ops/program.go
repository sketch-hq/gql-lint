package ops

import (
	"github.com/spf13/cobra"
)

const (
	outputFormatFlag = "output"
	jsonFormat       = "json"
	stdoutFormat     = "stdout"
)

var (
	//Global Flags
	outputFormat string

	//Binary program
	Program = &cobra.Command{
		Use:   "gql-lint",
		Short: "gql-lint is a tool to lint GraphQL queries and mutations",
	}
)

func init() {
	Program.CompletionOptions.DisableDefaultCmd = true
	Program.PersistentFlags().StringVar(&outputFormat, outputFormatFlag, stdoutFormat, "Output format. Choose between json and stdout. Defaults is stdout.")

	// TODO: This is required because the test suite doesn't finish the program and flags are not reset. Find a better way to do this.
	cobra.OnFinalize(func() {
		outputFormat = stdoutFormat
	})
}
