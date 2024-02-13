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
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/plugins"
	"github.com/willie68/wssg/internal/plugins/gallery"
	"github.com/willie68/wssg/internal/plugins/mdtohtml"
	"github.com/willie68/wssg/internal/plugins/plain"
	"github.com/willie68/wssg/internal/utils"
	"github.com/willie68/wssg/templates"
	"gopkg.in/yaml.v3"
)

const (
	rootSection = "_root"
)

// Generator this the main generator engine
type Generator struct {
	rootFolder string
	force      bool
	genConfig  config.Generate
	siteConfig config.Site
	sections   map[string]config.General
	pages      []model.Page
	log        *logging.Logger
	refreshed  bool
	autoreload bool
}

// New creates a new initialised generator
func New(rootFolder string, force bool) Generator {
	root, err := filepath.Abs(rootFolder)
	if err != nil {
		logging.Root.Errorf("wrong format for root folder: %s \r\n %v", rootFolder, err)
		os.Exit(-1)
	}
	wssgfld := filepath.Join(root, config.WssgFolder)
	if _, err := os.Stat(wssgfld); err != nil {
		logging.Root.Errorf("folder is not an wssg root folder: %s \r\n %v", rootFolder, err)
		os.Exit(-1)
	}
	g := Generator{
		rootFolder: root,
		force:      force,
		log:        logging.New().WithName("generator"),
	}
	g.init()
	return g
}

func (g *Generator) init() {
	g.sections = make(map[string]config.General)
	g.siteConfig = config.LoadSite(g.rootFolder)
	g.genConfig = config.LoadGenConfig(g.rootFolder)
	if g.autoreload {
		g.genConfig.Autoreload = templates.AutoreloadJS
	}
	g.pages = make([]model.Page, 0)
	g.genConfig.Force = g.force
}

// WithAutoreload adapt autorelaod
func (g *Generator) WithAutoreload() {
	g.autoreload = true
}

// SiteConfig return the configuration of the site
func (g *Generator) SiteConfig() config.Site {
	return g.siteConfig
}

// GenConfig return the configuration of the generator
func (g *Generator) GenConfig() config.Generate {
	return g.genConfig
}

// Execute walk thru the folders and register section/pages. After that processing each file.
func (g *Generator) Execute() error {
	g.log.Debug("init")
	g.init()
	g.log.Debug("prepare")
	err := g.prepare()
	if err != nil {
		g.log.Errorf("error prepare site: %V", err)
		return err
	}

	g.log.Debugf("process pages: %d", len(g.pages))
	for _, pg := range g.pages {
		err := g.processPage(pg)
		if err != nil {
			g.log.Errorf("error processing site: %V", err)
			return err
		}
	}
	g.refreshed = true
	g.log.Debug("finished")
	return nil
}

func (g *Generator) IsRefreshed() bool {
	rf := g.refreshed
	g.refreshed = false
	return rf
}
func (g *Generator) prepare() error {
	return filepath.Walk(g.rootFolder, g.doWalk)
}

func (g *Generator) doWalk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return nil
	}
	name := ""
	if info != nil {
		name = info.Name()
	}
	g.log.Debugf("walk: %s", path)
	section := strings.ReplaceAll(path, "\\", "/")
	if path == section || name == "" {
		return nil
	}
	// skip directories with . prefix
	if strings.HasPrefix(name, ".") && info.IsDir() {
		return filepath.SkipDir
	}
	// skip files with . and _ prefix
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
		return nil
	}
	rootPath := strings.ReplaceAll(g.rootFolder, "\\", "/")
	section = strings.TrimPrefix(section, rootPath)
	sections := strings.Split(section, "/")
	if info.IsDir() {
		return nil
	}
	section = strings.Join(sections[1:len(sections)-1], "/")
	if g.isTemplate(name) {
		err := g.registerPage(section, path, info)
		if err != nil {
			g.log.Errorf("error registering page \"%s/%s\": %v", section, name, err)
		}
	} else {
		// copy as static file to output
		err := g.copy2Output(section, path, info)
		if err != nil {
			g.log.Errorf("error copying file: %v", err)
		}
	}
	return nil
}

func (g *Generator) isTemplate(name string) bool {
	if strings.HasSuffix(strings.ToLower(name), ".md") {
		return true
	}
	if strings.HasSuffix(strings.ToLower(name), ".html") {
		return true
	}
	return false
}

// registerPage this will only process the page config and cache information about the page
func (g *Generator) registerPage(section string, path string, info os.FileInfo) error {
	g.log.Debugf("start register page: %s/%s", section, info.Name())
	secCnf := g.getSectionConfig(section)
	//g.log.Debugf("used config: %v", secCnf)
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
	proc, ok := secCnf["processor"].(string)
	if !ok {
		proc = config.ProcMarkdown
	}
	// process pageCnf
	defaults := make(config.General)
	name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
	title := name
	if title == "index" {
		title = secCnf["title"].(string)
	}
	defaults["name"] = name
	defaults["processor"] = proc
	defaults["title"] = title
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
	srcFolder := filepath.Dir(path)
	sections := strings.Split(section, "/")
	dstFolder := filepath.Join(g.rootFolder, g.genConfig.Output, filepath.Join(sections...))
	pg := &model.Page{
		Name:         pageCnf["name"].(string),
		Title:        pageCnf["title"].(string),
		Filename:     info.Name(),
		Section:      section,
		Path:         path,
		Cnf:          pageCnf,
		Order:        order,
		Processor:    pageCnf["processor"].(string),
		SourceFolder: srcFolder,
		DestFolder:   dstFolder,
	}
	pg = g.pageURLPath(pg)
	g.pages = append(g.pages, *pg)
	return nil
}

func (g *Generator) pageURLPath(pg *model.Page) *model.Page {
	if pg.Section == "" || pg.Section == rootSection {
		pg.URLPath = fmt.Sprintf("%s.html", pg.Name)
		return pg
	}
	pg.URLPath = fmt.Sprintf("/%s/%s.html", pg.Section, pg.Name)
	return pg
}

// processPage will now generate the desired html file
func (g *Generator) processPage(pg model.Page) error {
	g.log.Debugf("start processing page: %s/%s (%s)", pg.Section, pg.Name, pg.Filename)
	secCnf := g.getSectionConfig(pg.Section)
	//g.log.Debugf("used config: %v", secCnf)

	pg.Cnf["page"] = pg
	pg.Cnf["section"] = secCnf
	pg.Cnf["site"] = g.siteConfig
	pages := g.filterSortPages(pg.Section)
	pg.Cnf["pages"] = pages
	pg.Cnf["sections"] = g.filterSortSections()
	pg.Cnf["generator"] = g.genConfig

	// load file
	dt, err := os.ReadFile(pg.Path)
	if err != nil {
		return err
	}
	var processor plugins.Plugin

	switch pg.Processor {
	case gallery.PluginName:
		processor = gallery.New(pg.Cnf)
	case config.ProcMarkdown:
		processor = mdtohtml.New()
	default:
		processor = plain.New()
	}

	// now process page with plugin
	// set converted md as body
	res, err := processor.CreateBody(dt, pg)
	if err != nil {
		return err
	}
	pg.Cnf["body"] = res.Body
	pg.Cnf["style"] = res.Style
	pg.Cnf["script"] = res.Script
	banner := ""
	if bn, ok := pg.Cnf["cookiebanner"].(config.General); ok {
		if en, ok := bn["enabled"].(bool); ok && en {
			banner = templates.Cookiebanner
			if txt, ok := bn["text"].(string); !ok || txt == "" {
				bn["text"] = templates.CookiebannerText
			}
		}
	}
	pg.Cnf["cbanner"] = banner

	// load html layout
	//TODO layout.html should be in the site config
	layFile := filepath.Join(g.rootFolder, config.WssgFolder, processor.HTMLTemplateName())
	layout, err := os.ReadFile(layFile)
	if err != nil {
		return err
	}
	ht, err := g.mergeHTML(string(layout), pg.Cnf)
	if err != nil {
		return err
	}
	// write html to output
	var destPath string
	sections := strings.Split(pg.Section, "/")
	destPath = filepath.Join(g.rootFolder, g.genConfig.Output, filepath.Join(sections...))
	err = os.MkdirAll(destPath, 755)
	if err != nil {
		return err
	}
	pageHTMLFile := filepath.Join(destPath, fmt.Sprintf("%s.html", pg.Name))
	err = os.WriteFile(pageHTMLFile, ht, 0775)
	return err
}

// filterSortPages getting all pages, filtering the actual and unvisible pages, than sorting in index order
func (g *Generator) filterSortPages(sec string) []model.Page {
	ps := make([]model.Page, 0)
	for _, pg := range g.pages {
		if pg.Section == sec {
			ps = append(ps, pg)
		}
	}
	sort.Slice(ps, func(i, j int) bool {
		// less function
		return ps[i].Order < ps[j].Order
	})
	return ps
}

// filterSortSections getting all section names, filter actual, root and unvisible sections, sorting in index order
func (g *Generator) filterSortSections() []config.Section {
	sl := make([]config.Section, 0)
	for key, sec := range g.sections {
		if !strings.HasPrefix(key, "_") {
			sc := config.G2Section(sec)
			sl = append(sl, sc)
		}
	}
	sort.Slice(sl, func(i, j int) bool {
		// less function
		if sl[i].Order > 0 || sl[j].Order > 0 {
			return sl[i].Order < sl[j].Order
		}
		return sl[i].Name < sl[j].Name
	})
	return sl
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
	err = mergo.Merge(&pageCnf, secCnf)
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
		section = rootSection
	}
	if ok, _ := utils.FileExists(sectionFile); ok {
		err := utils.LoadYAML(sectionFile, &cnf)
		if err != nil {
			g.log.Errorf("error loading section file: %v", err)
		}
	}
	cnf["site"] = g.siteConfig
	err := mergo.Merge(&cnf, g.siteConfig.UserProperties)
	if err != nil {
		g.log.Errorf("error merging section config with site config: %v", err)
	}
	g.sections[section] = cnf
	return cnf
}

func (g *Generator) getRegisteredPageCnf(name string) (*model.Page, bool) {
	for _, v := range g.pages {
		if v.Name == name {
			return &v, true
		}
	}
	return nil, false
}
