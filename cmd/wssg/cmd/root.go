/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willie68/wssg/internal/logging"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "wssg",
		Short: "willies web site generator",
		Long:  `willies web site generator`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logging.Root.SetLevel(logging.Info)
			if verbose {
				logging.Root.SetLevel(logging.Debug)
			}
		},
	}
	rootFolder string
	verbose    bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&rootFolder, "root", "r", "./", "root folder of your website")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
