package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	insertionRe    = regexp.MustCompile(`\[([^\]]*)\]\{\.insertion\}`)
	deletionRe     = regexp.MustCompile(`\[[^\]]*\]\{\.deletion\}`)
	commentStartRe = regexp.MustCompile(`\{\.comment-start[^}]*\}`)
	commentEndRe   = regexp.MustCompile(`\[[^\]]*\]\{\.comment-end\}`)
	headingRe      = regexp.MustCompile(`^(#{1,6})\s+(.*)`)
	gridDividerRe  = regexp.MustCompile(`^\+[-=+]+\+\s*$`)
	emptyPipeRe    = regexp.MustCompile(`^\|(\s*\|)+\s*$`)
	imageWithAltRe = regexp.MustCompile(`!\[([^\]]+)\]\([^)]+\)(\{[^}]*\})?`)
	imageNoAltRe   = regexp.MustCompile(`!\[\]\([^)]+\)(\{[^}]*\})?`)
	blankLineRe    = regexp.MustCompile(`^\s*$`)
)

func postprocess(input, output string, stdout bool) error {
	data, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("read %s: %w", input, err)
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")

	text = transformTrackedChanges(text)
	text = transformHeadingHierarchy(text)
	text = transformTables(text)
	text = transformImages(text)
	text = collapseBlankLines(text)

	if stdout {
		_, err := fmt.Print(text)
		return err
	}

	if output == "" {
		output = input
	}
	return os.WriteFile(output, []byte(text), 0o644)
}

// transformTrackedChanges accepts insertions and drops deletions and comments.
func transformTrackedChanges(text string) string {
	text = insertionRe.ReplaceAllString(text, "$1")
	text = deletionRe.ReplaceAllString(text, "")
	text = commentStartRe.ReplaceAllString(text, "")
	text = commentEndRe.ReplaceAllString(text, "")
	return text
}

// transformHeadingHierarchy shifts all headings so the document starts at H1.
func transformHeadingHierarchy(text string) string {
	lines := strings.Split(text, "\n")

	minLevel := 7
	for _, line := range lines {
		if m := headingRe.FindStringSubmatch(line); m != nil {
			lvl := len(m[1])
			if lvl < minLevel {
				minLevel = lvl
			}
		}
	}

	if minLevel >= 7 || minLevel <= 1 {
		return text
	}

	shift := minLevel - 1
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if m := headingRe.FindStringSubmatch(line); m != nil {
			newLevel := len(m[1]) - shift
			if newLevel < 1 {
				newLevel = 1
			}
			line = strings.Repeat("#", newLevel) + " " + m[2]
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// transformTables removes grid-table dividers and empty pipe rows.
func transformTables(text string) string {
	lines := strings.Split(text, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if gridDividerRe.MatchString(line) || emptyPipeRe.MatchString(line) {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// transformImages replaces markdown image refs with text placeholders.
func transformImages(text string) string {
	text = imageWithAltRe.ReplaceAllString(text, "[IMAGE: $1]")
	text = imageNoAltRe.ReplaceAllString(text, "[IMAGE: no description]")
	return text
}

// collapseBlankLines reduces runs of 2+ blank lines to one.
func collapseBlankLines(text string) string {
	lines := strings.Split(text, "\n")
	result := make([]string, 0, len(lines))
	blank := 0
	for _, line := range lines {
		if blankLineRe.MatchString(line) {
			blank++
			if blank == 1 {
				result = append(result, "")
			}
			continue
		}
		blank = 0
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
