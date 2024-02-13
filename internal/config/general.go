package config

import (
	"github.com/willie68/wssg/internal/logging"
)

const (
	// WssgFolder folder name for the configuration
	WssgFolder = ".wssg"
)

var (
	log = logging.New().WithName("config")
)
