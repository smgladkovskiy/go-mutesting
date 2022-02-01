package main

import (
	"os"

	"github.com/smgladkovskiy/go-mutesting/cmd"
	log "github.com/spacetab-io/logs-go/v2"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}

	os.Exit(0)
}
