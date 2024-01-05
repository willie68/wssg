package plain

import (
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/plugins"
)

// Plain internal plain plugin do nothing
type Plain struct {
}

// New create a new internal plugin
func New() plugins.Plugin {
	return &Plain{}
}

// CreateBody interface method to create a html body from a markdown file
func (m *Plain) CreateBody(content []byte, _ model.Page) ([]byte, error) {
	return content, nil
}
