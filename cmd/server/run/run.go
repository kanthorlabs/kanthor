package run

import (
	"slices"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/kanthor/cmd/base"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:       "run",
		Short:     "run a single service or multiple services",
		ValidArgs: append(base.ServiceNames(), ALL),
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			if slices.Contains(args, ALL) || len(args) > 1 {
				return multiple(provider, args)
			}

			return single(provider, args[0])
		},
	}

	return command
}
