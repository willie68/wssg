package config

const SectionFileName = "section.yaml"

type Section struct {
	SectionName  string `yaml:"sectionname"`
	SectionTitle string `yaml:"sectiontitle"`
	Processor    string `yaml:"processor"`
}

var SectionDefault = Section{
	SectionName:  "{{.sectionname}}",
	SectionTitle: "{{.sectionname}}",
	Processor:    ProcInternal,
}

func (s Section) General() (output General) {
	output = make(General)
	output["sectionname"] = s.SectionName
	output["sectiontitle"] = s.SectionTitle
	output["processor"] = s.Processor
	return
}
