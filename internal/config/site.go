package config

import (
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/willie68/wssg/internal/logging"
	"gopkg.in/yaml.v3"
)

const SiteFile = "siteconfig.yaml"

type Site struct {
	BaseURL        string `yaml:"baseurl"`
	Title          string `yaml:"title"`
	Description    string `yaml:"description"`
	Keywords       string `yaml:"keywords"`
	Language       string `yaml:"language"`
	UserProperties General
}

var SiteDefault = Site{
	BaseURL:     "example.com",
	Title:       "example",
	Description: "a short description of this site",
	Keywords:    "tutorial basic static website",
	Language:    "en",
}

var (
	// SiteConfig this is the main configuration for this site
	SiteConfig Site
	siteLoaded bool
)

func init() {
	_ = mergo.Merge(&SiteConfig, SiteDefault)
	siteLoaded = false
}

// LoadSite loading the site config
func LoadSite(rootFolder string) Site {
	if siteLoaded {
		return SiteConfig
	}
	fd := filepath.Join(rootFolder, WssgFolder, SiteFile)
	if _, err := os.Stat(fd); err != nil {
		if os.IsNotExist(err) {
			return SiteConfig
		}
		log.Errorf("site config: %v", err)
	}
	dt, err := os.ReadFile(fd)
	if err != nil {
		log.Errorf("can't read site config: %v", err)
		panic(1)
	}
	// Load main parts
	err = yaml.Unmarshal(dt, &SiteConfig)
	if err != nil {
		log.Errorf("can't read site config: %v", err)
		panic(1)
	}
	// Load user properties
	err = yaml.Unmarshal(dt, &SiteConfig.UserProperties)
	if err != nil {
		log.Errorf("can't read site config: %v", err)
		panic(1)
	}
	siteLoaded = true
	return SiteConfig
}

func (s *Site) General() (output General) {
	log := logging.New().WithName("siteconfig")
	output = make(General)
	output["baseurl"] = s.BaseURL
	output["title"] = s.Title
	output["description"] = s.Description
	output["keywords"] = s.Keywords
	output["language"] = s.Language
	err := mergo.Merge(&output, s.UserProperties)
	if err != nil {
		log.Errorf("error merging user properties: %v", err)
	}
	return
}
