package config

import (
	"dario.cat/mergo"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/processors"
)

const (
	// SectionFileName name of the section file
	SectionFileName = "section.yaml"
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
