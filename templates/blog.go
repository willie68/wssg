package templates

import _ "embed"

// Page the page template
var (
	//go:embed blog/page.md
	BlogPage string
	//go:embed blog/index.md
	BlogIndex string
)
