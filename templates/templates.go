package templates

import _ "embed"

// Page the page template
var (
	//go:embed page.md
	PageMD string
	//go:embed layout.html
	LayoutHTML string
	//go:embed reload.js
	AutoreloadJS string
)
