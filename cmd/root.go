package cmd

import (
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/kanthor/cmd/check"
	"github.com/kanthorlabs/kanthor/cmd/config"
	"github.com/kanthorlabs/kanthor/cmd/serve"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{}

	command.AddCommand(check.New())
	command.AddCommand(config.New(provider))
	command.AddCommand(serve.New(provider))

	command.PersistentFlags().BoolP("verbose", "", false, "--verbose | show more information")
	return command
}
