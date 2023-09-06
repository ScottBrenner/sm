/*
Copyright © 2023 Scott Brenner <scott@scottbrenner.me>
*/
package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download pack",
	Long:  "Download pack from the pack's source defined in source.txt",
	Run: func(cmd *cobra.Command, args []string) {
		downloadPack()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

func openSourceFile() (sourceURL string, err error) {
	dat, err := os.ReadFile("source.txt")
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func downloadFromURL(sourceURL string) (err error) {
	fmt.Println("Downloading pack...")
	// Create the file
	out, err := os.Create("pack.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(sourceURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return err
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func unzipDownloadedPack() (err error) {
	fmt.Println("Unpacking pack...")
	dst := "./"
	archive, err := zip.OpenReader("pack.zip")
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, f := range archive.File {
		if !strings.Contains(f.Name, "..") {
			filePath := filepath.Join(dst, f.Name)
			fmt.Println("unzipping file ", filePath)

			if f.FileInfo().IsDir() {
				fmt.Println("creating song folder...")
				os.MkdirAll(filePath, os.ModePerm)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				panic(err)
			}

			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				panic(err)
			}

			fileInArchive, err := f.Open()
			if err != nil {
				panic(err)
			}

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				panic(err)
			}

			dstFile.Close()
			fileInArchive.Close()
		}
	}
	return nil
}

func removeZip() (err error) {
	err = os.Remove("pack.zip")
	if err != nil {
		return err
	}
	return nil
}

func downloadPack() (err error) {
	sourceURL, err := openSourceFile()
	if err != nil {
		return err
	}

	err = downloadFromURL(sourceURL)
	if err != nil {
		return err
	}

	err = unzipDownloadedPack()
	if err != nil {
		return err
	}

	err = removeZip()
	if err != nil {
		return err
	}

	return nil
}
