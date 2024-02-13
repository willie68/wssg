package templates

import _ "embed"

const (
	ImageEntry = "<div style=\"display: inline-block;overflow: hidden;width:{{.thumbswidth}}px;padding: 5px 5px 5px 5px;\"><a href=\"{{.source}}\" target=\"_blank\"><img loading=\"lazy\" src=\"{{.thumbnail}}\" alt=\"{{.name}}\"><span>{{.name}}%s</span></a></div><br/>"
	ImageTag   = "<br/>%[1]s: {{.%[1]s}}"
)

// Page the page template
var (
	//go:embed gallery/page.md
	GalleryPage string
	//go:embed gallery/style.css
	GalleryStyle string
	//go:embed gallery/style_fluid.css
	GalleryFluidStyle string
)
