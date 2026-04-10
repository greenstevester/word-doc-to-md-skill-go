package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func convert(input, output string, stdout bool) error {
	if err := bootstrap(defaultPandocVersion); err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}

	pandocPath, err := findPandoc()
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp("", "pandoc-out-*.md")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

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
	dir := binDir()
	for _, name := range []string{"pandoc", "pandoc.exe"} {
		p := filepath.Join(dir, name)
		if isExecutable(p) {
			return p, nil
		}
	}
	return "", fmt.Errorf("pandoc binary not found in %s", dir)
}
