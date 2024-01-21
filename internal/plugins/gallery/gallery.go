package gallery

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"image"
	"os"
	"path/filepath"
	"slices"
	"sync"
	txttpl "text/template"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/plugins"
	"github.com/willie68/wssg/internal/plugins/mdtohtml"
	"github.com/willie68/wssg/internal/utils"
	"github.com/willie68/wssg/templates"
	"gopkg.in/yaml.v3"
)

// PluginName name of this plugin
const (
	PluginName = "gallery"
)

var (
	exts = []string{".jpeg", ".jpg", ".bmp", ".png"}
)

// Gallery struct for the processor
type Gallery struct {
	cnf           config.General
	log           *logging.Logger
	width         int
	force         bool
	crop          bool
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

// New creating a new gallery processor
func New(cnf config.General) plugins.Plugin {
	return &Gallery{
		cnf: cnf,
		log: logging.New().WithName("gallery"),
	}
}

// CreateBody creating ths body for this gallery page
func (g *Gallery) CreateBody(content []byte, pg model.Page) (*plugins.Response, error) {
	// getting all image file names
	imgFld, ok := pg.Cnf["images"].(string)
	if !ok {
		return nil, errors.New("can't determine image folder")
	}
	if !filepath.IsAbs(imgFld) {
		imgFld = filepath.Join(pg.SourceFolder, imgFld)
	}
	g.imgFolder = imgFld

	imgProps := pg.Cnf["imgproperties"]
	props := utils.ConvertArrIntToArrString(imgProps)
	imgs, err := g.prepareImageList(props)
	if err != nil {
		return nil, err
	}
	g.images = imgs
	g.dstFolder = filepath.Join(pg.DestFolder, "images")
	err = g.ensureImageCopy()
	if err != nil {
		return nil, err
	}

	// generating thumbs in output folder
	g.width = 100
	if w, ok := pg.Cnf["thumbswidth"].(int); ok {
		g.width = w
	}
	g.crop = false
	if c, ok := pg.Cnf["crop"].(bool); ok {
		g.crop = c
	}
	g.force = false
	if genCnf, ok := pg.Cnf["generator"].(config.Generate); ok {
		g.force = genCnf.Force
	}

	g.log.Info("generating thumbs")

	g.generateThumbs()

	// generating the gallery page with htmx
	tplImgEntry, ok := pg.Cnf["imageentry"].(string)
	if !ok {
		return nil, fmt.Errorf("something wrong with the gallery imageentry on page \"%s\"", pg.Name)
	}
	tplImg, err := template.New("galleryentry").Parse(tplImgEntry)
	if err != nil {
		return nil, err
	}
	g.templateImage = tplImg

	imgContainer, ok := pg.Cnf["imagecontainer"].(string)
	if !ok {
		g.log.Info("no image container given")
		imgContainer = "{{ .images }}"
	}
	tplImgContainer, err := txttpl.New("gallerycontainer").Parse(imgContainer)
	if err != nil {
		return nil, err
	}

	imagesHTML, err := g.writeImageHTMLList()
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
	res.Style = templates.GalleryStyle
	res.Script = ""
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g *Gallery) writeImageHTMLList() (string, error) {
	var b bytes.Buffer
	for _, i := range g.images {
		m := make(map[string]string)
		m["name"] = utils.FileNameWOExt(i.Name)
		m["source"] = fmt.Sprintf("images/%s", i.Source)
		m["thumbnail"] = fmt.Sprintf("images/%s", i.Thumbnail)
		m["sizebytes"] = fmt.Sprintf("%d", i.Size)
		m["size"] = utils.ByteCountBinary(i.Size)
		for k, v := range i.UserProperties {
			m[k] = v
		}

		var bb bytes.Buffer
		err := g.templateImage.Execute(&bb, m)
		if err != nil {
			return "", err
		}
		_, err = b.WriteString(bb.String())
		if err != nil {
			return "", err
		}
		_, err = b.WriteString("\r\n")
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func (g *Gallery) generateThumbs() {
	var wg sync.WaitGroup
	for _, i := range g.images {
		wg.Add(1)
		img := i
		go func() {
			defer wg.Done()
			err := g.creatThumb(g.dstFolder, img, g.width, g.force, g.crop)
			if err != nil {
				g.log.Errorf("error creating thumbnail: %s, %v", img.Name, err)
			}
		}()
	}
	wg.Wait()
}

func (g *Gallery) ensureImageCopy() error {
	err := os.MkdirAll(g.dstFolder, os.ModePerm)
	if err != nil {
		return err
	}

	for _, de := range g.images {
		// copy to output folder
		err = g.ensureCopy(g.imgFolder, g.dstFolder, de.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gallery) prepareImageList(props []string) ([]img, error) {
	imageDescriptions, err := g.readImageDescription(g.imgFolder)
	if err != nil {
		return nil, err
	}
	imgs, err := os.ReadDir(g.imgFolder)
	if err != nil {
		return nil, err
	}
	images := make([]img, 0)
	for _, de := range imgs {
		if !slices.Contains(exts, filepath.Ext(de.Name())) {
			continue
		}
		name := utils.FileNameWOExt(de.Name())
		thb := fmt.Sprintf("%s_thb.png", name)
		info, err := de.Info()
		if err != nil {
			g.log.Errorf("can't get file info of %s: %v", de.Name(), err)
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
		images = append(images, i)
	}
	if len(props) > 0 {
		err = g.writeImageDescription(g.imgFolder, imageDescriptions)
		if err != nil {
			return nil, err
		}
	}
	return images, nil
}

func getUserproperties(props []string, imageDescriptions config.General, name string) map[string]string {
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

func (g *Gallery) writeImageDescription(imgFolder string, descs config.General) error {
	imgDescription := getImageDescriptionFile(imgFolder)
	ya, err := yaml.Marshal(descs)
	if err != nil {
		return err
	}
	err = os.WriteFile(imgDescription, ya, 0666)
	return err
}

func (g *Gallery) readImageDescription(imgFolder string) (config.General, error) {
	descs := make(config.General)

	imgDescription := getImageDescriptionFile(imgFolder)
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

func getImageDescriptionFile(imgFolder string) string {
	return filepath.Join(imgFolder, "_content.yaml")
}

func (g *Gallery) ensureCopy(imgFolder, dstFolder, name string) error {
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
func (g *Gallery) HTMLTemplateName() string {
	return "layout.html"
}

func (g *Gallery) creatThumb(fld string, i img, width int, force, crop bool) error {
	dst := filepath.Join(fld, i.Thumbnail)
	ok, _ := utils.FileExists(dst)
	if ok && !force {
		return nil
	}
	g.log.Debugf("generating thumb: %s", i.Name)
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
