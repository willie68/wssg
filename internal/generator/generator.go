package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
	"github.com/willie68/wssg/internal/utils"
	"gopkg.in/yaml.v3"
)

type Generator struct {
	rootFolder string
	force      bool
	genConfig  config.Generate
	sideConfig config.Site
	sections   map[string]config.General
	pages      []page
	log        *logging.Logger
}

type page struct {
	Name     string         `json:"name" yaml:"name"`
	Title    string         `json:"title" yaml:"title"`
	filename string         `json:"filename" yaml:"filename"`
	section  string         `json:"section" yaml:"section"`
	Path     string         `json:"path" yaml:"path"`
	cnf      config.General `json:"cnf" yaml:"cnf"`
	URLPath  string         `json:"urlpath" yaml:"urlpath"`
	Order    int            `json:"order" yaml:"order"`
}

var ()

func New(rootFolder string, force bool) Generator {
	g := Generator{
		rootFolder: rootFolder,
		force:      force,
		log:        logging.New().WithName("generator"),
		pages:      make([]page, 0),
	}
	g.init()
	return g
}

func (g *Generator) init() {
	g.sections = make(map[string]config.General)
	g.sideConfig = config.LoadSite(g.rootFolder)
	g.genConfig = config.LoadGenConfig(g.rootFolder)
}

func (g *Generator) Execute() error {
	g.init()
	err := filepath.Walk(g.rootFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == g.rootFolder {
				return nil
			}
			name := info.Name()
			if strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			section := strings.ReplaceAll(path, "\\", "/")
			sections := strings.Split(section, "/")
			if info.IsDir() {
				return nil
			}
			fmt.Println(path, info.Name(), info.Size())
			section = strings.Join(sections[1:len(sections)-1], "/")
			if g.isTemplate(name) {
				g.registerPage(section, path, info)
			} else {
				// copy as static file to output
				g.copy2Output(section, path, info)
			}
			return nil
		})
	if err != nil {
		g.log.Errorf("error processing site: %V", err)
		return err
	}
	for _, pg := range g.pages {
		err = g.processPage(pg)
		if err != nil {
			g.log.Errorf("error processing site: %V", err)
			return err
		}
	}
	return nil
}

func (g *Generator) isTemplate(name string) bool {
	if strings.HasSuffix(strings.ToLower(name), ".md") {
		return true
	}
	return false
}

// registerPage this will only process the page config and cache information about the page
func (g *Generator) registerPage(section string, path string, info os.FileInfo) error {
	g.log.Debugf("start processing file: %s", info.Name())
	secCnf := g.getSectionConfig(section)
	g.log.Debugf("used config: %v", secCnf)
	dt, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// extract front matter yaml and md
	pageCnf := make(config.General)
	_, err = frontmatter.Parse(strings.NewReader(string(dt)), &pageCnf)
	if err != nil {
		return err
	}

	// process pageCnf
	defaults := make(config.General)
	defaults["name"] = strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
	defaults["processor"] = config.ProcInternal
	defaults["title"] = defaults["name"]
	err = mergo.Merge(&pageCnf, defaults)
	if err != nil {
		return err
	}
	pageCnf, err = g.processPageCnf(pageCnf, secCnf)
	if err != nil {
		return err
	}
	order, ok := pageCnf["order"].(int)
	if !ok {
		order = 0
	}
	pg := &page{
		Name:     pageCnf["name"].(string),
		Title:    pageCnf["title"].(string),
		filename: info.Name(),
		section:  section,
		Path:     path,
		cnf:      pageCnf,
		Order:    order,
	}
	pg = g.pageURLPath(pg)
	g.pages = append(g.pages, *pg)
	return nil
}

func (g *Generator) pageURLPath(pg *page) *page {
	pg.URLPath = fmt.Sprintf("%s.html", pg.Name)
	return pg
}

// processPage will now generate the desired html file
func (g *Generator) processPage(pg page) error {
	g.log.Debugf("start processing file: %s", pg.filename)
	secCnf := g.getSectionConfig(pg.section)
	g.log.Debugf("used config: %v", secCnf)
	dt, err := os.ReadFile(pg.Path)
	if err != nil {
		return err
	}

	// extract md
	ignore := make(config.General)
	md, err := frontmatter.Parse(strings.NewReader(string(dt)), &ignore)
	if err != nil {
		return err
	}
	// convert md to html
	ht := mdToHTML(md)
	// set converted md as body
	pg.cnf["body"] = string(ht)
	pg.cnf["page"] = pg
	pg.cnf["section"] = secCnf
	pg.cnf["site"] = g.sideConfig
	pages := g.filterSortPages(pg.section)
	pg.cnf["pages"] = pages

	// load html layout
	//TODO layout.html should be in the site config
	layFile := filepath.Join(g.rootFolder, config.WssgFolder, "layout.html")
	layout, err := os.ReadFile(layFile)
	if err != nil {
		return err
	}
	ht, err = g.mergeHTML(string(layout), pg.cnf)
	if err != nil {
		return err
	}
	// write html to output
	var destPath string
	sections := strings.Split(pg.section, "/")
	destPath = filepath.Join(g.rootFolder, g.genConfig.Output, filepath.Join(sections...))
	err = os.MkdirAll(destPath, 755)
	if err != nil {
		return err
	}
	pageHTMLFile := filepath.Join(destPath, fmt.Sprintf("%s.html", pg.Name))
	err = os.WriteFile(pageHTMLFile, ht, 0775)
	return err
}

func (g *Generator) filterSortPages(sec string) []page {
	ps := make([]page, 0)
	for _, pg := range g.pages {
		if pg.section == sec {
			ps = append(ps, pg)
		}
	}
	sort.Slice(ps, func(i, j int) bool {
		// less function
		return ps[i].Order < ps[j].Order
	})
	return ps
}

func (g *Generator) mergeHTML(layout string, cnf config.General) ([]byte, error) {
	// merge html template
	tmpl, err := template.New("htmltemplate").Parse(layout)
	if err != nil {
		return nil, err
	}
	var bb bytes.Buffer
	err = tmpl.Execute(&bb, cnf)
	if err != nil {
		return nil, err
	}
	// merge resulting html
	tmpl, err = template.New("htmlpage").Parse(bb.String())
	if err != nil {
		return nil, err
	}
	bb.Reset()
	err = tmpl.Execute(&bb, cnf)
	if err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}

func (g *Generator) processPageCnf(pageCnf config.General, secCnf config.General) (config.General, error) {
	err := mergo.Merge(&pageCnf, config.PageDefault.General())
	if err != nil {
		return nil, err
	}
	ym, err := yaml.Marshal(pageCnf)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("pageConfig").Parse(string(ym))
	if err != nil {
		return nil, err
	}
	var bb bytes.Buffer
	err = tmpl.Execute(&bb, secCnf)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bb.Bytes(), &pageCnf)
	if err != nil {
		return nil, err
	}
	return pageCnf, nil
}

func (g *Generator) copy2Output(section string, path string, info os.FileInfo) error {
	sections := strings.Split(section, "/")
	destPath := filepath.Join(g.rootFolder, g.genConfig.Output, filepath.Join(sections...))

	err := os.MkdirAll(destPath, 755)
	if err != nil {
		return err
	}
	_, err = utils.FileCopy(path, filepath.Join(destPath, info.Name()))
	return err
}

func (g *Generator) getSectionConfig(section string) config.General {
	if cnf, ok := g.sections[section]; ok {
		return cnf
	}
	cnf := make(config.General)
	sections := strings.Split(section, "/")
	sectionFile := filepath.Join(g.rootFolder, filepath.Join(sections...), config.WssgFolder, config.SectionFileName)
	if section == "" {
		section = "_root"
	}
	if ok, _ := utils.FileExists(sectionFile); ok {
		err := utils.LoadYAML(sectionFile, &cnf)
		if err != nil {
			g.log.Errorf("error loading section file: %v", err)
		}
	}
	cnf["site"] = g.sideConfig
	g.sections[section] = cnf
	return cnf
}

func (g *Generator) getRegisteredPageCnf(name string) (*page, bool) {
	for _, v := range g.pages {
		if v.Name == name {
			return &v, true
		}
	}
	return nil, false
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
