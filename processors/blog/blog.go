package blog

// this is a generator generating a blog with pagination
// every blogentry is a single markdown file. The index.md is the starting page for this.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/adrg/frontmatter"
	"github.com/goodsign/monday"
	"github.com/samber/do"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/interfaces"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/utils"
	"github.com/willie68/wssg/processors/mdtohtml"
	"github.com/willie68/wssg/processors/processor"
	"gopkg.in/yaml.v3"
)

// Page the page template
var (
	//go:embed templates/page.md
	BlogPage string
	//go:embed templates/index.md
	BlogIndex string
)

// BlogEntry an entry for the blog
type BlogEntry struct {
	Name    string    `yaml:"name"`
	Created time.Time `yaml:"created"`
}

// Processor the blog processor itself
type Processor struct {
}

func init() {
	proc := New()
	do.ProvideNamedValue[processor.Processor](nil, proc.Name(), proc)
}

// New create a new plain processor
func New() processor.Processor {
	return &Processor{}
}

// GetPageTemplate getting the right template for the named page
func (p *Processor) GetPageTemplate(name string) string {
	if name == "index" {
		return BlogIndex
	}
	return BlogPage
}

// AddPage adding the new blog page to the list of pages with the actual time.
// if already there the entry will be overwriten.
// After that, the _content.yaml will be sorted descending by time (created) and saved.
func (p *Processor) AddPage(folder, pagefile string) (objx.Map, error) {
	// index.md von der Verarbeitung ausschliessen
	if pagefile == "index.md" {
		return nil, nil
	}
	// die aktuelle Seite als neueste Seite in eine _content.yaml schreiben
	contentFile := filepath.Join(folder, "_content.yaml")
	entries, err := readEntries(contentFile)
	if err != nil {
		return nil, err
	}

	entries = slices.DeleteFunc(entries, func(e BlogEntry) bool {
		return e.Name == pagefile
	})

	entry := BlogEntry{
		Name:    pagefile,
		Created: time.Now(),
	}
	entries = append(entries, entry)

	slices.SortFunc(entries, func(a, b BlogEntry) int {
		return b.Created.Compare(a.Created)
	})

	err = writeEntries(contentFile, entries)
	if err != nil {
		return nil, err
	}
	return objx.Map{"created": entry.Created}, nil
}

func readEntries(file string) ([]BlogEntry, error) {
	entries := make([]BlogEntry, 0)
	if ok, _ := utils.FileExists(file); ok {
		contentYaml, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(contentYaml, &entries)
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}

func writeEntries(file string, be []BlogEntry) error {
	contentYaml, err := yaml.Marshal(be)
	if err != nil {
		return err
	}
	return os.WriteFile(file, contentYaml, 0777)
}

// Name returning the name of this processor
func (p *Processor) Name() string {
	return "blog"
}

// CanRenderPage all pages of the processor should be rendered
func (p *Processor) CanRenderPage(pg model.Page) bool {
	return pg.Name == "index"
}

// CreateBody interface method to create a html body from a markdown file
func (p *Processor) CreateBody(content []byte, pg model.Page) (*processor.Response, error) {
	rdr := pg.Name == "index"
	if rdr {
		et := pg.Cnf.Get("entrytemplate").Str("{{.content}}")

		contentFile := filepath.Join(pg.SourceFolder, "_content.yaml")
		entries, err := readEntries(contentFile)
		if err != nil {
			return nil, err
		}
		or := false // descending
		v := pg.Cnf.Get("order").Str("desc")
		if v != "desc" {
			or = true // ascending
		}
		slices.SortFunc(entries, func(a, b BlogEntry) int {
			cmp := b.Created.Compare(a.Created)
			if or {
				cmp = -cmp
			}
			return cmp
		})

		bpp := pg.Cnf.Get("pagination").Int(1)
		pc := 0
		ress := make([]string, 0)
		pg.Cnf["pageCount"] = (len(entries) / bpp) + 1
		if len(entries)%bpp == 0 {
			pg.Cnf["pageCount"] = (len(entries) / bpp)
		}
		pg.Cnf["entryCount"] = len(entries)

		for x, be := range entries {
			mdf := filepath.Join(pg.SourceFolder, be.Name)
			dt, err := os.ReadFile(mdf)
			if err != nil {
				return nil, err
			}
			res, err := p.md2html(dt, pg, be)
			if err != nil {
				return nil, err
			}

			tmpl, err := template.New("entry").Parse(et)
			if err != nil {
				return nil, err
			}

			var bb bytes.Buffer
			bemap := make(objx.Map)
			_, err = frontmatter.Parse(strings.NewReader(string(dt)), &bemap)
			if err != nil {
				return nil, err
			}
			bemap = pg.Cnf.Copy().Merge(bemap)
			bemap["created"] = be.Created
			bemap["content"] = string(res)
			bemap["entrynumber"] = x
			bemap["entryeven"] = (x%2 == 0)

			err = tmpl.Execute(&bb, bemap)
			if err != nil {
				return nil, err
			}

			ress = append(ress, bb.String())

			if x%bpp == (bpp - 1) {
				err := p.savePage(content, pc, ress, pg, x+1 < len(entries))
				if err != nil {
					return nil, err
				}
				pc++
				ress = make([]string, 0)
			}
		}
		if len(ress) > 0 {
			err := p.savePage(content, pc, ress, pg, false)
			if err != nil {
				return nil, err
			}
		}
	}
	return &processor.Response{
		Body:   string(content),
		Render: false,
	}, nil
}

func (p *Processor) md2html(content []byte, pg model.Page, be BlogEntry) (string, error) {
	// extract md
	bemap := make(objx.Map)
	md, err := frontmatter.Parse(strings.NewReader(string(content)), &bemap)
	if err != nil {
		return "", err
	}
	// for macro substitution
	bemap = pg.Cnf.Copy().Merge(bemap)
	bemap["created"] = be.Created

	tmpl, err := template.New("blogmd").Funcs(template.FuncMap{
		"dtFormat": func(dt time.Time, f, l string) string {
			return monday.Format(dt, f, monday.Locale(l))
		},
	}).Parse(string(md))

	if err != nil {
		return "", err
	}
	var bb bytes.Buffer
	err = tmpl.Execute(&bb, bemap)
	if err != nil {
		return "", err
	}

	// convert md to html
	ht := mdtohtml.MdToHTML(bb.Bytes())

	// now process all macros in the html

	return string(ht), nil
}

// HTMLTemplateName returning the used html template
func (p *Processor) HTMLTemplateName() string {
	return "layout.html"
}

func getPageName(pc int) string {
	if pc == 0 {
		return "index"
	}
	return fmt.Sprintf("page%d", pc)
}

func (p *Processor) savePage(content []byte, pc int, ress []string, pg model.Page, hasNext bool) error {
	pg.Cnf["prevPage"] = ""
	pg.Cnf["nextPage"] = ""
	// creating the right page name
	pg.Name = getPageName(pc)
	if pc == 0 {
		pg.Name = "index"
	} else {
		pg.Cnf["prevPage"] = getPageName(pc-1) + ".html"
	}
	if hasNext {
		pg.Cnf["nextPage"] = getPageName(pc+1) + ".html"
	}
	pg.Cnf["actualPage"] = pc + 1

	// merging the blog eintries to one html part
	var bb bytes.Buffer
	for _, res := range ress {
		_, _ = bb.WriteString(res)
	}
	pg.Cnf["blogentries"] = bb.String()

	// converting the md page to html
	md2html := do.MustInvokeNamed[processor.Processor](nil, "markdown")
	res, err := md2html.CreateBody(content, pg)
	if err != nil {
		return err
	}
	// set the html as body for the layout.html
	pg.Cnf["body"] = res.Body

	// generate the final html page
	gen := do.MustInvoke[interfaces.Generator](nil)
	return gen.RenderHTML(p.HTMLTemplateName(), pg)
}
