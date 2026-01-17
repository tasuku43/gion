package cli

import (
	"context"
	"fmt"
	"strings"
)

type provider interface {
	Name() string
	FetchIssues(ctx context.Context, host, owner, repoName string) ([]issueSummary, error)
	FetchIssue(ctx context.Context, host, owner, repoName string, number int) (issueSummary, error)
	FetchPRs(ctx context.Context, host, owner, repoName string) ([]prSummary, error)
	FetchPR(ctx context.Context, host, owner, repoName string, number int) (prSummary, error)
}

type githubProvider struct{}

func (githubProvider) Name() string {
	return "github"
}

func (githubProvider) FetchIssues(ctx context.Context, host, owner, repoName string) ([]issueSummary, error) {
	return fetchGitHubIssues(ctx, host, owner, repoName)
}

func (githubProvider) FetchIssue(ctx context.Context, host, owner, repoName string, number int) (issueSummary, error) {
	return fetchGitHubIssue(ctx, host, owner, repoName, number)
}

func (githubProvider) FetchPRs(ctx context.Context, host, owner, repoName string) ([]prSummary, error) {
	return fetchGitHubPRs(ctx, host, owner, repoName)
}

func (githubProvider) FetchPR(ctx context.Context, host, owner, repoName string, number int) (prSummary, error) {
	return fetchGitHubPR(ctx, host, owner, repoName, number)
}

var providers = map[string]provider{
	"github": githubProvider{},
}

func providerByName(name string) (provider, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if key == "" {
		return nil, fmt.Errorf("provider is required")
	}
	p, ok := providers[key]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", key)
	}
	return p, nil
}

func providerNameForHost(host string) string {
	lower := strings.ToLower(strings.TrimSpace(host))
	if strings.Contains(lower, "gitlab") {
		return "gitlab"
	}
	if strings.Contains(lower, "bitbucket") {
		return "bitbucket"
	}
	return "github"
}
