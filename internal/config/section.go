package config

import (
	"fmt"

	"dario.cat/mergo"
	"github.com/stretchr/objx"
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
	UserProperties General
}

// SectionDefault the default configuration of a section, used for creating a new one
var SectionDefault = Section{
	Name:      "{{.name}}",
	Title:     "{{.name}}",
	Processor: ProcMarkdown,
}

// General convert this section to general
func (s Section) General() (output General) {
	output = make(General)
	output["name"] = s.Name
	output["title"] = s.Title
	output["processor"] = s.Processor
	err := mergo.Merge(&output, s.UserProperties)
	if err != nil {
		log.Errorf("error merging user properties: %v", err)
	}
	return
}

// G2Section convert a general struct to a section
func G2Section(g General) Section {
	m := objx.New(map[string]any(g))
	up := make(General)
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
	name := m.Get("name").Str("no_name")
	return Section{
		Name:           name,
		Title:          m.Get("title").Str("no title given"),
		Processor:      m.Get("processor").Str(ProcMarkdown),
		URLPath:        fmt.Sprintf("/%s", name),
		Order:          m.Get("order").Int(0),
		UserProperties: up,
	}
}
