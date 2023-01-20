package ops

import (
	"fmt"

	"github.com/spf13/cobra"
)

func exactlyNArgsValidator(n int, errMsg string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf(errMsg)
		}
		return nil
	}
}

func atLeastNArgsValidator(n int, errMsg string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) >= n {
			return fmt.Errorf(errMsg)
		}
		return nil
	}
}
