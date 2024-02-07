package cmd

import (
	"github.com/kanthorlabs/kanthor/cmd/check"
	"github.com/kanthorlabs/kanthor/cmd/config"
	"github.com/kanthorlabs/kanthor/cmd/migrate"
	"github.com/kanthorlabs/kanthor/cmd/serve"
	"github.com/kanthorlabs/kanthor/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{}

	command.AddCommand(check.New())
	command.AddCommand(config.New(provider))
	command.AddCommand(migrate.New(provider))
	command.AddCommand(serve.New(provider))

	command.PersistentFlags().BoolP("verbose", "", false, "--verbose | show more information")
	return command
}
