/*
Copyright © 2023 Scott Brenner <EMAIL ADDRESS>
*/
package cmd

import (
	"io/fs"
	"reflect"
	"testing"
)

func Test_getPacks(t *testing.T) {
	tests := []struct {
		name      string
		wantPacks []fs.DirEntry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPacks := getPacks(); !reflect.DeepEqual(gotPacks, tt.wantPacks) {
				t.Errorf("getPacks() = %v, want %v", gotPacks, tt.wantPacks)
			}
		})
	}
}

func Test_promptPackSource(t *testing.T) {
	type args struct {
		packSlice []fs.DirEntry
	}
	tests := []struct {
		name        string
		args        args
		wantSources map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSources := promptPackSource(tt.args.packSlice); !reflect.DeepEqual(gotSources, tt.wantSources) {
				t.Errorf("promptPackSource() = %v, want %v", gotSources, tt.wantSources)
			}
		})
	}
}

func Test_setPackSource(t *testing.T) {
	type args struct {
		sourceMap map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setPackSource(tt.args.sourceMap)
		})
	}
}

func Test_packSource(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packSource()
		})
	}
}
