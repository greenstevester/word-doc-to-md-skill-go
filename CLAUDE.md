# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CLI tool that converts Word documents (.docx) to clean, agent-readable Markdown. Uses pandoc (auto-downloaded) for initial conversion, then applies post-processing transforms to clean the output.

This repo builds the **platform-specific binaries** (Windows, macOS, Linux × amd64/arm64) consumed by the skill package at [greenstevester/word-doc-to-md-skill](https://github.com/greenstevester/word-doc-to-md-skill). When someone installs that skill, it detects OS/arch and selects the correct binary from this repo's releases.

Module name: `docx-to-agent-md` | Go 1.24+ (CI runs 1.26.1) | Single `main` package

## Build & Test Commands

```bash
make build                       # Build for current platform -> build/docx-to-md
make build-all                   # Cross-compile all 5 targets -> build/
make test                        # Run all tests with -race
make lint                        # golangci-lint
make coverage                    # Generate coverage report
make check                       # Quick pre-commit: fmt + vet + test-short
make ci                          # Full CI pipeline locally
go test -run TestTransformImages # Run a single test by name
```

## Architecture

Three-stage pipeline, all in package `main`:

1. **bootstrap.go** — Downloads pandoc from GitHub releases on first run (version-pinned at `defaultPandocVersion` in main.go). Detects OS/arch automatically. Binary lands in `./bin/pandoc`.
2. **convert.go** — Runs pandoc with `--track-changes=all -t gfm` to produce raw GFM markdown in a temp file, then pipes through postprocess.
3. **postprocess.go** — Regex-based transforms applied in order: accept tracked insertions / drop deletions & comments → normalize heading hierarchy (shift to start at H1) → strip grid-table dividers & empty pipe rows → replace image refs with `[IMAGE: ...]` placeholders → collapse consecutive blank lines.

`main.go` dispatches subcommands: default (full convert), `postprocess` (postprocess-only), `bootstrap` (download pandoc only). Output goes to file by default or stdout with `--stdout`.

## Build & Release

Pattern copied from [greenstevester/fast-cc-git-hooks](https://github.com/greenstevester/fast-cc-git-hooks):

- **Makefile** — local builds, cross-compile, lint, test, coverage
- **`.goreleaser.yml`** — `CGO_ENABLED=0`, `-trimpath`, archives as `docx-to-md_{version}_{os}_{arch}.tar.gz` (`.zip` for Windows), checksums.txt
- **`.github/workflows/ci.yml`** — matrix test (ubuntu/macos/windows), lint, security scan, cross-platform build
- **`.github/workflows/release.yml`** — push to main triggers auto-semver tag (BREAKING→major, feat→minor, fix→patch) + GoReleaser

The skill repo's `install.sh` downloads from these releases based on the user's OS/arch.

## Key Design Decisions

- Pandoc is **not** a system dependency — it is self-bootstrapped into `./bin/` and should not be committed.
- Images are intentionally replaced with text placeholders (the output targets AI agents, not humans).
- Tracked changes are resolved (insertions accepted, deletions removed) rather than preserved.
