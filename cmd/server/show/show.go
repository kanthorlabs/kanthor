package show

import (
	"github.com/kanthorlabs/common/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:   "show",
		Short: "show information about the server",
	}

	command.AddCommand(NewVersion(provider))
	return command
}
