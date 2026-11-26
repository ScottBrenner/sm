/*
Copyright © 2023 Scott Brenner <scott@scottbrenner.me>
*/
package pack

import (
	"archive/zip"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SourceFileName is the file inside a pack directory that stores its download URL.
const SourceFileName = "source.txt"

const (
	dirMode  os.FileMode = 0o755
	fileMode os.FileMode = 0o644
)

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
		},
	}
}

// ValidateSourceURL trims whitespace and verifies that rawURL is an http or https URL.
func ValidateSourceURL(rawURL string) (string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", errors.New("download source URL cannot be empty")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid download source URL: %w", err)
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", fmt.Errorf("download source URL must be a valid http or https URL: %q", trimmed)
	}

	return trimmed, nil
}

// PromptSource prompts the user for a pack download URL and validates the input.
func PromptSource(in io.Reader, out io.Writer) (string, error) {
	if _, err := fmt.Fprint(out, "Enter download source for pack: "); err != nil {
		return "", fmt.Errorf("write prompt: %w", err)
	}

	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read download source: %w", err)
	}
	if errors.Is(err, io.EOF) && strings.TrimSpace(line) == "" {
		return "", fmt.Errorf("read download source: %w", err)
	}

	return ValidateSourceURL(line)
}

// WriteSource validates rawURL and writes it to SourceFileName inside dir.
func WriteSource(dir, rawURL string) error {
	validURL, err := ValidateSourceURL(rawURL)
	if err != nil {
		return err
	}

	sourcePath := filepath.Join(dir, SourceFileName)
	if err := os.WriteFile(sourcePath, []byte(validURL), fileMode); err != nil {
		return fmt.Errorf("write %s: %w", sourcePath, err)
	}

	return nil
}

// ReadSource reads and validates the download URL from SourceFileName inside dir.
func ReadSource(dir string) (string, error) {
	sourcePath := filepath.Join(dir, SourceFileName)
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", sourcePath, err)
	}

	return ValidateSourceURL(string(data))
}

// Init prompts for a pack download URL and persists it in dir.
func Init(dir string, in io.Reader, out io.Writer) error {
	sourceURL, err := PromptSource(in, out)
	if err != nil {
		return err
	}

	return WriteSource(dir, sourceURL)
}

// Download reads SourceFileName from dir, downloads the pack archive, extracts it into dir,
// and cleans up the temporary archive.
func Download(ctx context.Context, client *http.Client, dir string, out io.Writer) error {
	sourceURL, err := ReadSource(dir)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp("", "sm-pack-*.zip")
	if err != nil {
		return fmt.Errorf("create temporary archive: %w", err)
	}
	tmpPath := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close temporary archive: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := fmt.Fprintln(out, "Downloading pack..."); err != nil {
		return err
	}
	if err := downloadToFile(ctx, client, sourceURL, tmpPath); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(out, "Unpacking pack..."); err != nil {
		return err
	}
	if err := unzipArchive(tmpPath, dir, out); err != nil {
		return err
	}

	return nil
}

func downloadToFile(ctx context.Context, client *http.Client, sourceURL, dstPath string) error {
	if client == nil {
		client = defaultHTTPClient()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download pack: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download pack: unexpected HTTP status %s", resp.Status)
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create archive file: %w", err)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		_ = out.Close()
		return fmt.Errorf("write archive file: %w", err)
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("close archive file: %w", err)
	}

	return nil
}

func unzipArchive(zipPath, dstDir string, out io.Writer) error {
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer archive.Close()

	cleanDst := filepath.Clean(dstDir)
	for _, f := range archive.File {
		if err := extractZipEntry(cleanDst, f, out); err != nil {
			return err
		}
	}

	return nil
}

func extractZipEntry(dstDir string, f *zip.File, out io.Writer) error {
	targetPath := filepath.Join(dstDir, f.Name)
	rel, err := filepath.Rel(dstDir, targetPath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("illegal file path in archive: %s", f.Name)
	}
	if rel == "." {
		return nil
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(targetPath, dirMode); err != nil {
			return fmt.Errorf("create directory %s: %w", targetPath, err)
		}
		return nil
	}

	if _, err := fmt.Fprintf(out, "unzipping file %s\n", targetPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), dirMode); err != nil {
		return fmt.Errorf("create parent directory for %s: %w", targetPath, err)
	}

	srcFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("open archive entry %s: %w", f.Name, err)
	}
	defer srcFile.Close()

	mode := f.Mode().Perm()
	if mode == 0 {
		mode = fileMode
	}

	dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create extracted file %s: %w", targetPath, err)
	}

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		_ = dstFile.Close()
		return fmt.Errorf("extract file %s: %w", targetPath, err)
	}

	if err := dstFile.Close(); err != nil {
		return fmt.Errorf("close extracted file %s: %w", targetPath, err)
	}

	return nil
}
