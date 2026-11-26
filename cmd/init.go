/*
Copyright © 2026 Scott Brenner <scott@scottbrenner.me>
*/
package cmd

import (
	"os"

	"github.com/ScottBrenner/sm/internal/pack"
	"github.com/spf13/cobra"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sm functionality",
	Long: `Initialize sm functionality in a pack's directory.

Prompts for pack's download source and stores it in a text file
to be referenced for future operations.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		return pack.Init(dir, cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
