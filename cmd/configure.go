/*
Copyright © 2023 Scott Brenner <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure sm-cli",
	Long:  "Configure sm-cli to interact with your local filesystem",
	Run: func(cmd *cobra.Command, args []string) {
		packSource()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.PersistentFlags().String("sources", "", "Configure sources for local packs")
}

func getPacks() (packs []fs.DirEntry) {
	packs, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	var packSlice []fs.DirEntry

	for _, file := range packs {
		if file.IsDir() {
			packSlice = append(packSlice, file)
		}
	}

	return packSlice
}

func checkPacks(packSlice []fs.DirEntry) {
	for _, file := range packSlice {
		songFolder, err := os.Open(file.Name())
		if err != nil {
			log.Fatalf("failed opening directory: %s", err)
		}
		defer songFolder.Close()

		numSimFiles := 0

		list, _ := songFolder.Readdirnames(0)
		for _, name := range list {
			if strings.Contains(name, ".sm") {
				numSimFiles += 1
			}
		}
		if numSimFiles == 0 {
			log.Printf("No simfiles detected in directory: %s", file.Name())
		}
	}
}

func promptPackSource(packSlice []fs.DirEntry) (sources map[string]string) {
	sourceMap := make(map[string]string)

	for _, file := range packSlice {
		var packURL string
		fmt.Printf("Enter download source for pack %s: ", file.Name())
		fmt.Scanln(&packURL)
		sourceMap[file.Name()] = packURL
	}

	return sourceMap
}

func setPackSource(sourceMap map[string]string) {
	for pack, source := range sourceMap {
		sourceFile := fmt.Sprintf("%s/%s", pack, "source.txt")
		f, err := os.Create(sourceFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = f.WriteString(source)
		if err != nil {
			log.Fatal(err)
		}
		f.Sync()
	}
}

func packSource() {
	packs := getPacks()
	checkPacks(packs)
	// packSources := promptPackSource(packs)
	// setPackSource(packSources)
}
