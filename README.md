# git-rndocs

**git-rndocs** automatically generates professional release notes from Git history. Instead of manually drafting changelogs, you run a single command and get structured, categorized Markdown release notes ready for distribution.

It analyzes commits, detects versions/tags, categorizes changes using Conventional Commits, and generates structured Markdown release notes.

## Table of Contents

- [What Is git-rndocs?](#what-is-git-rndocs)
- [Use Cases](#use-cases)
- [Benefits for Developers](#benefits-for-developers)
- [Advantages Over Other Tools](#advantages-over-other-tools)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Commands](#cli-commands)
- [Configuration](#configuration)
- [Templates](#templates)
- [Conventional Commits](#conventional-commits)
- [Output Structure](#output-structure)
- [Example Output](#example-output)
- [CI/CD Integration](#cicd-integration)
- [Development](#development)
- [Architecture](#architecture)

## What Is git-rndocs?

git-rndocs is a **Git-native CLI tool** that turns your commit history into polished, human-readable release notes. It plugs directly into your Git workflow — run it as `git rndocs generate` just like any other Git subcommand or as a standalone binary.

Under the hood it:
- Reads your Git tags to detect versions (SemVer, custom prefixes, annotated or lightweight)
- Collects commits between versions
- Parses every commit against the Conventional Commits specification
- Categorizes changes into sections (Features, Bug Fixes, Performance, etc.)
- Detects breaking changes, issue references, and pull request links
- Counts contributors and generates statistics
- Outputs beautiful Markdown files or renders them through Go templates
- Optionally uploads directly to GitHub Releases

## Use Cases

| Scenario | How git-rndocs Helps |
|---|---|
| **Open-source maintainers** | Automate changelog generation for every release. Never forget to credit a contributor. |
| **CI/CD pipelines** | Run `git rndocs generate` on every tag push. Publish release notes automatically. |
| **Monorepo teams** | Generate release notes per component or per version range. |
| **Release managers** | Preview notes before publishing with `--dry-run`. Control exactly what gets included. |
| **SaaS / internal tools** | Keep stakeholders informed with professional release notes on every deploy. |
| **Compliance / auditing** | Maintain a permanent, structured record of every change per release. |

## Benefits for Developers

- **Zero manual effort** — one command replaces hours of copy-pasting commit messages
- **Consistent formatting** — every release note follows the same structure, every time
- **Less context switching** — stay in the terminal; no need to open a docs tool
- **Catch missing context** — breaking changes and issue references are surfaced automatically
- **CI-native** — runs in GitHub Actions, GitLab CI, CircleCI, Jenkins, or any shell
- **Contributor recognition** — automatically lists and sorts contributors by commit count
- **GitHub-ready** — uploads directly as a GitHub Release with `git rndocs release --upload`
- **Fully customizable** — templates, config, include/exclude filters — adapt it to your team's style

## Advantages Over Other Tools

| Aspect | git-rndocs | git-cliff | semantic-release | handwritten CHANGELOG.md |
|---|---|---|---|---|
| **Setup time** | ~10 seconds (`git rndocs init`) | Minutes (config file required) | Hours (full pipeline setup) | None, but ongoing effort |
| **Conventional Commits** | Full parser with scopes, footers, breaking changes | Basic regex matching | Limited to bump logic | N/A (manual) |
| **GitHub Releases** | Built-in (`--upload` via CLI or API) | Plugin required | Native | Manual copy-paste |
| **Templates** | Go templates (3 built-in + custom) | Tera templates | Fixed format | Any format but manual |
| **Contributor detection** | Automatic, sorted by commit count | Basic | Not included | Manual |
| **Statistics** | Files changed, insertions/deletions, category counts | Limited | Not included | Not included |
| **Monorepo support** | `--from` / `--to` ranges + custom config | Yes | Per-package setup | Manual |
| **Dry-run / preview** | `--dry-run` and `preview` subcommand | `--dry-run` | Not available | N/A |
| **Commit filtering** | Include/exclude by type, group unknown commits | Regex-based filtering | Not available | N/A |
| **Release automation** | `release` subcommand with draft/prerelease | Git hooks | Full pipeline | None |
| **Go binary** | Single static binary, no runtime deps | Rust binary | Requires Node.js | None |

If you already write Conventional Commits, git-rndocs gives you release notes for free — no extra tooling, no configuration rabbit holes, no runtime dependencies.

## Installation

```bash
go install github.com/marcuwynu23/git-rndocs@latest
```

Or download a binary from the [releases page](https://github.com/marcuwynu23/git-rndocs/releases).

## Quick Start

```bash
# Initialize in your project
git rndocs init

# Generate release notes
git rndocs generate

# Preview in terminal
git rndocs preview

# Validate setup
git rndocs validate
```

## CLI Commands

### `generate`

Generate release notes from Git history.

```bash
git rndocs generate
git rndocs generate --latest
git rndocs generate --version v2.0.0
git rndocs generate --from v1.0.0 --to HEAD
git rndocs generate --all
git rndocs generate --output ./docs/releases
git rndocs generate --template github
git rndocs generate --dry-run
git rndocs generate --overwrite
```

### `preview`

Preview release notes in the terminal without writing files.

```bash
git rndocs preview
git rndocs preview --latest
```

### `validate`

Validate repository tags, history, and configuration.

```bash
git rndocs validate
```

### `init`

Initialize git-rndocs in your project.

```bash
git rndocs init
```

Creates `.git-rndocs.yaml`, `docs/releases/`, and `templates/default.md`.

### `config`

View and manage configuration.

```bash
git rndocs config
git rndocs config --get output
```

### `release`

Generate release notes and optionally create a GitHub Release.

```bash
git rndocs release
git rndocs release --upload
git rndocs release --draft
git rndocs release --prerelease
```

## Configuration

`.git-rndocs.yaml`:

```yaml
output: docs/releases
template: default
release_name: "Version {{ .Version }}"
include:
  - feat
  - fix
  - perf
exclude:
  - chore
group_unknown: true
github:
  upload: false
contributors: true
statistics: true
```

## Templates

Built-in templates:
- `default` - Professional Markdown release notes
- `github` - GitHub-flavored release notes
- `minimal` - Minimal release notes

Custom templates use Go templates:

```markdown
# Release Notes

## Version

{{ .Version }}

{{ range .Sections }}
## {{ .Title }}

{{ range .Commits }}
- {{ .Header }}
{{ end }}
{{ end }}
```

## Conventional Commits

Recognized commit types:
- `feat` - Features
- `fix` - Bug Fixes
- `perf` - Performance
- `docs` - Documentation
- `refactor` - Refactoring
- `style` - Style
- `build` - Build
- `ci` - CI/CD
- `test` - Tests
- `chore` - Maintenance
- `revert` - Reverts

Breaking changes are detected via `!` or `BREAKING CHANGE` footer.

## Output Structure

```
docs/
└── releases/
    ├── v1.0.0/
    │   └── RELEASE-NOTES.md
    ├── v1.1.0/
    │   └── RELEASE-NOTES.md
    └── v2.0.0/
        └── RELEASE-NOTES.md
```

## Example Output

```markdown
# Release Notes

## Version

v2.0.0

## Release Date

2026-07-07

---

## Features

- **auth:** Added OAuth2 login (#123)
- Added GitHub Release upload

## Bug Fixes

- Fixed Windows path detection

## Breaking Changes

- CLI configuration renamed

## Contributors

- Alice (5 commits)
- Bob (3 commits)

## Full Changelog

v1.1.0...v2.0.0
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Release Notes
on:
  push:
    tags:
      - "v*"
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go install github.com/marcuwynu23/git-rndocs@latest
      - run: git rndocs generate
      - run: git rndocs release --upload
```

## Development

### Prerequisites

- Go 1.22+

### Build

```bash
make build
```

### Test

```bash
make test
```

### Coverage

```bash
make cover
```

### Project Structure

```
cmd/             - CLI commands
internal/
  app/           - Application orchestration
  config/        - Configuration management
  contributors/  - Contributor detection
  git/           - Git operations
  github/        - GitHub Releases integration
  markdown/      - Markdown generation
  output/        - File output
  parser/        - Conventional Commit parsing
  releasenotes/  - Release notes generation
  stats/         - Statistics collection
  template/      - Template engine
templates/       - Built-in templates
```

## Architecture

The project follows Clean Architecture principles:

- **Git access** is isolated behind interfaces
- **Business logic** is separated from CLI commands
- **Dependency injection** connects components
- **Single responsibility** per package

## License

MIT
