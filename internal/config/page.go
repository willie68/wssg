package config

type Page struct {
	Title     string `yaml:"title"`
	Processor string `yaml:"processor"`
}

var (
	PageDefault = Page{
		Title:     "{{.title}}",
		Processor: ProcInternal,
	}
)

func (p Page) General() (output General) {
	output = make(General)
	output["title"] = p.Title
	output["processor"] = p.Processor
	return
}
