# Contributing to git-rndocs

Thank you for considering contributing to git-rndocs. This document outlines the development workflow, code standards, and processes we follow.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Development Prerequisites](#development-prerequisites)
- [Project Structure](#project-structure)
- [Makefile Reference](#makefile-reference)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Conventions](#commit-conventions)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)
- [Questions & Support](#questions--support)

---

## Code of Conduct

This project adheres to the [Contributor Covenant](https://www.contributor-covenant.org/). By participating, you are expected to uphold this code. Please report unacceptable behaviour to the project maintainers.

---

## Development Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.22+ | Compiler and toolchain |
| Make | Any | Build automation |
| golangci-lint | Latest | Static analysis |
| Git | 2.30+ | Version control |

Optional tools used by specific Makefile targets:

| Tool | Purpose | Install |
|---|---|---|
| NSIS (`makensis`) | Windows installer | `sudo apt install nsis` / `winget install NSIS.NSIS` |
| FPM | Debian package builder | `sudo gem install fpm` |

---

## Project Structure

```
├── cmd/                    # CLI command definitions (Cobra)
│   ├── root.go
│   ├── generate.go
│   ├── preview.go
│   ├── validate.go
│   ├── init.go
│   ├── config.go
│   └── release.go
├── internal/
│   ├── app/                # Application orchestrator
│   ├── config/             # Configuration (Viper)
│   ├── contributors/       # Contributor detection
│   ├── git/                # Git interface + go-git implementation
│   ├── github/             # GitHub Releases integration
│   ├── markdown/           # Markdown output builder
│   ├── output/             # File writer
│   ├── parser/             # Conventional Commit parser
│   ├── releasenotes/       # Release notes generation engine
│   ├── stats/              # Statistics collection
│   └── template/           # Go template engine
├── installers/             # Platform installer scripts (NSIS)
├── templates/              # Built-in Go templates
├── docs/
│   ├── assets/             # Images, logos
│   └── releases/           # Generated release notes
├── .github/workflows/      # CI/CD pipelines
├── main.go                 # Entry point
└── Makefile                # Build automation
```

Each package has a single responsibility. The dependency flow is:

```
cmd → app → releasenotes → { parser, markdown, template, stats, contributors }
                             ↓
                          git → { go-git }
```

---

## Makefile Reference

The Makefile is the primary interface for development tasks. It auto-detects your OS and appends `.exe` on Windows where appropriate.

### Common Commands

```bash
make build          # Compile the binary into ./build/
make test           # Run all tests with the race detector
make lint           # Run golangci-lint (if installed)
make cover          # Run tests and generate HTML coverage report
make clean          # Remove build artifacts and coverage output
make install        # Run 'go install' to $GOPATH/bin
make all            # clean → lint → test → build
```

### Symlink Commands

Use these to place the binary on your `PATH` without copying:

```bash
make link           # Symlink ./build/git-rndocs → C:/Bin/tools/git-rndocs (Windows)
                    # Symlink ./build/git-rndocs → /usr/local/bin/git-rndocs (Unix)
make unlink         # Remove the symlink
```

The link target directory defaults to `C:/Bin/tools` on Windows. Override with:

```bash
make link LINK_DIR=/custom/path
```

### Release Commands

```bash
make installer-nsis NSIS_VERSION=1.0.0    # Build Windows NSIS installer
make deb DEB_VERSION=1.0.0                # Build Debian .deb package
```

### OS Detection

The Makefile detects Windows via the `OS` environment variable and automatically:

- Appends `.exe` to the binary name
- Uses PowerShell for symlink creation
- Uses `cmd.exe` syntax for directory checks

No special shell is required — works with `cmd.exe`, PowerShell, Git Bash, and WSL.

---

## Development Workflow

### 1. Fork and Clone

```bash
git clone https://github.com/your-username/git-rndocs.git
cd git-rndocs
```

### 2. Create a Feature Branch

```bash
git checkout -b feat/my-feature
```

### 3. Make Changes

Follow the [coding standards](#coding-standards) below. Write or update tests for every change.

### 4. Run Tests

```bash
make test
```

### 5. Run Linter

```bash
make lint
```

### 6. Verify Build

```bash
make build
```

### 7. Commit

Follow the [commit conventions](#commit-conventions).

### 8. Push and Open a Pull Request

```bash
git push origin feat/my-feature
```

---

## Coding Standards

### General

- Write **idiomatic Go** — run `gofmt` and `go vet` before committing.
- Follow the [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) guide.
- Keep functions small and focused. One responsibility per function.
- Use **context** for cancellation and timeouts in all I/O and Git operations.
- Surface errors with enough context for debugging. Use `fmt.Errorf("operation: %w", err)`.
- Avoid global state. Use dependency injection through interfaces.

### Naming

- Packages: lowercase, one word, no underscores or mixed caps.
- Interfaces: `Repository`, `Writer`, `Parser` — name by behaviour.
- Exported names: use godoc-style comments on all exported symbols.
- Unexported names: be descriptive but concise.

### Imports

Group imports in three blocks separated by a blank line:

```go
import (
    // standard library
    "context"
    "fmt"

    // third-party
    "github.com/spf13/cobra"

    // internal
    "github.com/marcuwynu23/git-rndocs/internal/config"
)
```

### Error Handling

- Never silence errors. If a function returns an error, handle or propagate it.
- Use `errors.Is` / `errors.As` for sentinel error checks.
- Prefer `fmt.Errorf` with `%w` for error wrapping.

---

## Testing

- Every package must have test coverage. Aim for 80%+.
- Use `_test.go` files alongside the package they test (white-box).
- Use `t.TempDir()` for filesystem tests — it auto-cleans up.
- Use `t.Helper()` on test helpers.
- Write table-driven tests where appropriate.
- Integration tests that require a real Git repository go in `internal/git/gogit_test.go`.

Run tests:

```bash
make test        # all tests with race detector
make cover       # tests + HTML coverage report
```

Example test pattern:

```go
func TestParseCommit(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  CommitType
    }{
        {name: "feat", input: "feat: add x", want: TypeFeat},
        {name: "fix",  input: "fix: resolve y", want: TypeFix},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ParseCommit(tt.input, "", "", "")
            if got.Type != tt.want {
                t.Errorf("got %s, want %s", got.Type, tt.want)
            }
        })
    }
}
```

---

## Commit Conventions

This project generates release notes from commits, so every commit message must follow the **Conventional Commits** specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Usage |
|---|---|
| `feat` | A new feature |
| `fix` | A bug fix |
| `perf` | Performance improvement |
| `docs` | Documentation only |
| `refactor` | Code restructuring |
| `style` | Formatting, linting |
| `build` | Build system, dependencies |
| `ci` | CI/CD configuration |
| `test` | Adding or fixing tests |
| `chore` | Maintenance, tooling |
| `revert` | Reverting a previous change |

### Scope

The scope should be the package or area affected:

- `config` — configuration changes
- `parser` — commit parsing
- `git` — Git operations
- `cmd` — CLI commands
- `template` — template engine
- `markdown` — markdown generation
- `release` — release workflow
- `project` — project scaffolding, meta-files

### Examples

```
feat(parser): add support for BREAKING CHANGE footer

fix(git): handle detached HEAD gracefully

docs: add API reference to README

test(config): add tests for LoadConfig edge cases

ci: add Windows runner to CI matrix
```

Write commits that tell a story. Each commit should be a logical unit — small enough to review but large enough to be meaningful.

---

## Pull Request Process

1. **Title** must follow conventional commits (e.g., `feat(parser): add monorepo support`).
2. **Description** should explain what and why, with screenshots if UI-related.
3. **Related issues** — link any relevant issues with `Closes #123`.
4. **Checklist** in the PR body:

```markdown
- [ ] Tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Build succeeds (`make build`)
- [ ] Commits follow conventional commits
- [ ] Docs updated (if applicable)
- [ ] Changes are backward-compatible (or breaking changes documented)
```

5. A maintainer will review within 3 business days.
6. Address review feedback with additional commits — no force-pushing to shared branches.
7. Once approved, a maintainer will squash-merge into `main`.

### What Gets Merged

- Bug fixes and feature additions with tests
- Performance improvements with benchmarks
- Documentation improvements
- Refactoring that improves code quality

### What Doesn't

- Breaking changes without prior discussion (open an issue first)
- Untested code
- Changes that reduce test coverage
- Large refactors mixed with features

---

## Release Process

Releases are automated via GitHub Actions. A maintainer triggers a release by pushing a tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The `Release` workflow will:

1. Build binaries for Windows, Linux, and macOS (Intel + ARM)
2. Create a Windows NSIS installer
3. Build a Debian `.deb` package
4. Publish all artifacts to the GitHub Release

The `release-notes` workflow will also run `git rndocs generate` to produce the changelog.

---

## Questions & Support

- Open a [GitHub Discussion](https://github.com/marcuwynu23/git-rndocs/discussions) for questions.
- Open an [Issue](https://github.com/marcuwynu23/git-rndocs/issues) for bug reports and feature requests.
- For security vulnerabilities, email the maintainers directly — do not open a public issue.
