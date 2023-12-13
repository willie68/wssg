package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
)

type Generator struct {
	rootFolder string
	force      bool
	genConfig  config.Generate
	sideConfig config.General
	sections   map[string]config.General
	log        *logging.Logger
}

var ()

func New(rootFolder string, force bool) Generator {
	g := Generator{
		rootFolder: rootFolder,
		force:      force,
		log:        logging.New().WithName("generator"),
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
			if info.IsDir() {
				// create a section config
			}

			fmt.Println(path, info.Name(), info.Size())
			if !info.IsDir() && g.isProceeable(name) {
				g.processFile(info)
			}
			return nil
		})
	if err != nil {
		g.log.Errorf("error rocessing site: %V", err)
		return err
	}
	return nil
}

func (g *Generator) isProceeable(name string) bool {
	if strings.HasSuffix(strings.ToLower(name), ".md") {
		return true
	}
	return false
}

func (g *Generator) processFile(info os.FileInfo) error {
	g.log.Debugf("start processing file: %s", info.Name())
	return nil
}
