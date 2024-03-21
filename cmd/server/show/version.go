package show

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/project"
	"github.com/spf13/cobra"

	"github.com/kanthorlabs/kanthor/cmd/base"
)

func NewVersion(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "show components version of the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"#", "Name", "Version"})

			names := base.ServiceNames()
			for i := range names {
				t.AppendRows([]table.Row{
					{i + 1, names[i], project.GetVersion()},
				})
			}
			t.Render()

			return nil
		},
	}

	return command
}
