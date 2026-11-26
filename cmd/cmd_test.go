/*
Copyright © 2026 Scott Brenner <scott@scottbrenner.me>
*/
package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ScottBrenner/sm/internal/pack"
)

func TestInitCommand(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	var out bytes.Buffer
	rootCmd.SetIn(strings.NewReader("https://example.com/pack.zip\n"))
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"init"})
	t.Cleanup(func() {
		rootCmd.SetIn(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute() unexpected error: %v", err)
	}

	got, err := pack.ReadSource(dir)
	if err != nil {
		t.Fatalf("pack.ReadSource() unexpected error: %v", err)
	}
	if got != "https://example.com/pack.zip" {
		t.Errorf("pack.ReadSource() = %q, want %q", got, "https://example.com/pack.zip")
	}
}

func TestDownloadCommandMissingSource(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"download"})
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
	})

	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error when source.txt is missing, got nil")
	}
}
