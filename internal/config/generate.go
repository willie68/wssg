package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const GenerateFile = "generate.yaml"

type Generate struct {
	Output string `yaml:"output"`
}

var (
	GenDefault = Generate{
		Output: fmt.Sprintf("./%s/output", WssgFolder),
	}
	GenConfig Generate
	genLoded  bool
)

func init() {
	GenConfig = GenDefault
}

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
