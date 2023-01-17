package ops

import (
	"fmt"

	"github.com/spf13/cobra"
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
	Program.PersistentFlags().StringVar(&outputFormat, "output", "stdout", "Output format. Choose between json and stdout. Defaults is stdout.")

	// TODO: This is required because the test suite doesn't finish the program and flags are not reset. Find a better way to do this.
	cobra.OnFinalize(func() {
		outputFormat = "stdout"
	})
}

func OnFinalize() {
	fmt.Println("Finalizing gql-lint")
}
