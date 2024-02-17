package blog

// this is a generator generating a blog with pagination
// every blogentry is a single markdown file. The index.md is the starting page for this.

import (
	"os"
	"path/filepath"
	"slices"
	"time"

	_ "embed"

	"github.com/samber/do"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/utils"
	"github.com/willie68/wssg/processors/processor"
	"gopkg.in/yaml.v3"
)

// Page the page template
var (
	//go:embed templates/page.md
	BlogPage string
	//go:embed templates/index.md
	BlogIndex string
)

// BlogEntry an entry for the blog
type BlogEntry struct {
	Name    string    `yaml:"name"`
	Created time.Time `yaml:"created"`
}

// Processor the blog processor itself
type Processor struct {
}

func init() {
	proc := New()
	do.ProvideNamedValue[processor.Processor](nil, proc.Name(), proc)
}

// New create a new plain processor
func New() processor.Processor {
	return &Processor{}
}

// GetPageTemplate getting the right template for the named page
func (p *Processor) GetPageTemplate(name string) string {
	if name == "index" {
		return BlogIndex
	}
	return BlogPage
}

// AddPage adding the new blog page to the list of pages with the actual time.
// if already there the entry will be overwriten.
// After that, the _content.yaml will be sorted descending by time (created) and saved.
func (p *Processor) AddPage(folder, pagefile string) (objx.Map, error) {
	// index.md von der Verarbeitung ausschliessen
	if pagefile == "index.md" {
		return nil, nil
	}
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

// Name returning the name of this processor
func (b *Processor) Name() string {
	return "blog"
}

// CreateBody interface method to create a html body from a markdown file
func (b *Processor) CreateBody(content []byte, _ model.Page) (*processor.Response, error) {
	return &processor.Response{
		Body: string(content),
	}, nil
}

// HTMLTemplateName returning the used html template
func (b *Processor) HTMLTemplateName() string {
	return "layout.html"
}