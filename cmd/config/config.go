package config

import (
	"github.com/kanthorlabs/common/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use: "config",
	}
	command.AddCommand(NewValidate(provider))
	return command
}
