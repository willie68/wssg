/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"

	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"
	"github.com/stretchr/objx"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
	"github.com/willie68/wssg/internal/plugins/gallery"
	"github.com/willie68/wssg/internal/utils"
	"github.com/willie68/wssg/templates"
)

// pageCmd represents the page command
var (
	pageCmd = &cobra.Command{
		Use:   "page [pagename]",
		Short: "add a new page to a section",
		Long: `add a new page to a section.
		It automatically generates a new md file with an example config.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			plugin, _ := cmd.Flags().GetString("plugin")
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			return CreatePage(rootFolder, name, plugin, force)
		},
	}
)

func init() {
	newCmd.AddCommand(pageCmd)
	pageCmd.Flags().BoolP("force", "f", false, "force reinitialise. Page content maybe overwritten.")
	pageCmd.Flags().StringP("plugin", "p", "markdown", "new page with this plugin. Default is markdown.")
}

// CreatePage creates a new page in the site. Name should be prefixed with sections like gallerie/index
func CreatePage(rootFolder string, name string, plugin string, force bool) error {
	log := logging.New().WithName("newPage")
	config.LoadSite(rootFolder)
	sections := make([]string, 0)
	if name == "" {
		name = "index"
	}

	if strings.Contains(name, "/") {
		sections = strings.Split(name, "/")
		name = sections[len(sections)-1]
		sections = sections[:len(sections)-1]
	}

	log.Infof("creating a new page in section \"%v\" with name: %s", sections, name)
	pageFolder := filepath.Join(rootFolder, filepath.Join(sections...))
	if len(sections) > 0 {
		ok := checkSection(pageFolder)
		if !ok {
			return errors.New("new page can only be created in a already created section. Please create a new section before trying to create a new page. \r\n Example: wssg new section <name>")
		}
	}

	pageFile := filepath.Join(pageFolder, fmt.Sprintf("%s.md", name))
	ok, err := utils.FileExists(pageFile)
	if err != nil {
		return err
	}
	if ok && !force {
		return errors.New("page already exists")
	}
	pageGenerateConfig, err := buildPageDefault(name, plugin)
	if err != nil {
		return err
	}

	// Front matters extract page config
	var pageConfig objx.Map
	pageTemplate := templates.PageMD
	if gallery.PluginName == plugin {
		pageTemplate = templates.GalleryPage
		if err != nil {
			return err
		}
	}

	rest, err := frontmatter.Parse(strings.NewReader(pageTemplate), &pageConfig)
	if err != nil {
		return err
	}

	if gallery.PluginName == plugin {
		images, ok := pageConfig["images"].(string)
		if !ok {
			return fmt.Errorf("gallery template broken: %v", err)
		}
		images = filepath.Join(pageFolder, images)
		err = os.MkdirAll(images, os.ModePerm)
		if err != nil {
			return err
		}
	}
	// process config
	err = mergo.Merge(&pageConfig, config.PageDefault.MSA())
	if err != nil {
		return err
	}
	// frontmatter part
	fm, err := yaml.Marshal(pageConfig)
	if err != nil {
		return err
	}
	// process with template engine
	tmpl, err := template.New("page").Parse(string(fm))
	if err != nil {
		return err
	}
	var page bytes.Buffer
	err = tmpl.Execute(&page, pageGenerateConfig)
	if err != nil {
		return err
	}
	// write to site folder
	return os.WriteFile(pageFile, []byte(fmt.Sprintf("---\n%s---\n%s", page.String(), rest)), 0775)
}

func buildPageDefault(name, processor string) (cnf objx.Map, err error) {
	cnf = make(objx.Map)
	err = mergo.Merge(&cnf, config.SiteConfig.MSA())
	if err != nil {
		return nil, err
	}
	cnf["pagename"] = name
	cnf["processor"] = processor
	return cnf, nil
}

func checkSection(sectionFolder string) bool {
	ok, _ := utils.FileExists(sectionFolder)
	if !ok {
		return false
	}
	// create config folder
	configFolder := filepath.Join(sectionFolder, config.WssgFolder)
	ok, _ = utils.FileExists(configFolder)
	if !ok {
		return false
	}
	sectionConfigFile := filepath.Join(configFolder, config.SectionFileName)
	ok, _ = utils.FileExists(sectionConfigFile)
	if !ok {
		return false
	}
	return true
}
