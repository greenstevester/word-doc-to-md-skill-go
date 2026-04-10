package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadVersionFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "reads version",
			content: "3.9.0.2",
			want:    "3.9.0.2",
		},
		{
			name:    "trims whitespace",
			content: "  3.9.0.2\n",
			want:    "3.9.0.2",
		},
		{
			name:    "empty file",
			content: "",
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, ".pandoc-version")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}
			got := readVersionFile(path)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadVersionFileMissing(t *testing.T) {
	got := readVersionFile("/nonexistent/path/.pandoc-version")
	if got != "" {
		t.Errorf("got %q for missing file, want empty string", got)
	}
}

func TestDetectPlatform(t *testing.T) {
	platform, err := detectPlatform()
	if err != nil {
		t.Fatalf("detectPlatform() failed: %v", err)
	}
	if platform == "" {
		t.Fatal("detectPlatform() returned empty string")
	}
}

func TestAssetForPlatform(t *testing.T) {
	tests := []struct {
		platform string
		version  string
		want     string
	}{
		{"linux-amd64", "3.9.0.2", "pandoc-3.9.0.2-linux-amd64.tar.gz"},
		{"linux-arm64", "3.9.0.2", "pandoc-3.9.0.2-linux-arm64.tar.gz"},
		{"macos-x86_64", "3.9.0.2", "pandoc-3.9.0.2-x86_64-macOS.zip"},
		{"macos-arm64", "3.9.0.2", "pandoc-3.9.0.2-arm64-macOS.zip"},
		{"windows-x86_64", "3.9.0.2", "pandoc-3.9.0.2-windows-x86_64.zip"},
	}
	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			got := assetForPlatform(tt.platform, tt.version)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPandocBinaryName(t *testing.T) {
	tests := []struct {
		platform string
		want     string
	}{
		{"linux-amd64", "pandoc"},
		{"macos-arm64", "pandoc"},
		{"windows-x86_64", "pandoc.exe"},
	}
	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			got := pandocBinaryName(tt.platform)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInnerPath(t *testing.T) {
	tests := []struct {
		platform string
		version  string
		want     string
	}{
		{"linux-amd64", "3.9.0.2", "pandoc-3.9.0.2/bin/pandoc"},
		{"macos-arm64", "3.9.0.2", "pandoc-3.9.0.2/bin/pandoc"},
		{"windows-x86_64", "3.9.0.2", "pandoc-3.9.0.2/pandoc.exe"},
	}
	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			got := innerPath(tt.platform, tt.version)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsExecutable(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		if isExecutable("/nonexistent/path/binary") {
			t.Error("expected false for non-existent file")
		}
	})

	t.Run("executable file", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "test-bin")
		if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
		if !isExecutable(path) {
			t.Error("expected true for executable file")
		}
	})

	t.Run("non-executable file", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "test-file")
		if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}
		if isExecutable(path) {
			t.Error("expected false for non-executable file")
		}
	})
}
