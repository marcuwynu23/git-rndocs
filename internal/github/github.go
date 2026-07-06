package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ReleaseOptions struct {
	Repo        string
	Tag         string
	Name        string
	Body        string
	Draft       bool
	Prerelease  bool
	Token       string
}

type GitHubRelease struct {
	ID          int64  `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	Draft       bool   `json:"draft"`
	Prerelease  bool   `json:"prerelease"`
	HTMLURL     string `json:"html_url"`
}

func CreateRelease(ctx context.Context, opts *ReleaseOptions) (*GitHubRelease, error) {
	if hasGHCli() {
		return createViaGHCli(ctx, opts)
	}
	return createViaAPI(ctx, opts)
}

func hasGHCli() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

func createViaGHCli(ctx context.Context, opts *ReleaseOptions) (*GitHubRelease, error) {
	args := []string{"release", "create", opts.Tag}
	args = append(args, "--title", opts.Name)
	args = append(args, "--notes", opts.Body)

	if opts.Draft {
		args = append(args, "--draft")
	}
	if opts.Prerelease {
		args = append(args, "--prerelease")
	}

	if opts.Repo != "" {
		args = append(args, "--repo", opts.Repo)
	}

	cmd := exec.CommandContext(ctx, "gh", args...)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh CLI failed: %w", err)
	}

	url := strings.TrimSpace(string(output))
	return &GitHubRelease{
		TagName:    opts.Tag,
		Name:       opts.Name,
		Body:       opts.Body,
		Draft:      opts.Draft,
		Prerelease: opts.Prerelease,
		HTMLURL:    url,
	}, nil
}

func createViaAPI(ctx context.Context, opts *ReleaseOptions) (*GitHubRelease, error) {
	token := opts.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is required for API release creation")
	}

	repo := opts.Repo
	if repo == "" {
		repo = os.Getenv("GITHUB_REPOSITORY")
	}
	if repo == "" {
		return nil, fmt.Errorf("repository not specified and GITHUB_REPOSITORY not set")
	}

	payload := map[string]interface{}{
		"tag_name":   opts.Tag,
		"name":       opts.Name,
		"body":       opts.Body,
		"draft":      opts.Draft,
		"prerelease": opts.Prerelease,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var release GitHubRelease
	if err := json.Unmarshal(respBody, &release); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &release, nil
}
