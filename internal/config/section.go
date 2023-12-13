package config

const SectionFileName = "section.yaml"

type Section struct {
	Name      string `yaml:"name"`
	Title     string `yaml:"title"`
	Processor string `yaml:"processor"`
}

var SectionDefault = Section{
	Name:      "{{.name}}",
	Title:     "{{.name}}",
	Processor: ProcInternal,
}

func (s Section) General() (output General) {
	output = make(General)
	output["name"] = s.Name
	output["title"] = s.Title
	output["processor"] = s.Processor
	return
}
