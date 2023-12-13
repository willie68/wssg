package config

type Page struct {
	Title     string `yaml:"title"`
	Name      string `yaml:"name"`
	Processor string `yaml:"processor"`
}

var (
	PageDefault = Page{
		Title:     "{{.name}}",
		Name:      "{{.name}}",
		Processor: ProcInternal,
	}
)

func (p Page) General() (output General) {
	output = make(General)
	output["title"] = p.Title
	output["name"] = p.Name
	output["processor"] = p.Processor
	return
}
