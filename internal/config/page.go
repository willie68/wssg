package config

// Page the configuration of a single page, used by frontmatter
type Page struct {
	Title     string `yaml:"title"`
	Name      string `yaml:"name"`
	Processor string `yaml:"processor"`
}

var (
	// PageDefault the page default config
	PageDefault = Page{
		Title:     "{{.name}}",
		Name:      "{{.name}}",
		Processor: ProcMarkdown,
	}
)

// General converting the page sturct into a general
func (p Page) General() (output General) {
	output = make(General)
	output["title"] = p.Title
	output["name"] = p.Name
	output["processor"] = p.Processor
	return
}
