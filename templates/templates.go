package templates

import _ "embed"

// Page the page template
var (
	//go:embed page.md
	PageMD string
	//go:embed layout.html
	LayoutHTML string

	//go:embed gallery/page.md
	GalleryPage string
	//go:embed gallery/gallery.html
	GalleryHTML string
)
