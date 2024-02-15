/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willie68/wssg/internal/generator"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate the web site",
	Long:  `generate the web site`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		clean, _ := cmd.Flags().GetBool("force")
		return Generate(rootFolder, force, clean)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().BoolP("force", "f", false, "force build. Unchanged page content will be overwritten.")
	generateCmd.Flags().BoolP("clear", "c", false, "clear output folder. Tha output folder will be delete before.")
}

// Generate creates a new generator and generate the whole site
func Generate(rootFolder string, force, clean bool) error {
	gen := generator.New(rootFolder, force)
	if clean {
		gen.CleanOutput()
	}
	return gen.Execute()
}
