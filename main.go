package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultPandocVersion = "3.9.0.2"

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage()
		os.Exit(0)
	}

	var err error
	switch args[0] {
	case "bootstrap":
		version := defaultPandocVersion
		if len(args) > 1 {
			version = args[1]
		}
		err = bootstrap(version)
	case "postprocess":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "ERROR: postprocess requires an input file")
			os.Exit(1)
		}
		stdout, output := parseOutputArgs(args[2:], args[1])
		err = postprocess(args[1], output, stdout)
	default:
		input := args[0]
		if _, statErr := os.Stat(input); statErr != nil {
			fmt.Fprintf(os.Stderr, "ERROR: file not found: %s\n", input)
			os.Exit(1)
		}
		stdout, output := parseOutputArgs(args[1:], input)
		err = convert(input, output, stdout)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func parseOutputArgs(args []string, input string) (bool, string) {
	var stdout bool
	var output string
	for _, a := range args {
		if a == "--stdout" {
			stdout = true
		} else {
			output = a
		}
	}
	if output == "" && !stdout {
		ext := filepath.Ext(input)
		output = strings.TrimSuffix(input, ext) + ".md"
	}
	return stdout, output
}

// binDir returns the directory where pandoc is stored.
// It lives next to the docx-to-md executable (i.e. inside the skill plugin dir),
// NOT in the user's current working directory.
func binDir() string {
	exe, err := os.Executable()
	if err != nil {
		// Fall back to CWD if we can't resolve the executable path.
		wd, _ := os.Getwd()
		return filepath.Join(wd, "bin")
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		wd, _ := os.Getwd()
		return filepath.Join(wd, "bin")
	}
	return filepath.Join(filepath.Dir(exe), "bin")
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage:
  docx-to-md <input.docx> [output.md] [--stdout]
  docx-to-md postprocess <input.md> [output.md] [--stdout]
  docx-to-md bootstrap [version]

Convert Word documents to clean, agent-readable Markdown.
Pandoc is auto-downloaded on first run (default version: %s).
`, defaultPandocVersion)
}
