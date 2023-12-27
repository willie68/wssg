/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willie68/wssg/internal/generator"
	"github.com/willie68/wssg/internal/logging"
	"github.com/willie68/wssg/internal/server"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "auto generate, watch and start a http server on port 8080",
	Long:  `auto generate, watch and start a http server on port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Serve(rootFolder)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// Serve starting a local http server serving the generated files
func Serve(rootFolder string) error {
	log := logging.New().WithName("serve")
	log.Info("generate web site")
	gen := generator.New(rootFolder, true)
	err := gen.Execute()
	if err != nil {
		return err
	}
	s := server.New(rootFolder, gen)
	return s.Serve()
}
