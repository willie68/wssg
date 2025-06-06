package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// GenerateFile the file name of the generator config
const GenerateFile = "generate.yaml"

// Generate the configuration for the generator
type Generate struct {
	// Output where to output the generated site files
	Output string `yaml:"output"`
	// Processors conect mime types with processors
	ProcMime map[string]string `yaml:"procmime"`
	// Autoreload script
	Autoreload string `yaml:"autoreload"`
	// force forces to create everything newly
	Force bool `yaml:"force"`
	// this is the part before all of the website starts, used when a page is deployed with a subpath
	Basepath string `yaml:"basepath"`
}

var (
	// GenDefault the default generator config
	GenDefault = Generate{
		Output: fmt.Sprintf("./%s/output", WssgFolder),
		ProcMime: map[string]string{
			"text/html":     "plain",
			"text/markdown": "markdown",
			"text/plain":    "plain",
		},
		Autoreload: "",
		Basepath:   "/",
	}
	// cfg the actual generator config
	cfg      Generate
	genLoded bool
)

func init() {
	cfg = GenDefault
}

// LoadGenConfig loading the generator config from the site
func LoadGenConfig(rootFolder string) Generate {
	if genLoded {
		return cfg
	}
	fd := filepath.Join(rootFolder, WssgFolder, GenerateFile)
	if _, err := os.Stat(fd); err != nil {
		if os.IsNotExist(err) {
			return cfg
		}
		log.Errorf("site config: %v", err)
	}
	dt, err := os.ReadFile(fd)
	if err != nil {
		log.Errorf("can't read generator config: %v", err)
		os.Exit(-1)
	}
	err = yaml.Unmarshal(dt, &cfg)
	if err != nil {
		log.Errorf("can't read generator config: %v", err)
		os.Exit(-1)
	}
	siteLoaded = true
	if cfg.Basepath == "" {
		cfg.Basepath = "/"
	}
	if len(cfg.Basepath) > 1 && !strings.HasSuffix(cfg.Basepath, "/") {
		cfg.Basepath = cfg.Basepath + "/"
	}
	return cfg
}
