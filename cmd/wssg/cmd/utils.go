package cmd

import (
	"os"

	"gopkg.in/yaml.v3"
)

func fileExists(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func writeAsYaml(file string, v any) error {
	dt, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, dt, 0755)
	return err
}
