package mdtohtml

import (
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/plugins"
)

// Md2HTML internal plugin for converting a markdown file to html
type Md2HTML struct {
}

// New create a new internal plugin
func New() plugins.Plugin {
	return &Md2HTML{}
}

// CreateBody interface method to create a html body from a markdown file
func (m *Md2HTML) CreateBody(content []byte, pg model.Page) ([]byte, error) {

	// extract md
	ignore := make(config.General)
	md, err := frontmatter.Parse(strings.NewReader(string(content)), &ignore)
	if err != nil {
		return nil, err
	}
	// convert md to html
	ht := mdToHTML(md)
	return ht, nil
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
