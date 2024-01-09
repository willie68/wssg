package plugins

import "github.com/willie68/wssg/internal/model"

//Plugin this is the interface for a plugin
type Plugin interface {
	CreateBody(content []byte, pg model.Page) ([]byte, error)
	HTMLTemplateName() string
}
