package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func convert(input, output string, stdout bool) error {
	pandocPath, err := findPandoc()
	if err != nil {
		// No pandoc found anywhere — bootstrap it
		if bErr := bootstrap(defaultPandocVersion); bErr != nil {
			return fmt.Errorf("bootstrap: %w", bErr)
		}
		pandocPath, err = findPandoc()
		if err != nil {
			return err
		}
	}

	tmpFile, err := os.CreateTemp("", "pandoc-out-*.md")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpPath := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp: %w", err)
	}
	defer func() { _ = os.Remove(tmpPath) }()

	fmt.Fprintf(os.Stderr, "Running pandoc on %s ...\n", filepath.Base(input))

	cmd := exec.Command(pandocPath, "--track-changes=all", "-t", "gfm", input, "-o", tmpPath)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pandoc: %w", err)
	}

	if err := postprocess(tmpPath, output, stdout); err != nil {
		return err
	}
	if !stdout {
		fmt.Fprintf(os.Stderr, "Written to %s\n", output)
	}
	return nil
}

func findPandoc() (string, error) {
	// 1. Check the bundled bin directory first
	dir := binDir()
	for _, name := range []string{"pandoc", "pandoc.exe"} {
		p := filepath.Join(dir, name)
		if isExecutable(p) {
			return p, nil
		}
	}

	// 2. Fall back to system PATH
	if p, err := exec.LookPath("pandoc"); err == nil {
		return p, nil
	}

	return "", fmt.Errorf("pandoc not found in %s or system PATH", dir)
}
