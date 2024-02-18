package html

import (
	"fmt"

	"github.com/samber/do"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/processors/processor"
)

// Processor internal plain processor do nothing
type Processor struct {
}

func init() {
	proc := New()
	do.ProvideNamedValue[processor.Processor](nil, proc.Name(), proc)
}

// New create a new plain processor
func New() processor.Processor {
	return &Processor{}
}

// Name returning the name of this processor
func (p *Processor) Name() string {
	return "html"
}

// AddPage adding the new page
func (p *Processor) AddPage(folder, pagefile string) (m objx.Map, err error) {
	return
}

// GetPageTemplate getting the right template for the named page
func (p *Processor) GetPageTemplate(name string) string {
	return fmt.Sprintf("this is a file with the name %s", name)
}

// CreateBody interface method to create a html body from a markdown file
func (p *Processor) CreateBody(content []byte, _ model.Page) (*processor.Response, error) {
	return &processor.Response{
		Render: true,
		Body:   string(content),
	}, nil
}

// HTMLTemplateName returning the used html template
func (p *Processor) HTMLTemplateName() string {
	return "layout.html"
}
