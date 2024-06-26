package main

import (
	"log"

	_ "embed"

	"github.com/kanthorlabs/kanthor/cmd/base"
	"github.com/kanthorlabs/kanthor/cmd/server/run"
	"github.com/kanthorlabs/kanthor/cmd/server/show"
)

func main() {
	provider, command := base.New()
	command.AddCommand(run.New(provider))
	command.AddCommand(show.New(provider))

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
