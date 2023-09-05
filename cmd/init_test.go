/*
Copyright © 2023 Scott Brenner <scott@scottbrenner.me>
*/
package cmd

import (
	"testing"
)

func Test_getPackSource(t *testing.T) {
	tests := []struct {
		name            string
		wantDownloadURL string
		wantErr         bool
	}{
		{"Invalid URL", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDownloadURL, err := getPackSource()
			if (err != nil) != tt.wantErr {
				t.Errorf("getPackSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDownloadURL != tt.wantDownloadURL {
				t.Errorf("getPackSource() = %v, want %v", gotDownloadURL, tt.wantDownloadURL)
			}
		})
	}
}

func Test_setPackSource(t *testing.T) {
	type args struct {
		downloadURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"DanceDanceRevolution GRAND PRIX (PC)", args{"https://zenius-i-vanisher.com/v5.2/download.php?type=ddrpack&categoryid=1456"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setPackSource(tt.args.downloadURL); (err != nil) != tt.wantErr {
				t.Errorf("setPackSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initPackSource(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initPackSource(); (err != nil) != tt.wantErr {
				t.Errorf("initPackSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
