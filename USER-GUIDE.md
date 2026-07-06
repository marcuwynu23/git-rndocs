# git-rndocs User Guide

A complete walkthrough of installing, configuring, and using git-rndocs.

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Commands Reference](#cli-commands-reference)
- [Configuration](#configuration)
- [Templates](#templates)
- [Conventional Commits](#conventional-commits)
- [CI/CD Integration](#cicd-integration)
- [Monorepo Workflows](#monorepo-workflows)
- [GitHub Releases](#github-releases)
- [Output Customisation](#output-customisation)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Installation

### From Source

```bash
go install github.com/marcuwynu23/git-rndocs@latest
```

### From Binary

Download the latest release for your platform from the [releases page](https://github.com/marcuwynu23/git-rndocs/releases), extract it, and place the binary somewhere on your `PATH`.

```bash
# Linux / macOS
chmod +x git-rndocs
sudo mv git-rndocs /usr/local/bin/

# Windows
# Move git-rndocs.exe to a directory in your PATH, e.g. C:\Bin\tools
```

### Verify

```bash
git rndocs --help
```

You should see the list of available commands.

---

## Quick Start

### 1. Initialise in Your Project

```bash
cd my-project
git rndocs init
```

This creates:

- `.git-rndocs.yaml` — configuration file
- `docs/releases/` — output directory
- `templates/default.md` — customisable template

### 2. Generate Release Notes

```bash
git rndocs generate
```

This scans tags, collects commits since the last tag, categorises them, and writes release notes to `docs/releases/<version>/RELEASE-NOTES.md`.

### 3. Preview Before Writing

```bash
git rndocs preview
```

Shows the generated release notes in the terminal without creating any files.

### 4. Validate Your Setup

```bash
git rndocs validate
```

Checks that tags, history, configuration, and templates are all in good shape.

---

## CLI Commands Reference

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

| Flag | Default | Description |
|---|---|---|
| `--repo` | `.` | Path to the Git repository |
| `--output` | `docs/releases` | Output directory |
| `--version` | — | Specific version/tag to generate notes for |
| `--from` | — | Starting commit or tag |
| `--to` | `HEAD` | Ending commit or tag |
| `--all` | `false` | Generate notes for every version (tag) |
| `--latest` | `false` | Generate notes for the latest version only |
| `--template` | `default` | Template name (`default`, `github`, `minimal`, or custom) |
| `--config` | — | Path to a config file |
| `--json` | `false` | Output as JSON |
| `--overwrite` | `false` | Overwrite existing release notes |
| `--dry-run` | `false` | Show what would be generated without writing |
| `-v`, `--verbose` | `false` | Verbose output |

#### Examples by Use Case

**Latest release only** (typical CI use):

```bash
git rndocs generate --latest
```

**All versions, full history:**

```bash
git rndocs generate --all
```

**Custom version range:**

```bash
git rndocs generate --from v1.0.0 --to v2.0.0
```

**Specific tag, up to HEAD:**

```bash
git rndocs generate --from v1.0.0
```

**Different output location:**

```bash
git rndocs generate --output ./CHANGELOG
```

---

### `preview`

Show release notes in the terminal without writing files.

```bash
git rndocs preview
git rndocs preview --latest
git rndocs preview --from v1.0.0 --to HEAD
```

Useful for checking what will be generated before committing to a release.

---

### `validate`

Check that everything is configured correctly.

```bash
git rndocs validate
```

Validates:

- Tags exist and are parseable
- Git history is accessible
- Configuration file is valid
- Templates are loadable
- Repository is healthy

---

### `init`

Scaffold git-rndocs in a project.

```bash
git rndocs init
```

Creates:

```
.git-rndocs.yaml
docs/releases/
templates/default.md
```

---

### `config`

View or query the current configuration.

```bash
git rndocs config
git rndocs config --get output
git rndocs config --get template
```

---

### `release`

Generate release notes and optionally publish them as a GitHub Release.

```bash
git rndocs release
git rndocs release --upload
git rndocs release --draft
git rndocs release --prerelease
git rndocs release --upload --draft
```

| Flag | Default | Description |
|---|---|---|
| `--upload` | `false` | Upload to GitHub Releases |
| `--draft` | `false` | Create as a draft release |
| `--prerelease` | `false` | Mark as a pre-release |

When `--upload` is set, the tool first tries to use the `gh` CLI. If that's not available, it falls back to the GitHub REST API (requires `GITHUB_TOKEN`).

---

## Configuration

The configuration file `.git-rndocs.yaml` controls all defaults.

### Reference

```yaml
# Output directory for generated release notes
output: docs/releases

# Template name (default, github, minimal, or custom)
template: default

# Release title template (Go template syntax)
release_name: "Version {{ .Version }}"

# Commit types to include (empty = all)
include:
  - feat
  - fix
  - perf
  - docs
  - refactor
  - style
  - build
  - ci
  - test
  - chore
  - revert

# Commit types to exclude (takes precedence over include)
exclude: []

# Group non-conventional commits under "Other Changes"
group_unknown: true

# GitHub integration settings
github:
  upload: false

# Automatically detect and list contributors
contributors: true

# Generate statistics (files changed, insertions, deletions)
statistics: true
```

### Configuration Precedence

1. CLI flags (highest priority)
2. `.git-rndocs.yaml` in the current directory
3. Custom config path via `--config`
4. Built-in defaults

---

## Templates

### Built-in Templates

| Name | Description |
|---|---|
| `default` | Professional release notes with full details, sections, contributors, and changelog link |
| `github` | Compact format optimised for GitHub Releases |
| `minimal` | Bare-bones — version, date, and bullet points only |

### Using a Built-in Template

```bash
git rndocs generate --template github
```

### Custom Templates

Templates use Go's `text/template` syntax. Place custom `.md` files in the `templates/` directory.

**Example custom template** (`templates/custom.md`):

```markdown
# Release {{ .Version }}

Released on {{ .ReleaseDate }}

{{ range .Sections }}
## {{ .Title }}

{{ range .Commits }}
- {{ .Header }}{{ if .Breaking }} ⚠️{{ end }}
{{ end }}
{{ end }}

{{ if .Contributors }}
**Contributors:** {{ range .Contributors }}{{ .Name }}, {{ end }}
{{ end }}
```

Use it with:

```bash
git rndocs generate --template custom
```

### Template Data Model

| Field | Type | Description |
|---|---|---|
| `.Version` | `string` | Version or tag name |
| `.ReleaseDate` | `string` | Date of generation |
| `.Summary` | `string` | Optional summary text |
| `.Sections` | `[]Section` | Categorised commit groups |
| `.Breaking` | `[]ParsedCommit` | All commits with breaking changes |
| `.Contributors` | `[]Contributor` | Contributors sorted by commit count |
| `.Statistics` | `*Statistics` | Code change statistics |
| `.From` | `string` | Start of range |
| `.To` | `string` | End of range |

Each `Section` has:

| Field | Type | Description |
|---|---|---|
| `.Title` | `string` | Section heading (e.g. "Features") |
| `.Commits` | `[]ParsedCommit` | Commits in this category |

Each `ParsedCommit` has:

| Field | Type | Description |
|---|---|---|
| `.Header` | `string` | Commit description |
| `.Type` | `string` | Commit type (`feat`, `fix`, etc.) |
| `.Scope` | `string` | Optional scope |
| `.Breaking` | `bool` | Whether this is a breaking change |
| `.Body` | `string` | Full commit body |
| `.Issues` | `[]string` | Referenced issue numbers |
| `.PRs` | `[]string` | Referenced PR numbers |

---

## Conventional Commits

git-rndocs parses commits following the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Recognised Types

| Type | Section |
|---|---|
| `feat` | Features |
| `fix` | Bug Fixes |
| `perf` | Performance |
| `docs` | Documentation |
| `refactor` | Refactoring |
| `style` | Style |
| `build` | Build |
| `ci` | CI/CD |
| `test` | Tests |
| `chore` | Maintenance |
| `revert` | Reverts |
| *(non-conventional)* | Other Changes |

### Breaking Changes

Detected in two ways:

1. **`!` before the colon**: `feat!: remove deprecated API`
2. **`BREAKING CHANGE` footer**:

```
feat: add new API

BREAKING CHANGE: old API removed
```

### Issue and PR References

- `#123` in commit messages are extracted as issue references
- `(#123)` at the end of the header is extracted as a PR reference
- Both are included in the generated Markdown output

---

## CI/CD Integration

### GitHub Actions

Create `.github/workflows/release-notes.yml`:

```yaml
name: Generate Release Notes

on:
  push:
    tags:
      - "v*"

jobs:
  release-notes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - run: go install github.com/marcuwynu23/git-rndocs@latest

      - name: Generate release notes
        run: git rndocs generate --overwrite

      - name: Upload to GitHub Release
        run: git rndocs release --upload
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Key points for CI:

- Use `fetch-depth: 0` so git-rndocs can see all tags and history
- Use `--overwrite` if the output directory already exists from a previous run
- Set `GITHUB_TOKEN` for API-based release uploads

### Other CI Systems

git-rndocs is a single static binary — it runs anywhere Go compiles:

```yaml
# GitLab CI
before_script:
  - go install github.com/marcuwynu23/git-rndocs@latest

script:
  - git rndocs generate --overwrite
```

```yaml
# CircleCI
steps:
  - run: go install github.com/marcuwynu23/git-rndocs@latest
  - run: git rndocs generate --overwrite
```

---

## Monorepo Workflows

### Generate for a Subdirectory

Use `--from` and `--to` to target specific version ranges:

```bash
git rndocs generate --from v1.0.0 --to v2.0.0
```

### Multiple Configs

Create separate config files for each component:

```bash
git rndocs generate --config .git-rndocs-frontend.yaml
git rndocs generate --config .git-rndocs-backend.yaml
```

### Named Ranges

Use any reference that `git` understands:

```bash
git rndocs generate --from v1.0.0 --to v1.1.0
git rndocs generate --from v1.1.0 --to HEAD
git rndocs generate --from <commit-hash>
```

---

## GitHub Releases

Automatic publishing to GitHub Releases is handled by the `release` command.

### Prerequisites

**Option A: GitHub CLI (`gh`)**

```bash
gh auth login
git rndocs release --upload
```

**Option B: API Token**

```bash
export GITHUB_TOKEN=ghp_xxxxx
git rndocs release --upload
```

If both are available, the `gh` CLI is preferred.

### Workflow

```bash
# Preview first
git rndocs release --dry-run

# Create and upload
git rndocs release --upload

# Draft, upload later
git rndocs release --upload --draft

# Pre-release
git rndocs release --upload --prerelease
```

---

## Output Customisation

### Output Structure

By default:

```
docs/
└── releases/
    └── v2.0.0/
        └── RELEASE-NOTES.md
```

Custom output location:

```bash
git rndocs generate --output ./changelogs
```

```
changelogs/
└── v2.0.0/
    └── RELEASE-NOTES.md
```

### Including / Excluding Commit Types

```yaml
# Only features and fixes
include:
  - feat
  - fix

# But exclude chore commits
exclude:
  - chore
```

### Statistics and Contributors

These are enabled by default. Disable them in `.git-rndocs.yaml`:

```yaml
contributors: false
statistics: false
```

---

## Troubleshooting

### "No release notes generated"

**Cause:** No tags found in the repository.

**Fix:** Create at least one tag:

```bash
git tag v1.0.0
git rndocs generate
```

### "Unknown commits in output"

**Cause:** Commits that don't follow Conventional Commits appear under "Other Changes".

**Fix:** Update commit messages to follow the `type: description` format.

### "Output directory already exists"

**Cause:** The output directory from a previous run still exists.

**Fix:** Use `--overwrite` to replace existing files.

### "Cannot create symlink — administrator privilege required" (Windows)

**Cause:** Windows requires admin or Developer Mode for symlinks.

**Fix:** Run the terminal as Administrator, or enable Developer Mode (Settings → For Developers → Developer Mode).

### "GITHUB_TOKEN not set"

**Cause:** No GitHub token is available for API-based release uploads.

**Fix:** Set the `GITHUB_TOKEN` environment variable, or install and authenticate the `gh` CLI.

### "Template not found"

**Cause:** The specified template doesn't exist.

**Fix:** Use one of the built-in names (`default`, `github`, `minimal`) or create a `.md` file in the `templates/` directory.

### `go install` fails

**Cause:** Missing or outdated Go toolchain.

**Fix:** Ensure Go 1.22+ is installed:

```bash
go version
```

---

## FAQ

**Q: Does git-rndocs modify my Git history?**

No. It only reads your repository. It never writes to Git objects, refs, or the working tree (other than the output directory).

**Q: Can I use it without tags?**

Yes. Running `git rndocs generate` without tags will generate notes for the entire commit history.

**Q: Does it support Git submodules?**

Yes. Point it at the submodule path: `git rndocs generate --repo path/to/submodule`.

**Q: Can I use my own template engine?**

The built-in Go template engine supports conditionals, ranges, and functions. For most use cases, custom templates are sufficient. If you need something else, the Markdown output can be piped to other tools.

**Q: Does it work on Windows?**

Yes. git-rndocs is fully cross-platform. Symlink features require Administrator or Developer Mode on Windows.

**Q: What if I use emoji commits like `✨ feat: add login`?**

The parser looks for the conventional `type:` prefix. Emoji before it (`✨ feat:`) is ignored and the commit is parsed normally. Emoji within the description is preserved in the output.

**Q: Can I generate release notes for a private repository?**

Yes. The tool only needs read access to the Git repository. GitHub Release uploads require a valid token with `repo` scope.

**Q: How do I contribute?**

See [CONTRIBUTING.md](CONTRIBUTING.md) for development workflow, coding standards, and pull request process.

---

*For CLI reference, see `git rndocs --help` or `git rndocs <command> --help`.*
