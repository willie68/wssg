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

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
	"github.com/willie68/wssg/internal/model"
	"github.com/willie68/wssg/internal/plugins"
	"github.com/willie68/wssg/internal/plugins/mdtohtml"
	"github.com/willie68/wssg/internal/utils"
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
	cnf config.General
	log *logging.Logger
}

type img struct {
	Name      string
	Source    string
	Thumbnail string
	Size      int64
}

// New creating a new gallery processor
func New(cnf config.General) plugins.Plugin {
	return &Gallery{
		cnf: cnf,
		log: logging.New().WithName("gallery"),
	}
}

// CreateBody creating ths body for this gallery page
func (g *Gallery) CreateBody(content []byte, pg model.Page) ([]byte, error) {
	// getting all image file names
	imgFolder, ok := pg.Cnf["images"].(string)
	if !ok {
		return nil, errors.New("can't determine image folder")
	}
	if !filepath.IsAbs(imgFolder) {
		imgFolder = filepath.Join(pg.SourceFolder, imgFolder)
	}
	imgs, err := os.ReadDir(imgFolder)
	if err != nil {
		return nil, err
	}
	dstFolder := filepath.Join(pg.DestFolder, "images")
	err = os.MkdirAll(dstFolder, os.ModePerm)
	if err != nil {
		return nil, err
	}

	images := make([]img, 0)
	for _, de := range imgs {
		if slices.Contains(exts, filepath.Ext(de.Name())) {
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
				Name:      de.Name(),
				Source:    de.Name(),
				Thumbnail: thb,
				Size:      size,
			}
			images = append(images, i)
			// copy to output folder
			err = g.ensureCopy(imgFolder, dstFolder, de.Name())
			if err != nil {
				return nil, err
			}

		}
	}

	// generating thumbs in output folder
	width := 100
	if w, ok := pg.Cnf["thumbswidth"].(int); ok {
		width = w
	}
	g.log.Info("generating thumbs")
	var wg sync.WaitGroup
	for _, i := range images {
		wg.Add(1)
		img := i
		go func() {
			defer wg.Done()
			err := g.creatThumb(dstFolder, img, width)
			if err != nil {
				g.log.Errorf("error creating thumbnail: %s, %v", img.Name, err)
			}
		}()
	}
	wg.Wait()
	// generating the gallery page with htmx
	var b bytes.Buffer
	tplEntry, ok := pg.Cnf["imageentry"].(string)
	if !ok {
		return nil, fmt.Errorf("something wrong with the gallery imageentry on page \"%s\"", pg.Name)
	}
	tpl, err := template.New("galleryentry").Parse(tplEntry)
	if err != nil {
		return nil, err
	}
	for _, i := range images {
		m := make(map[string]string)
		m["name"] = utils.FileNameWOExt(i.Name)
		m["source"] = fmt.Sprintf("images/%s", i.Source)
		m["thumbnail"] = fmt.Sprintf("images/%s", i.Thumbnail)
		m["sizebytes"] = fmt.Sprintf("%d", i.Size)
		m["size"] = utils.ByteCountBinary(i.Size)

		var bb bytes.Buffer
		err = tpl.Execute(&bb, m)
		if err != nil {
			return nil, err
		}
		b.WriteString(bb.String())
	}
	pg.Cnf["images"] = b.String()

	// extract md
	return mdtohtml.New().CreateBody(content, pg)
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
	return "gallery.html"
}

func (g *Gallery) creatThumb(fld string, i img, width int) error {
	g.log.Debugf("generating thumb: %s", i.Name)
	src := filepath.Join(fld, i.Source)
	img, err := imgio.Open(src)
	if err != nil {
		return err
	}
	var rect image.Rectangle
	bd := img.Bounds()
	if bd.Dx() < bd.Dy() {
		delta := bd.Dy() - bd.Dx()
		rect = image.Rect(0, delta/2, bd.Dx(), (delta/2)+bd.Dx())
	} else {
		delta := bd.Dx() - bd.Dy()
		rect = image.Rect(delta/2, 0, (delta/2)+bd.Dy(), bd.Dy())
	}
	img = transform.Crop(img, rect)

	thb := transform.Resize(img, width, width, transform.NearestNeighbor)

	dst := filepath.Join(fld, i.Thumbnail)
	err = imgio.Save(dst, thb, imgio.PNGEncoder())
	return err
}
