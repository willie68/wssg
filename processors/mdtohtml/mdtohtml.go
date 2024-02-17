package mdtohtml

import (
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/samber/do"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/processors/processor"

	_ "embed"
)

// Page the page template
var (
	//go:embed templates/page.md
	PageMD string
)

// Processor markdown processor for converting a markdown file to html
type Processor struct {
}

func init() {
	proc := New()
	do.ProvideNamedValue[processor.Processor](nil, proc.Name(), proc)
}

// New create a new markdown processor
func New() processor.Processor {
	return &Processor{}
}

// Name returning the name of this processor
func (p *Processor) Name() string {
	return "markdown"
}

// AddPage adding the new page
func (p *Processor) AddPage(folder, pagefile string) (m objx.Map, err error) {
	return
}

// GetPageTemplate getting the right template for the named page
func (p *Processor) GetPageTemplate(name string) string {
	return PageMD
}

// CreateBody interface method to create a html body from a markdown file
func (p *Processor) CreateBody(content []byte, _ model.Page) (*processor.Response, error) {
	// extract md
	ignore := make(objx.Map)
	md, err := frontmatter.Parse(strings.NewReader(string(content)), &ignore)
	if err != nil {
		return nil, err
	}
	// convert md to html
	ht := mdToHTML(md)
	return &processor.Response{
		Body: string(ht),
	}, nil
}

// HTMLTemplateName returning the used html template
func (p *Processor) HTMLTemplateName() string {
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
