/*
Copyright © 2023 Scott Brenner <scott@scottbrenner.me>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sm functionality",
	Long: `Initialize sm functionality in a pack's directory.

Prompts for pack's download source and stores it in a text file
to be referenced for future operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		initPackSource()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func getPackSource() (downloadURL string, err error) {
	_, err = fmt.Print("Enter download source for pack: ")
	if err != nil {
		return "", err
	}

	_, err = fmt.Scanln(&downloadURL)
	if err != nil {
		return "", err
	}

	return downloadURL, nil
}

func setPackSource(downloadURL string) (err error) {
	f, err := os.Create("source.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(downloadURL)
	if err != nil {
		return err
	}
	f.Sync()

	return nil
}

func initPackSource() (err error) {
	sourceURL, err := getPackSource()
	if err != nil {
		return err
	}

	err = setPackSource(sourceURL)
	if err != nil {
		return err
	}

	return nil
}
