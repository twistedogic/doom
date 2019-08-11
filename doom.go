package main

import (
	"log"

	"github.com/twistedogic/doom/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
}
