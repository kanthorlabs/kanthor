package main

import (
	_ "embed"
	"log"

	"github.com/kanthorlabs/kanthor/cmd"
	"github.com/kanthorlabs/kanthor/configuration"
	"github.com/kanthorlabs/kanthor/project"
)

//go:embed .version
var version string

func main() {
	project.SetVersion(version)

	provider, err := configuration.New()
	if err != nil {
		panic(err)
	}
	command := cmd.New(provider)

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
