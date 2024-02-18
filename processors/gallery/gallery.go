package gallery

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	txttpl "text/template"

	_ "embed"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"github.com/samber/do"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/utils"
	"github.com/willie68/wssg/processors/mdtohtml"
	"github.com/willie68/wssg/processors/processor"
	"gopkg.in/yaml.v3"
)

const (
	ImageEntry = "<div style=\"display: inline-block;overflow: hidden;width:{{.thumbswidth}}px;padding: 5px 5px 5px 5px;\"><a href=\"{{.source}}\" target=\"_blank\"><img loading=\"lazy\" src=\"{{.thumbnail}}\" alt=\"{{.name}}\"><span>{{.name}}%s</span></a></div><br/>"
	ImageTag   = "<br/>%[1]s: {{.%[1]s}}"
)

// Page the page template
var (
	//go:embed templates/page.md
	GalleryPage string
	//go:embed templates/style.css
	GalleryStyle string
	//go:embed templates/style_fluid.css
	GalleryFluidStyle string

	exts = []string{".jpeg", ".jpg", ".bmp", ".png"}
)

// Processor struct for the processor
type Processor struct {
	log *logging.Logger
}

type galleryPage struct {
	name          string
	width         int
	force         bool
	crop          bool
	fluid         bool
	imgFolder     string
	dstFolder     string
	images        []img
	templateImage *template.Template
}

type img struct {
	Name           string
	Source         string
	Thumbnail      string
	Size           int64
	UserProperties map[string]string
}

func init() {
	proc := New()
	do.ProvideNamedValue[processor.Processor](nil, proc.Name(), proc)
}

// New creating a new gallery processor
func New() processor.Processor {
	return &Processor{
		log: logging.New().WithName("gallery"),
	}
}

// Name returning the name of this processor
func (p *Processor) Name() string {
	return "gallery"
}

// AddPage adding the new page
func (p *Processor) AddPage(folder, pagefile string) (m objx.Map, err error) {
	return
}

// GetPageTemplate getting the right template for the named page
func (p *Processor) GetPageTemplate(name string) string {
	return GalleryPage
}

// CreateBody creating ths body for this gallery page
func (p *Processor) CreateBody(content []byte, pg model.Page) (*processor.Response, error) {
	// getting all image file names
	imgFld := pg.Cnf.Get("images").String()
	if imgFld == "" {
		return nil, errors.New("can't determine image folder")
	}
	if !filepath.IsAbs(imgFld) {
		imgFld = filepath.Join(pg.SourceFolder, imgFld)
	}
	g := galleryPage{
		name: pg.Name,
	}
	g.imgFolder = imgFld

	imgProps := pg.Cnf["imageproperties"]
	props := utils.ConvertArrIntToArrString(imgProps)
	imgs, err := p.prepareImageList(pg, g, props)
	if err != nil {
		return nil, err
	}
	g.images = imgs
	g.dstFolder = filepath.Join(pg.DestFolder, "images", g.name)
	err = p.ensureImageCopy(g)
	if err != nil {
		return nil, err
	}

	// generating thumbs in output folder
	g.width = pg.Cnf.Get("thumbswidth").Int(100)
	g.crop = pg.Cnf.Get("crop").Bool(false)
	g.force = false
	if genCnf, ok := pg.Cnf["generator"].(config.Generate); ok {
		g.force = genCnf.Force
	}
	g.fluid = pg.Cnf.Get("fluid").Bool(false)

	p.log.Info("generating thumbs")

	p.generateThumbs(g)

	// generating the gallery page with htmx
	tplImgEntry := pg.Cnf.Get("imageentry").String()
	if tplImgEntry == "" {
		var b bytes.Buffer
		for _, property := range props {
			_, _ = b.WriteString(fmt.Sprintf(ImageTag, property))
		}
		tplImgEntry = fmt.Sprintf(ImageEntry, b.String())
		p.log.Infof("page %s: using build in image entry template", pg.Name)
	}
	tplImg, err := template.New("galleryentry").Parse(tplImgEntry)
	if err != nil {
		return nil, err
	}
	g.templateImage = tplImg

	imgContainer := pg.Cnf.Get("imagecontainer").Str("{{ .images }}")
	tplImgContainer, err := txttpl.New("gallerycontainer").Parse(imgContainer)
	if err != nil {
		return nil, err
	}
	var imagesHTML string
	if g.fluid {
		imagesHTML, err = p.writeFluidImageHTMLList(g)
	} else {
		imagesHTML, err = p.writeImageHTMLList(g)
	}
	if err != nil {
		return nil, err
	}
	var bc bytes.Buffer
	m := make(map[string]any)
	m["images"] = imagesHTML
	err = tplImgContainer.Execute(&bc, m)
	if err != nil {
		return nil, err
	}
	pg.Cnf["images"] = bc.String()

	// extract md
	res, err := mdtohtml.New().CreateBody(content, pg)
	res.Style = GalleryStyle
	if g.fluid {
		res.Style = GalleryFluidStyle
	}
	res.Style = pg.Cnf.Get("style").Str(res.Style)
	res.Script = ""
	if err != nil {
		return nil, err
	}
	res.Render = true
	return res, nil
}

func (p *Processor) writeImageHTMLList(g galleryPage) (string, error) {
	var b bytes.Buffer
	for _, ig := range g.images {
		m := p.makeImageMap(g, ig)

		var bb bytes.Buffer
		err := g.templateImage.Execute(&bb, m)
		if err != nil {
			return "", err
		}
		_, err = b.WriteString(bb.String() + "\r\n")
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func (p *Processor) writeFluidImageHTMLList(g galleryPage) (string, error) {
	colCount := 3

	orderedImgs := make([][]int, 0)
	for range colCount {
		orderedImgs = append(orderedImgs, make([]int, 0))
	}
	for i := range g.images {
		x := i % colCount
		orderedImgs[x] = append(orderedImgs[x], i)
	}

	var b bytes.Buffer
	_, err := b.WriteString("<div class=\"galrow\">\r\n  <div class=\"galcolumn\">\r\n")
	if err != nil {
		return "", err
	}
	for x := range orderedImgs {
		for y := range orderedImgs[x] {
			m := p.makeImageMap(g, g.images[orderedImgs[x][y]])

			var bb bytes.Buffer
			err := g.templateImage.Execute(&bb, m)
			if err != nil {
				return "", err
			}
			_, err = b.WriteString(bb.String() + "\r\n")
			if err != nil {
				return "", err
			}
		}
		_, err := b.WriteString("  </div>\r\n  <div class=\"galcolumn\">\r\n")
		if err != nil {
			return "", err
		}
	}
	_, err = b.WriteString("  </div>\r\n</div>\r\n")
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (p *Processor) generateThumbs(g galleryPage) {
	var wg sync.WaitGroup
	for _, i := range g.images {
		wg.Add(1)
		img := i
		go func() {
			defer wg.Done()
			err := p.creatThumb(g.dstFolder, img, g.width, g.force, g.crop)
			if err != nil {
				p.log.Errorf("error creating thumbnail: %s, %v", img.Name, err)
			}
		}()
	}
	wg.Wait()
}

func (p *Processor) makeImageMap(g galleryPage, i img) map[string]string {
	m := make(map[string]string)
	m["name"] = utils.FileNameWOExt(i.Name)
	m["source"] = fmt.Sprintf("images/%s/%s", g.name, i.Source)
	m["thumbnail"] = fmt.Sprintf("images/%s/%s", g.name, i.Thumbnail)
	m["sizebytes"] = fmt.Sprintf("%d", i.Size)
	m["size"] = utils.ByteCountBinary(i.Size)
	m["thumbswidth"] = fmt.Sprintf("%d", g.width)
	for k, v := range i.UserProperties {
		m[k] = v
	}
	return m
}

func (p *Processor) ensureImageCopy(g galleryPage) error {
	err := os.MkdirAll(g.dstFolder, os.ModePerm)
	if err != nil {
		return err
	}

	for _, de := range g.images {
		// copy to output folder
		err = p.ensureCopy(g.imgFolder, g.dstFolder, de.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) prepareImageList(pg model.Page, g galleryPage, props []string) ([]img, error) {
	imageDescriptions, err := p.readImageDescription(pg.SourceFolder, pg.Name)
	if err != nil {
		return nil, err
	}
	imgs, err := os.ReadDir(g.imgFolder)
	if err != nil {
		return nil, err
	}
	imgs = p.filterAllowedImages(imgs)
	// First sort all images after name
	slices.SortFunc(imgs, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})
	// Than check if there is another sort order given
	order := utils.ConvertArrIntToArrString(pg.Cnf["imagelist"])
	listonly := pg.Cnf.Get("listonly").Bool(false)
	p.log.Debugf("listonly: %v", listonly)
	images := make([]img, len(order))
	for _, de := range imgs {
		name := utils.FileNameWOExt(de.Name())
		thb := fmt.Sprintf("%s_thb.png", name)
		info, err := de.Info()
		if err != nil {
			p.log.Errorf("can't get file info of %s: %v", de.Name(), err)
		}
		size := int64(0)
		if info != nil {
			size = info.Size()
		}
		i := img{
			Name:           de.Name(),
			Source:         de.Name(),
			Thumbnail:      thb,
			Size:           size,
			UserProperties: getUserproperties(props, imageDescriptions, name),
		}
		if slices.Contains(order, name) {
			pos := slices.Index(order, name)
			images[pos] = i
			continue
		}
		if !listonly {
			images = append(images, i)
		}
	}
	err = p.writeImageDescription(pg.SourceFolder, pg.Name, imageDescriptions)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (p *Processor) filterAllowedImages(imgs []fs.DirEntry) []fs.DirEntry {
	res := make([]fs.DirEntry, 0)
	for _, img := range imgs {
		if !slices.Contains(exts, strings.ToLower(filepath.Ext(img.Name()))) {
			continue
		}
		if strings.HasPrefix(img.Name(), "_") || strings.HasPrefix(img.Name(), ".") {
			continue
		}
		res = append(res, img)
	}
	return res
}

func getUserproperties(props []string, imageDescriptions objx.Map, name string) map[string]string {
	var up map[string]string
	if len(props) > 0 {
		u := imageDescriptions[name]
		if u == nil {
			up = make(map[string]string)
		} else {
			up = utils.ConvertMapIntToMapString(u)
		}
		for _, k := range props {
			_, ok := up[k]
			if !ok {
				up[k] = k
			}
		}
		imageDescriptions[name] = up
	}
	return up
}

func (p *Processor) writeImageDescription(srcFolder, galName string, descs objx.Map) error {
	imgDescription := getImageDescriptionFile(srcFolder, galName)
	if ok, _ := utils.FileExists(imgDescription); ok {
		return nil
	}
	ya, err := yaml.Marshal(descs)
	if err != nil {
		return err
	}
	err = os.WriteFile(imgDescription, ya, 0666)
	return err
}

func (p *Processor) readImageDescription(srcFolder, galName string) (objx.Map, error) {
	descs := make(objx.Map)

	imgDescription := getImageDescriptionFile(srcFolder, galName)
	if ok, _ := utils.FileExists(imgDescription); ok {
		rd, err := os.ReadFile(imgDescription)
		if err == nil {
			err := yaml.Unmarshal(rd, &descs)
			if err != nil {
				return nil, err
			}
		}
	}
	return descs, nil
}

func getImageDescriptionFile(srcFolder, galName string) string {
	return filepath.Join(srcFolder, fmt.Sprintf("_%s.props", galName))
}

func (p *Processor) ensureCopy(imgFolder, dstFolder, name string) error {
	src := filepath.Join(imgFolder, name)
	dst := filepath.Join(dstFolder, name)
	ok, _ := utils.FileExists(dst)
	if !ok {
		_, err := utils.FileCopy(src, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

// HTMLTemplateName returning the used html template
func (p *Processor) HTMLTemplateName() string {
	return "layout.html"
}

func (p *Processor) creatThumb(fld string, i img, width int, force, crop bool) error {
	dst := filepath.Join(fld, i.Thumbnail)
	ok, _ := utils.FileExists(dst)
	if ok && !force {
		return nil
	}
	p.log.Debugf("generating thumb: %s", i.Name)
	src := filepath.Join(fld, i.Source)

	img, err := imgio.Open(src)
	if err != nil {
		return err
	}
	ori, err := orientation(src)
	if err != nil {
		return err
	}
	switch ori {
	case 3: // rotate 180
		img = transform.Rotate(img, 180.0, nil)
	case 6: // rotate 90
		img = transform.Rotate(img, 90.0, nil)
	case 8: // rotate 270
		img = transform.Rotate(img, 270.0, nil)
	}
	bd := img.Bounds()
	height := int(float64(width) * (float64(bd.Dy()) / float64(bd.Dx())))
	if crop {
		height = width
		var rect image.Rectangle
		if bd.Dx() < bd.Dy() {
			delta := bd.Dy() - bd.Dx()
			rect = image.Rect(0, delta/2, bd.Dx(), (delta/2)+bd.Dx())
		} else {
			delta := bd.Dx() - bd.Dy()
			rect = image.Rect(delta/2, 0, (delta/2)+bd.Dy(), bd.Dy())
		}
		img = transform.Crop(img, rect)
	}

	thb := transform.Resize(img, width, height, transform.NearestNeighbor)

	err = imgio.Save(dst, thb, imgio.PNGEncoder())
	return err
}

func orientation(filename string) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return 0, nil
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return 0, nil
	}
	if tag.Count == 1 && tag.Format() == tiff.IntVal {
		orientation, err := tag.Int(0)
		if err != nil {
			return 0, nil
		}
		return orientation, nil
	}
	return 0, nil
}
