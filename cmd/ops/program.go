package ops

import (
	"github.com/spf13/cobra"
)

const (
	outputFormatFlag = "output"
	jsonFormat       = "json"
	stdoutFormat     = "stdout"
	xcodeFormat      = "xcode"
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
	Program.PersistentFlags().StringVar(&outputFormat, outputFormatFlag, stdoutFormat, "Output format. Choose between stdout, json, xcode.")

	// This is required because the test suite doesn't finish the program and flags are not reset
	cobra.OnFinalize(func() {
		outputFormat = stdoutFormat
	})
}
