package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// GenerateFile the file name of the generator config
const GenerateFile = "generate.yaml"

// Generate the configuration for the generator
type Generate struct {
	Output string `yaml:"output"`
}

var (
	// GenDefault the default generator config
	GenDefault = Generate{
		Output: fmt.Sprintf("./%s/output", WssgFolder),
	}
	// GenConfig the actual generator config
	GenConfig Generate
	genLoded  bool
)

func init() {
	GenConfig = GenDefault
}

// LoadGenConfig loading the generator config from the site
func LoadGenConfig(rootFolder string) Generate {
	if genLoded {
		return GenConfig
	}
	fd := filepath.Join(rootFolder, WssgFolder, GenerateFile)
	if _, err := os.Stat(fd); err != nil {
		if os.IsNotExist(err) {
			return GenConfig
		}
		log.Errorf("site config: %v", err)
	}
	dt, err := os.ReadFile(fd)
	if err != nil {
		log.Errorf("can't read generator config: %v", err)
		panic(1)
	}
	err = yaml.Unmarshal(dt, &GenConfig)
	if err != nil {
		log.Errorf("can't read generator config: %v", err)
		panic(1)
	}
	siteLoaded = true
	return GenConfig
}
