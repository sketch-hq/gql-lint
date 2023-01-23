package ops

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ExactArgs(n int, errMsg string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(n)(cmd, args); err != nil {
			return fmt.Errorf(errMsg)
		}
		return nil
	}
}

func MinimumNArgs(n int, errMsg string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(n)(cmd, args); err != nil {
			return fmt.Errorf(errMsg)
		}
		return nil
	}
}
