package parser

import (
	"regexp"
	"strings"
)

type CommitType string

const (
	TypeFeat     CommitType = "feat"
	TypeFix      CommitType = "fix"
	TypePerf     CommitType = "perf"
	TypeDocs     CommitType = "docs"
	TypeRefactor CommitType = "refactor"
	TypeStyle    CommitType = "style"
	TypeBuild    CommitType = "build"
	TypeCI       CommitType = "ci"
	TypeTest     CommitType = "test"
	TypeChore    CommitType = "chore"
	TypeRevert   CommitType = "revert"
	TypeUnknown  CommitType = "unknown"
)

type ParsedCommit struct {
	Raw       string
	Type      CommitType
	Scope     string
	Breaking  bool
	Header    string
	Body      string
	Footers   []Footer
	Issues    []string
	PRs       []string
	Author    string
	AuthorEmail string
	Hash      string
}

type Footer struct {
	Token string
	Value string
}

var conventionalRe = regexp.MustCompile(`^(\w+)(?:\(([^)]+)\))?(!)?:\s*(.*)`)

var breakingRe = regexp.MustCompile(`(?i)BREAKING CHANGE:\s*(.*)`)

var issueRe = regexp.MustCompile(`(?:#(\d+)|(?:fixe[sd]|close[sd]?|resolve[sd]?)\s+#(\d+))`)

var prRe = regexp.MustCompile(`\(#(\d+)\)`)

var footerRe = regexp.MustCompile(`^([\w -]+):\s*(.*)`)

func ParseCommit(msg string, hash, author, authorEmail string) *ParsedCommit {
	pc := &ParsedCommit{
		Raw:         msg,
		Hash:        hash,
		Author:      author,
		AuthorEmail: authorEmail,
		Type:        TypeUnknown,
	}

	lines := strings.Split(strings.TrimSpace(msg), "\n")
	if len(lines) == 0 {
		return pc
	}

	header := lines[0]
	pc.Header = header

	matches := conventionalRe.FindStringSubmatch(header)
	if matches != nil {
		switch CommitType(matches[1]) {
		case TypeFeat, TypeFix, TypePerf, TypeDocs, TypeRefactor, TypeStyle,
			TypeBuild, TypeCI, TypeTest, TypeChore, TypeRevert:
			pc.Type = CommitType(matches[1])
		default:
			pc.Type = TypeUnknown
		}
		pc.Scope = matches[2]
		if matches[3] == "!" {
			pc.Breaking = true
		}
		pc.Header = matches[4]
	}

	if len(lines) > 1 {
		pc.Body = strings.Join(lines[1:], "\n")
	}

	pc.Footers = parseFooters(lines)
	for _, f := range pc.Footers {
		if breakingRe.MatchString(f.Token + ": " + f.Value) {
			pc.Breaking = true
		}
	}

	pc.Issues = extractIssues(msg, pc.Issues)
	pc.PRs = extractPRs(msg, pc.PRs)

	return pc
}

func parseFooters(lines []string) []Footer {
	var footers []Footer
	inFooter := false

	for _, line := range lines[1:] {
		if matches := footerRe.FindStringSubmatch(line); matches != nil {
			footers = append(footers, Footer{
				Token: matches[1],
				Value: matches[2],
			})
			inFooter = true
		} else if inFooter && strings.HasPrefix(line, " ") {
			if len(footers) > 0 {
				footers[len(footers)-1].Value += "\n" + strings.TrimSpace(line)
			}
		} else {
			inFooter = false
		}
	}

	return footers
}

func extractIssues(msg string, existing []string) []string {
	seen := make(map[string]bool)
	for _, id := range existing {
		seen[id] = true
	}

	matches := issueRe.FindAllStringSubmatch(msg, -1)
	for _, m := range matches {
		for _, g := range m[1:] {
			if g != "" && !seen[g] {
				seen[g] = true
				existing = append(existing, g)
			}
		}
	}
	return existing
}

func extractPRs(msg string, existing []string) []string {
	seen := make(map[string]bool)
	for _, id := range existing {
		seen[id] = true
	}

	matches := prRe.FindAllStringSubmatch(msg, -1)
	for _, m := range matches {
		if !seen[m[1]] {
			seen[m[1]] = true
			existing = append(existing, m[1])
		}
	}
	return existing
}

func CategorizeCommits(commits []*ParsedCommit, include, exclude []string) map[CommitType][]*ParsedCommit {

	includeSet := make(map[CommitType]bool)
	for _, t := range include {
		includeSet[CommitType(t)] = true
	}
	excludeSet := make(map[CommitType]bool)
	for _, t := range exclude {
		excludeSet[CommitType(t)] = true
	}

	categorized := make(map[CommitType][]*ParsedCommit)

	for _, c := range commits {
		ct := c.Type
		if ct == TypeBreaking {
			ct = TypeUnknown
		}

		if len(include) > 0 && !includeSet[ct] {
			continue
		}
		if excludeSet[ct] {
			continue
		}

		categorized[ct] = append(categorized[ct], c)
	}

	return categorized
}

func (ct CommitType) DisplayName() string {
	switch ct {
	case TypeFeat:
		return "Features"
	case TypeFix:
		return "Bug Fixes"
	case TypePerf:
		return "Performance"
	case TypeDocs:
		return "Documentation"
	case TypeRefactor:
		return "Refactoring"
	case TypeStyle:
		return "Style"
	case TypeBuild:
		return "Build"
	case TypeCI:
		return "CI/CD"
	case TypeTest:
		return "Tests"
	case TypeChore:
		return "Maintenance"
	case TypeRevert:
		return "Reverts"
	case TypeBreaking:
		return "Breaking Changes"
	default:
		return "Other Changes"
	}
}

const TypeBreaking CommitType = "breaking"
