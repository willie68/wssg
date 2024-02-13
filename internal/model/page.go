package model

import "github.com/stretchr/objx"

// Page this is the internal data model for a page
type Page struct {
	Name         string `json:"name" yaml:"name"`
	Title        string `json:"title" yaml:"title"`
	Path         string `json:"path" yaml:"path"`
	URLPath      string `json:"urlpath" yaml:"urlpath"`
	Order        int    `json:"order" yaml:"order"`
	Processor    string `json:"processor" yaml:"processor"`
	Filename     string
	Section      string
	SourceFolder string
	DestFolder   string
	Cnf          objx.Map
}
