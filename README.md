# git-rndocs

**git-rndocs** automatically generates professional release notes from Git history.

It analyzes commits, detects versions/tags, categorizes changes using Conventional Commits, and generates structured Markdown release notes.

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
