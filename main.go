package main

import (
	"os"
	// blank imports for every processor, so that the processor is registered to the di framework
	_ "github.com/willie68/wssg/processors/blog"
	_ "github.com/willie68/wssg/processors/gallery"
	_ "github.com/willie68/wssg/processors/html"
	_ "github.com/willie68/wssg/processors/mdtohtml"
	_ "github.com/willie68/wssg/processors/plain"

	"github.com/willie68/wssg/cmd/wssg/cmd"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ver := config.NewVersion().WithCommit(commit).WithDate(date).WithVersion(version)
	cmd.CmdVersion = *ver
	err := cmd.Execute()
	if err != nil {
		logging.Root.Errorf("error on command: %v", err)
		os.Exit(1)
	}
}
