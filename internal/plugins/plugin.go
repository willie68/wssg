package plugins

import "github.com/willie68/wssg/internal/model"

// Response contains scripts, styles and the body
type Response struct {
	Body   string
	Script string
	Style  string
}

//Plugin this is the interface for a plugin
type Plugin interface {
	CreateBody(content []byte, pg model.Page) (*Response, error)
	HTMLTemplateName() string
}
