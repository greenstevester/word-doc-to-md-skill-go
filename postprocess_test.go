package main

import (
	"testing"
)

func TestTransformTrackedChanges(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "accept insertions",
			input: "Hello [world]{.insertion} foo",
			want:  "Hello world foo",
		},
		{
			name:  "drop deletions",
			input: "Hello [old text]{.deletion} world",
			want:  "Hello  world",
		},
		{
			name:  "drop comment start markers",
			input: "text {.comment-start id=\"1\" author=\"Bob\"} more",
			want:  "text  more",
		},
		{
			name:  "drop comment end markers",
			input: "text [my comment]{.comment-end} more",
			want:  "text  more",
		},
		{
			name:  "mixed insertions and deletions",
			input: "[new]{.insertion} replaces [old]{.deletion}",
			want:  "new replaces ",
		},
		{
			name:  "no tracked changes",
			input: "plain text with no changes",
			want:  "plain text with no changes",
		},
		{
			name:  "empty insertion",
			input: "before []{.insertion} after",
			want:  "before  after",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformTrackedChanges(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTransformHeadingHierarchy(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "shift H3 and H5 to H1 and H3",
			input: "### Title\n\nSome text\n\n##### Sub-section",
			want:  "# Title\n\nSome text\n\n### Sub-section",
		},
		{
			name:  "already starts at H1",
			input: "# Title\n\n## Section\n\n### Sub",
			want:  "# Title\n\n## Section\n\n### Sub",
		},
		{
			name:  "no headings",
			input: "Just some plain text\nwith no headings",
			want:  "Just some plain text\nwith no headings",
		},
		{
			name:  "single H2 becomes H1",
			input: "## Only heading",
			want:  "# Only heading",
		},
		{
			name:  "preserves non-heading lines",
			input: "## Title\n\nParagraph here\n\n### Sub\n\nMore text",
			want:  "# Title\n\nParagraph here\n\n## Sub\n\nMore text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformHeadingHierarchy(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTransformTables(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "remove grid dividers",
			input: "+----+----+\n| A  | B  |\n+----+----+\n| 1  | 2  |\n+----+----+",
			want:  "| A  | B  |\n| 1  | 2  |",
		},
		{
			name:  "remove empty pipe rows",
			input: "| A | B |\n|   |   |\n| 1 | 2 |",
			want:  "| A | B |\n| 1 | 2 |",
		},
		{
			name:  "preserve content rows",
			input: "| Name | Value |\n| foo  | bar   |",
			want:  "| Name | Value |\n| foo  | bar   |",
		},
		{
			name:  "no tables",
			input: "Just text\nwith no tables",
			want:  "Just text\nwith no tables",
		},
		{
			name:  "grid divider with equals",
			input: "+====+====+\n| H1 | H2 |\n+----+----+",
			want:  "| H1 | H2 |",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformTables(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTransformImages(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "image with alt text",
			input: "See ![diagram](images/arch.png) below",
			want:  "See [IMAGE: diagram] below",
		},
		{
			name:  "image without alt text",
			input: "See ![](images/photo.jpg) below",
			want:  "See [IMAGE: no description] below",
		},
		{
			name:  "image with width attribute",
			input: "![logo](img/logo.png){width=200}",
			want:  "[IMAGE: logo]",
		},
		{
			name:  "multiple images",
			input: "![a](1.png) and ![b](2.png)",
			want:  "[IMAGE: a] and [IMAGE: b]",
		},
		{
			name:  "no images",
			input: "Just a [link](http://example.com)",
			want:  "Just a [link](http://example.com)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformImages(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCollapseBlankLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "collapse multiple blank lines",
			input: "line 1\n\n\n\nline 2",
			want:  "line 1\n\nline 2",
		},
		{
			name:  "preserve single blank line",
			input: "line 1\n\nline 2",
			want:  "line 1\n\nline 2",
		},
		{
			name:  "no blank lines",
			input: "line 1\nline 2\nline 3",
			want:  "line 1\nline 2\nline 3",
		},
		{
			name:  "whitespace-only lines count as blank",
			input: "line 1\n   \n  \n\nline 2",
			want:  "line 1\n\nline 2",
		},
		{
			name:  "multiple groups of blanks",
			input: "a\n\n\nb\n\n\n\nc",
			want:  "a\n\nb\n\nc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collapseBlankLines(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
