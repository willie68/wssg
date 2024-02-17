package config

import (
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/processors"
)

// Page the configuration of a single page, used by frontmatter
type Page struct {
	Title     string `yaml:"title"`
	Name      string `yaml:"name"`
	Processor string `yaml:"processor"`
}

var (
	// PageDefault the page default config
	PageDefault = Page{
		Title:     "{{.name}}",
		Name:      "{{.name}}",
		Processor: processors.DefaultProcessor,
	}
)

// MSA converting the page sturct into a general
func (p Page) MSA() (output objx.Map) {
	return objx.Map{
		"title":     p.Title,
		"name":      p.Name,
		"processor": p.Processor,
	}
}
