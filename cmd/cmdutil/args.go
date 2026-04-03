package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("%d個の引数が必要ですが、%d個指定されました\n\nUsage:\n  %s", n, len(args), cmd.UseLine())
		}
		return nil
	}
}
