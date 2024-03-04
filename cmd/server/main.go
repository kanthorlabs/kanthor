package main

import (
	"log"

	_ "embed"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/kanthor/cmd/server/run"
	"github.com/spf13/cobra"
)

//go:embed .version
var version string

func main() {
	project.SetVersion(version)

	provider, err := configuration.New(project.Namespace())
	if err != nil {
		panic(err)
	}
	command := New(provider)

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Println("--- error ---")
				log.Println(err.Error())
			}
		}
	}()

	if err := command.Execute(); err != nil {
		panic(err)
	}
}

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}

	// sub-command
	command.AddCommand(run.New(provider))

	command.PersistentFlags().BoolP("verbose", "v", false, "show verbose output including debug information")
	return command
}

var short = "Open-source Webhook Gateway: Delivery your message with precision and ease"
var long = " _  __           _   _                " + "\n" +
	"| |/ /__ _ _ __ | |_| |__   ___  _ __ " + "\n" +
	"| ' // _` | '_ \\| __| '_ \\ / _ \\| '__|" + "\n" +
	"| . \\ (_| | | | | |_| | | | (_) | |   " + "\n" +
	"|_|\\_\\__,_|_| |_|\\__|_| |_|\\___/|_|   " + "\n" +
	"\n" +
	short
