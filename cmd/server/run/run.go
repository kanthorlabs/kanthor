package run

import (
	"errors"
	"slices"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/opentelemetry"
	"github.com/kanthorlabs/kanthor/cmd/base"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:       "run",
		Short:     "run a single service or multiple services",
		ValidArgs: append(base.ServiceNames(), ALL),
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = opentelemetry.Setup(cmd.Context()); err != nil {
				return
			}
			defer func() {
				if teardownerr := opentelemetry.Teardown(cmd.Context()); teardownerr != nil {
					err = errors.Join(err, teardownerr)
				}
			}()

			if slices.Contains(args, ALL) || len(args) > 1 {
				err = multiple(provider, args)
				return
			}

			err = single(provider, args[0])
			return
		},
	}

	return command
}
