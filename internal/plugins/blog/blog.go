package blog

import (
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/utils"
	"github.com/willie68/wssg/templates"
	"gopkg.in/yaml.v3"
)

// this is a generator generating a blog with pagination
// every blogentry is a single markdown file. The index.md is the starting page for this.

const (
	PluginName = "blog"
)

type BlogEntry struct {
	Name    string    `yaml:"name"`
	Created time.Time `yaml:"created"`
}

func GetPageTemplate(name string) string {
	if name == "index" {
		return templates.BlogIndex
	}
	return templates.BlogPage
}

// AddBlogPage adding the new blog page to the list of pages with the actual time.
// if already there the entry will be overwriten.
// After that, the _content.yaml will be sorted descending by time (created) and saved.
func AddBlogPage(folder, pagefile string) (objx.Map, error) {
	// index.md von der Verarbeitung ausschliessen
	// die aktuelle Seite als neueste Seite in eine _content.yaml schreiben
	contentFile := filepath.Join(folder, "_content.yaml")
	entries, err := readEntries(contentFile)
	if err != nil {
		return nil, err
	}

	entries = slices.DeleteFunc(entries, func(e BlogEntry) bool {
		return e.Name == pagefile
	})

	entry := BlogEntry{
		Name:    pagefile,
		Created: time.Now(),
	}
	entries = append(entries, entry)

	slices.SortFunc(entries, func(a, b BlogEntry) int {
		return b.Created.Compare(a.Created)
	})

	err = writeEntries(contentFile, entries)
	if err != nil {
		return nil, err
	}
	return objx.Map{"created": entry.Created}, nil
}

func readEntries(file string) ([]BlogEntry, error) {
	entries := make([]BlogEntry, 0)
	if ok, _ := utils.FileExists(file); ok {
		contentYaml, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(contentYaml, &entries)
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}

func writeEntries(file string, be []BlogEntry) error {
	contentYaml, err := yaml.Marshal(be)
	if err != nil {
		return err
	}
	return os.WriteFile(file, contentYaml, 0777)
}
