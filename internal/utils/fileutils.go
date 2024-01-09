package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileExists checking if a file exists
func FileExists(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// WriteAsYaml writing a any struct as yaml
func WriteAsYaml(file string, v any) error {
	dt, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, dt, 0755)
	return err
}

// LoadYAML loading a file, yaml unmarshal into a any struct
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

// FileCopy convinient method for copy a file
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

func FileNameWOExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
