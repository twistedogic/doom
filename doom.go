package main

import (
	"log"

	"github.com/twistedogic/doom/cmd"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
}
