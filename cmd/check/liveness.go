package check

import (
	"github.com/kanthorlabs/common/healthcheck/background"
	"github.com/kanthorlabs/common/healthcheck/config"
	"github.com/kanthorlabs/kanthor/services"
	"github.com/spf13/cobra"
)

func NewLiveness() *cobra.Command {
	command := &cobra.Command{
		Use:       "liveness",
		ValidArgs: services.SERVICES,
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			client, err := background.NewClient(config.Default(serviceName, 5000))
			if err != nil {
				return err
			}
			if err := client.Liveness(); err != nil {
				return err
			}
			return nil
		},
	}
	return command
}
