package mdtohtml

import (
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/plugins"
)

// Md2HTML markdown plugin for converting a markdown file to html
type Md2HTML struct {
}

// New create a new markdown plugin
func New() plugins.Plugin {
	return &Md2HTML{}
}

// CreateBody interface method to create a html body from a markdown file
func (m *Md2HTML) CreateBody(content []byte, _ model.Page) (*plugins.Response, error) {
	// extract md
	ignore := make(objx.Map)
	md, err := frontmatter.Parse(strings.NewReader(string(content)), &ignore)
	if err != nil {
		return nil, err
	}
	// convert md to html
	ht := mdToHTML(md)
	return &plugins.Response{
		Body: string(ht),
	}, nil
}

// HTMLTemplateName returning the used html template
func (m *Md2HTML) HTMLTemplateName() string {
	return "layout.html"
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.LazyLoadImages
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
