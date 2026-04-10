package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const githubBase = "https://github.com/jgm/pandoc/releases/download"

func bootstrap(version string) error {
	dir := binDir()
	platform, err := detectPlatform()
	if err != nil {
		return err
	}

	binName := pandocBinaryName(platform)
	binPath := filepath.Join(dir, binName)

	if isExecutable(binPath) {
		fmt.Fprintf(os.Stderr, "pandoc already bootstrapped at %s\n", binPath)
		return nil
	}

	asset := assetForPlatform(platform, version)
	url := fmt.Sprintf("%s/%s/%s", githubBase, version, asset)
	inner := innerPath(platform, version)

	fmt.Fprintf(os.Stderr, "Bootstrapping pandoc %s for %s\n", version, platform)
	fmt.Fprintf(os.Stderr, "  Fetching %s ...\n", asset)

	tmpFile, err := os.CreateTemp("", "pandoc-download-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := downloadFile(url, tmpFile); err != nil {
		tmpFile.Close()
		return fmt.Errorf("download %s: %w", asset, err)
	}
	tmpFile.Close()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create bin dir: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Extracting %s ...\n", inner)

	if strings.HasSuffix(asset, ".tar.gz") {
		err = extractFromTarGz(tmpFile.Name(), inner, binPath)
	} else {
		err = extractFromZip(tmpFile.Name(), inner, binPath)
	}
	if err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	if err := os.Chmod(binPath, 0o755); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Installed -> %s\nDone.\n", binPath)
	return nil
}

func detectPlatform() (string, error) {
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "linux-amd64", nil
		case "arm64":
			return "linux-arm64", nil
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return "macos-x86_64", nil
		case "arm64":
			return "macos-arm64", nil
		}
	case "windows":
		if runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64" {
			return "windows-x86_64", nil
		}
	}
	return "", fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
}

func assetForPlatform(platform, version string) string {
	switch platform {
	case "linux-amd64":
		return fmt.Sprintf("pandoc-%s-linux-amd64.tar.gz", version)
	case "linux-arm64":
		return fmt.Sprintf("pandoc-%s-linux-arm64.tar.gz", version)
	case "macos-x86_64":
		return fmt.Sprintf("pandoc-%s-x86_64-macOS.zip", version)
	case "macos-arm64":
		return fmt.Sprintf("pandoc-%s-arm64-macOS.zip", version)
	case "windows-x86_64":
		return fmt.Sprintf("pandoc-%s-windows-x86_64.zip", version)
	}
	return ""
}

func pandocBinaryName(platform string) string {
	if strings.HasPrefix(platform, "windows") {
		return "pandoc.exe"
	}
	return "pandoc"
}

func innerPath(platform, version string) string {
	if strings.HasPrefix(platform, "windows") {
		return fmt.Sprintf("pandoc-%s/pandoc.exe", version)
	}
	return fmt.Sprintf("pandoc-%s/bin/pandoc", version)
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		return !info.IsDir()
	}
	return info.Mode()&0o111 != 0
}

func downloadFile(url string, dest *os.File) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	_, err = io.Copy(dest, resp.Body)
	return err
}

func extractFromTarGz(archivePath, innerFile, destPath string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Name == innerFile {
			out, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, tr)
			return err
		}
	}
	return fmt.Errorf("%s not found in archive", innerFile)
}

func extractFromZip(archivePath, innerFile, destPath string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == innerFile {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			out, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, rc)
			return err
		}
	}
	return fmt.Errorf("%s not found in archive", innerFile)
}
