# docx-to-md

[![CI](https://github.com/greenstevester/word-doc-to-md-skill-go/actions/workflows/ci.yml/badge.svg)](https://github.com/greenstevester/word-doc-to-md-skill-go/actions/workflows/ci.yml)
[![Release](https://github.com/greenstevester/word-doc-to-md-skill-go/actions/workflows/release.yml/badge.svg)](https://github.com/greenstevester/word-doc-to-md-skill-go/actions/workflows/release.yml)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue.svg)]()

> Convert Word documents to clean, agent-readable Markdown. One binary, any platform, zero dependencies.

## What It Does

```
.docx --> pandoc --> 5 transforms --> clean .md
```

| Transform | Before | After |
|-----------|--------|-------|
| **Tracked changes** | `[new text]{.insertion}` | `new text` |
| **Heading hierarchy** | H3 -> H5 (gaps) | H1 -> H3 (no gaps) |
| **Tables** | `+----+----+` grid noise | Clean pipe tables |
| **Images** | `![](media/image1.png)` | `[IMAGE: description]` |
| **Blank lines** | 3+ consecutive blanks | Single blank line |

Pandoc auto-downloads on first run (~30 MB). No manual install needed.

## Install

### Via the Claude Code Skill (recommended)

This binary powers the [word-doc-to-md-skill](https://github.com/greenstevester/word-doc-to-md-skill) plugin. Install the skill and the binary is fetched automatically for your platform:

```
/plugin marketplace add greenstevester/docx-to-agent-md
```

### Direct download

Grab the latest release for your platform from the [releases page](https://github.com/greenstevester/word-doc-to-md-skill-go/releases/latest), or use the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/greenstevester/word-doc-to-md-skill/main/install.sh | bash
```

This detects your OS/architecture and downloads the correct binary.

### From source

```bash
git clone https://github.com/greenstevester/word-doc-to-md-skill-go.git
cd word-doc-to-md-skill-go
go build -o docx-to-md .
```

## Usage

```bash
# Convert a Word doc
./docx-to-md document.docx

# Specify output path
./docx-to-md document.docx output/clean.md

# Pipe to stdout
./docx-to-md document.docx --stdout | your-tool

# Post-process existing markdown (no pandoc needed)
./docx-to-md postprocess raw.md cleaned.md

# Force re-download pandoc
./docx-to-md bootstrap
```

## Platform Support

| OS | Architecture | Archive |
|----|-------------|---------|
| Linux | x86_64 | `docx-to-md_*_linux_amd64.tar.gz` |
| Linux | ARM64 | `docx-to-md_*_linux_arm64.tar.gz` |
| macOS | Intel | `docx-to-md_*_darwin_amd64.tar.gz` |
| macOS | Apple Silicon | `docx-to-md_*_darwin_arm64.tar.gz` |
| Windows | x86_64 | `docx-to-md_*_windows_amd64.zip` |
| Windows | ARM64 | `docx-to-md_*_windows_arm64.zip` |

## Development

```bash
make build          # Build for current platform
make build-all      # Cross-compile all 6 targets
make test           # Run tests with -race
make lint           # golangci-lint
make check          # Quick pre-commit (fmt + vet + test-short)
make ci             # Full CI pipeline locally
```

## Architecture

Single `main` package, three-stage pipeline:

1. **bootstrap.go** -- Auto-downloads pandoc from GitHub releases, caches in `./bin/`
2. **convert.go** -- Runs pandoc (`--track-changes=all -t gfm`), pipes output through postprocess
3. **postprocess.go** -- Five sequential regex/line transforms to clean the markdown

## Related

- [word-doc-to-md-skill](https://github.com/greenstevester/word-doc-to-md-skill) -- Claude Code skill plugin that uses this binary
- [fast-cc-git-hooks](https://github.com/greenstevester/fast-cc-git-hooks) -- Build/release pattern this repo follows

## License

MIT
