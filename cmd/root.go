/*
Copyright © 2023 Scott Brenner <scott@scottbrenner.me>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sm",
	Short: "Download & update StepMania simfile packs",
	Long:  "A command-line interface for downloading and updating StepMania simfile packs",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
