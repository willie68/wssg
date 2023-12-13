/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
)

// sectionCmd represents the section command
var sectionCmd = &cobra.Command{
	Use:   "section <name>",
	Short: "creates a new section",
	Long: `creates a new section with that name. 
	In the folder there will be a new folder with that name and 
	some default configurations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		return CreateSection(rootFolder, args, force)
	},
}

func init() {
	newCmd.AddCommand(sectionCmd)
	sectionCmd.Flags().BoolP("force", "f", false, "force reinitialise. Page content maybe overwritten.")
}

func CreateSection(rootFolder string, args []string, force bool) error {
	log := logging.New().WithName("newSection")
	config.LoadSite(rootFolder)
	if len(args) == 0 {
		return errors.New("missing section name")
	}
	name := args[0]
	log.Infof("creating a new site in folder \"%s\" with name: %s", rootFolder, name)
	sectionFolder := filepath.Join(rootFolder, name)
	ok, err := fileExists(sectionFolder)
	if err != nil {
		return err
	}
	if ok && !force {
		return errors.New("site already exists")
	}
	err = os.MkdirAll(sectionFolder, 755)
	if err != nil {
		return err
	}
	// create config folder
	configFolder := filepath.Join(sectionFolder, config.WssgFolder)
	err = os.MkdirAll(configFolder, 755)
	if err != nil {
		return err
	}
	// generate default section config
	sectionDefault := config.Section{
		SectionName:  name,
		SectionTitle: name,
		Processor:    config.ProcInternal,
	}.General()
	sectionConfigFile := filepath.Join(configFolder, config.SectionFileName)
	err = writeAsYaml(sectionConfigFile, sectionDefault)
	return err
}
