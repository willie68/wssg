package templates

import _ "embed"

// Page the page template
var (
	//go:embed blog/page.md
	BlogPage string
)
