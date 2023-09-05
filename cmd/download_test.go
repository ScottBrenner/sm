/*
Copyright © 2023 Scott Brenner <scott@scottbrenner.me>
*/
package cmd

import "testing"

func Test_openSourceFile(t *testing.T) {
	tests := []struct {
		name          string
		wantSourceURL string
		wantErr       bool
	}{
		{"No such file", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSourceURL, err := openSourceFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("openSourceFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSourceURL != tt.wantSourceURL {
				t.Errorf("openSourceFile() = %v, want %v", gotSourceURL, tt.wantSourceURL)
			}
		})
	}
}

func Test_downloadFromURL(t *testing.T) {
	type args struct {
		sourceURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Valid", args{sourceURL: "https://zenius-i-vanisher.com/v5.2/download.php?type=ddrsimfilecustom&simfileid=48669"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := downloadFromURL(tt.args.sourceURL); (err != nil) != tt.wantErr {
				t.Errorf("downloadFromURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_unzipDownloadedPack(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"File does not exist", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := unzipDownloadedPack(); (err != nil) != tt.wantErr {
				t.Errorf("unzipDownloadedPack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_removeZip(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"No such file", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removeZip(); (err != nil) != tt.wantErr {
				t.Errorf("removeZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_downloadPack(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"No such file", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := downloadPack(); (err != nil) != tt.wantErr {
				t.Errorf("downloadPack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
