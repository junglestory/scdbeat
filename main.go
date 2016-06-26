package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/junglestory/scdbeat/beater"
)

func main() {
	err := beat.Run("scdbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
