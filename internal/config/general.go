package config

import "github.com/willie68/wssg/internal/logging"

const (
	WssgFolder = ".wssg"
	// buldin processors
	ProcInternal = "internal"
)

type General map[string]any

var (
	log = logging.New().WithName("config")
)
