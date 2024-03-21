package main

import (
	"log"

	_ "embed"

	"github.com/kanthorlabs/kanthor/cmd/base"
	"github.com/kanthorlabs/kanthor/cmd/data/migrate"
)

func main() {
	_, command := base.New()
	command.AddCommand(migrate.New())

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
