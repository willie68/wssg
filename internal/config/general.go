package config

import "github.com/willie68/wssg/internal/logging"

const (
	// WssgFolder folder name for the configuration
	WssgFolder = ".wssg"
)

// General this map is widely used for the templating
type General map[string]any

var (
	log = logging.New().WithName("config")
)
