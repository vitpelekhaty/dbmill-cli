package main

import (
	"log"

	"github.com/vitpelekhaty/dbmill-cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}
