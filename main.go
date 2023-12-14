package main

import (
	"os"

	"github.com/willie68/wssg/cmd/wssg/cmd"
	"github.com/willie68/wssg/internal/logging"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		logging.Root.Errorf("error on command: %v", err)
		os.Exit(1)
	}
}
