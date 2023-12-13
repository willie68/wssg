package utils

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func FileExists(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func WriteAsYaml(file string, v any) error {
	dt, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, dt, 0755)
	return err
}

func LoadYAML(file string, v any) error {
	if ok, _ := FileExists(file); !ok {
		return fmt.Errorf("file not found: %s", file)
	}
	dt, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(dt, v)
	return err
}

func FileCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
