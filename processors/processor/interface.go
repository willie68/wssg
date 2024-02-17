package processor

import (
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/model"
)

// Response contains scripts, styles and the body
type Response struct {
	Body   string
	Script string
	Style  string
}

// Processor this is the interface for a processor
type Processor interface {
	// Name Getting the name of the processor
	Name() string

	// AddPage this method will be called if a new page will be added to a secrion of this processor
	AddPage(folder, pagefile string) (objx.Map, error)

	// GetPageTemplate getting the page template for a page with that name
	GetPageTemplate(name string) string

	// CreateBody creates a body that should be saved to the output folder.
	// if the body part is empty, nothing is to do for this source file.
	CreateBody(content []byte, pg model.Page) (*Response, error)

	// HTMLTemplateName returnin the html template used for this processor
	HTMLTemplateName() string
}
