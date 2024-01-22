package templates

import _ "embed"

// Page the page template
var (
	//go:embed gallery/page.md
	GalleryPage string
	//go:embed gallery/style.css
	GalleryStyle string
	//go:embed gallery/style_fluid.css
	GalleryFluidStyle string
)
