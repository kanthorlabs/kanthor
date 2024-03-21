package check

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "check",
		Short: "check the health of backaground components",
	}

	command.AddCommand(NewReadiness())
	command.AddCommand(NewLiveness())
	return command
}
