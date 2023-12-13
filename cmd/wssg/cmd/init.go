/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/willie68/wssg/internal/config"
	"github.com/willie68/wssg/internal/logging"
	"github.com/willie68/wssg/templates"
)

// initCmd represents the init command
var (
	initCmd = &cobra.Command{
		Use:   "init <name>",
		Short: "initialise a new site",
		Long:  `initialise a new site`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initWebsite(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().BoolP("force", "f", false, "force reinitialise. Configs maybe overwritten.")
}

func initWebsite(cmd *cobra.Command, args []string) error {
	log := logging.New().WithName("init")
	// Checking if folder already exists
	if len(args) > 0 {
		rootFolder = filepath.Join(rootFolder, args[0])
	}
	force, _ := cmd.Flags().GetBool("force")
	log.Infof("starting with new website on: \"%s\"", rootFolder)
	err := initFolders(rootFolder, force)
	if err != nil {
		return err
	}
	err = initConfig(rootFolder)
	if err != nil {
		return err
	}
	err = addIndexPage(rootFolder)
	return err
}

func initFolders(rootFolder string, force bool) error {
	if _, err := os.Stat(rootFolder); err != nil {
		if !os.IsNotExist(err) {
			if err != nil {
				return errors.Join(errors.New("path error"), err)
			}
		}
	}
	err := os.MkdirAll(rootFolder, 0755)
	if err != nil {
		return err
	}
	fis, err := os.ReadDir(rootFolder)
	if err != nil {
		return errors.Join(errors.New("path error"), err)
	}
	if len(fis) > 0 && !force {
		return fmt.Errorf("path not empty: %s", rootFolder)
	}

	err = os.MkdirAll(filepath.Join(rootFolder, config.WssgFolder), 0755)
	if err != nil {
		return err
	}
	return nil
}

func initConfig(rootFolder string) error {
	siteConfigDir := filepath.Join(rootFolder, config.WssgFolder)
	siteConfigFile := filepath.Join(siteConfigDir, config.SiteFile)
	err := writeAsYaml(siteConfigFile, config.SiteDefault)
	if err != nil {
		return err
	}

	genConfigFile := filepath.Join(siteConfigDir, config.GenerateFile)
	err = writeAsYaml(genConfigFile, config.GenDefault)
	if err != nil {
		return err
	}

	layoutHTMLFile := filepath.Join(siteConfigDir, "layout.html")
	err = os.WriteFile(layoutHTMLFile, []byte(templates.LayoutHTML), 755)
	return err
}

func addIndexPage(rootFolder string) error {
	return CreatePage(rootFolder, []string{"index"}, true)
}
