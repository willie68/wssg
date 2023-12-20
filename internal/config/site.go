package config

import (
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

const SiteFile = "siteconfig.yaml"

type Site struct {
	BaseURL     string `yaml:"baseurl"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Keywords    string `yaml:"keywords"`
	Language    string `yaml:"language"`
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
	err = yaml.Unmarshal(dt, &SiteConfig)
	if err != nil {
		log.Errorf("can't read site config: %v", err)
		panic(1)
	}
	siteLoaded = true
	return SiteConfig
}
