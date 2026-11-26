/*
Copyright © 2026 Scott Brenner <scott@scottbrenner.me>
*/
package cmd

import (
	"os"

	"github.com/ScottBrenner/sm/internal/pack"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command.
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download pack",
	Long:  "Download pack from the pack's source defined in source.txt",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		return pack.Download(cmd.Context(), nil, dir, cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
