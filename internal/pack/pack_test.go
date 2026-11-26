/*
Copyright © 2023 Scott Brenner <scott@scottbrenner.me>
*/
package pack

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func buildTestZip(t *testing.T, entries map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range entries {
		if strings.HasSuffix(name, "/") {
			if _, err := zw.Create(name); err != nil {
				t.Fatalf("create zip dir %s: %v", name, err)
			}
			continue
		}
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create zip entry %s: %v", name, err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("write zip entry %s: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	return buf.Bytes()
}

func TestValidateSourceURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid https URL",
			input: "https://example.com/pack.zip",
			want:  "https://example.com/pack.zip",
		},
		{
			name:  "valid http URL with surrounding whitespace",
			input: "  http://example.com/download?id=123\r\n",
			want:  "http://example.com/download?id=123",
		},
		{
			name:    "empty string",
			input:   "   \n",
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			input:   "ftp://example.com/pack.zip",
			wantErr: true,
		},
		{
			name:    "missing host",
			input:   "https:///pack.zip",
			wantErr: true,
		},
		{
			name:    "not a URL",
			input:   "not-a-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateSourceURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateSourceURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ValidateSourceURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPromptSource(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid URL with newline",
			input: "https://example.com/pack.zip\n",
			want:  "https://example.com/pack.zip",
		},
		{
			name:  "valid URL at EOF without newline",
			input: "https://example.com/pack.zip",
			want:  "https://example.com/pack.zip",
		},
		{
			name:    "empty EOF",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			input:   "invalid\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			got, err := PromptSource(strings.NewReader(tt.input), &out)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PromptSource() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("PromptSource() = %q, want %q", got, tt.want)
			}
			if !strings.Contains(out.String(), "Enter download source for pack:") {
				t.Errorf("expected prompt in output, got %q", out.String())
			}
		})
	}
}

func TestWriteAndReadSource(t *testing.T) {
	dir := t.TempDir()

	if _, err := ReadSource(dir); err == nil {
		t.Fatal("expected error reading missing source.txt, got nil")
	}

	const rawURL = "https://example.com/simfiles.zip"
	if err := WriteSource(dir, "  "+rawURL+"\n"); err != nil {
		t.Fatalf("WriteSource() unexpected error: %v", err)
	}

	got, err := ReadSource(dir)
	if err != nil {
		t.Fatalf("ReadSource() unexpected error: %v", err)
	}
	if got != rawURL {
		t.Errorf("ReadSource() = %q, want %q", got, rawURL)
	}

	// Simulate manual edit with CRLF trailing newline.
	if err := os.WriteFile(filepath.Join(dir, SourceFileName), []byte(rawURL+"\r\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = ReadSource(dir)
	if err != nil {
		t.Fatalf("ReadSource() after CRLF edit unexpected error: %v", err)
	}
	if got != rawURL {
		t.Errorf("ReadSource() after CRLF edit = %q, want %q", got, rawURL)
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer

	err := Init(dir, strings.NewReader("https://example.com/pack.zip\n"), &out)
	if err != nil {
		t.Fatalf("Init() unexpected error: %v", err)
	}

	got, err := ReadSource(dir)
	if err != nil {
		t.Fatalf("ReadSource() unexpected error: %v", err)
	}
	if got != "https://example.com/pack.zip" {
		t.Errorf("ReadSource() = %q, want %q", got, "https://example.com/pack.zip")
	}
}

func TestDownload(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		zipData := buildTestZip(t, map[string]string{
			"MySong/":           "",
			"MySong/MySong.sm":  "#TITLE:MySong;",
			"MySong/MySong.mp3": "audio-bytes",
		})
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(zipData)
		}))
		defer srv.Close()

		dir := t.TempDir()
		if err := WriteSource(dir, srv.URL+"/pack.zip"); err != nil {
			t.Fatal(err)
		}

		var out bytes.Buffer
		if err := Download(context.Background(), srv.Client(), dir, &out); err != nil {
			t.Fatalf("Download() unexpected error: %v", err)
		}

		smBytes, err := os.ReadFile(filepath.Join(dir, "MySong", "MySong.sm"))
		if err != nil {
			t.Fatalf("read extracted file: %v", err)
		}
		if string(smBytes) != "#TITLE:MySong;" {
			t.Errorf("extracted content = %q, want %q", string(smBytes), "#TITLE:MySong;")
		}
		if !strings.Contains(out.String(), "Downloading pack...") || !strings.Contains(out.String(), "Unpacking pack...") {
			t.Errorf("unexpected output: %q", out.String())
		}
	})

	t.Run("http error status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "not found", http.StatusNotFound)
		}))
		defer srv.Close()

		dir := t.TempDir()
		if err := WriteSource(dir, srv.URL+"/missing.zip"); err != nil {
			t.Fatal(err)
		}

		var out bytes.Buffer
		err := Download(context.Background(), srv.Client(), dir, &out)
		if err == nil {
			t.Fatal("expected error for HTTP 404, got nil")
		}
		if !strings.Contains(err.Error(), "404") {
			t.Errorf("expected 404 in error, got %v", err)
		}
	})

	t.Run("missing source file", func(t *testing.T) {
		dir := t.TempDir()
		var out bytes.Buffer
		if err := Download(context.Background(), nil, dir, &out); err == nil {
			t.Fatal("expected error when source.txt is missing, got nil")
		}
	})

	t.Run("zip slip rejected", func(t *testing.T) {
		zipData := buildTestZip(t, map[string]string{
			"../escaped.txt": "bad",
		})
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(zipData)
		}))
		defer srv.Close()

		parent := t.TempDir()
		dir := filepath.Join(parent, "pack")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := WriteSource(dir, srv.URL+"/evil.zip"); err != nil {
			t.Fatal(err)
		}

		var out bytes.Buffer
		err := Download(context.Background(), srv.Client(), dir, &out)
		if err == nil {
			t.Fatal("expected zip slip error, got nil")
		}
		if _, statErr := os.Stat(filepath.Join(parent, "escaped.txt")); !os.IsNotExist(statErr) {
			t.Fatalf("zip slip wrote outside target directory")
		}
	})

	t.Run("canceled context", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(nil)
		}))
		defer srv.Close()

		dir := t.TempDir()
		if err := WriteSource(dir, srv.URL+"/pack.zip"); err != nil {
			t.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		var out bytes.Buffer
		if err := Download(ctx, srv.Client(), dir, &out); err == nil {
			t.Fatal("expected error for canceled context, got nil")
		}
	})
}
