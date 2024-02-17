package config

import (
	"fmt"

	"dario.cat/mergo"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/processors"
)

const (
	// SectionFileName name of the section file
	SectionFileName = "section.yaml"
)

var (
	keyNames = []string{"name", "title", "processor"}
)

// Section the configuration of a section
type Section struct {
	Name           string `yaml:"name"`
	Title          string `yaml:"title"`
	Processor      string `yaml:"processor"`
	Order          int    `yaml:"order"`
	URLPath        string
	UserProperties objx.Map
}

// SectionDefault the default configuration of a section, used for creating a new one
var SectionDefault = Section{
	Name:      "{{.name}}",
	Title:     "{{.name}}",
	Processor: processors.DefaultProcessor,
}

// MSA convert this section to general
func (s Section) MSA() objx.Map {
	output := objx.Map{
		"name":      s.Name,
		"title":     s.Title,
		"processor": s.Processor,
	}
	err := mergo.Merge(&output, s.UserProperties)
	if err != nil {
		log.Errorf("error merging user properties: %v", err)
	}
	return output
}

// G2Section convert a general struct to a section
func G2Section(g objx.Map) Section {
	up := make(objx.Map)
	for k, v := range g {
		use := true
		for _, f := range keyNames {
			if k == f {
				use = false
			}
		}
		if use {
			up[k] = v
		}
	}
	name := g.Get("name").Str("no_name")
	return Section{
		Name:           name,
		Title:          g.Get("title").Str("no title given"),
		Processor:      g.Get("processor").Str(processors.DefaultProcessor),
		URLPath:        fmt.Sprintf("/%s", name),
		Order:          g.Get("order").Int(0),
		UserProperties: up,
	}
}
