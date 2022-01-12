package client

import (
	"context"
	"net/http"

	"mash/pkg/github/config"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type ClientProvider func(gitHubAppConfig config.GitHubAppConfig, installationID int64) (tokenClient *GitHubClients, jwtClient *GitHubClients, err error)
type PersonalClientProvider func(personalOauthToken string) (personalClient *GitHubClients, err error)

type GitHubClients struct {
	Repositories RepositoriesClient
	PullRequests PullRequestsClient
	Apps         AppsClient
}

type RepositoriesClient interface {
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	GetByID(ctx context.Context, id int64) (*github.Repository, *github.Response, error)
	ListCollaborators(ctx context.Context, owner, repo string, opts *github.ListCollaboratorsOptions) ([]*github.User, *github.Response, error)
}

type AppsClient interface {
	CreateInstallationToken(ctx context.Context, id int64, opts *github.InstallationTokenOptions) (*github.InstallationToken, *github.Response, error)
}

type PullRequestsClient interface {
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	Create(ctx context.Context, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error)
	Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error)
	Edit(ctx context.Context, owner string, repo string, number int, pull *github.PullRequest) (*github.PullRequest, *github.Response, error)
}

func NewClient(gitHubAppConfig config.GitHubAppConfig, installationID int64) (tokenClient *GitHubClients, jwtClient *GitHubClients, err error) {
	jwtTransport, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, gitHubAppConfig.GitHubAppID, gitHubAppConfig.GitHubAppPrivateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	installationTokenTransport := ghinstallation.NewFromAppsTransport(jwtTransport, installationID)

	ghClient := github.NewClient(&http.Client{Transport: installationTokenTransport})
	appsGhClient := github.NewClient(&http.Client{Transport: jwtTransport})

	return &GitHubClients{
			Repositories: ghClient.Repositories,
			PullRequests: ghClient.PullRequests,
		},
		&GitHubClients{
			Apps: appsGhClient.Apps,
		}, nil
}

// NewPersonalClient returns a client that users the users GitHub OAuth token to act on their behalf
// On GitHub these actions shows up as "$USERNAME via Sturdy"
func NewPersonalClient(personalOauthToken string) (personalClient *GitHubClients, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: personalOauthToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubClients{
		Repositories: client.Repositories,
		PullRequests: client.PullRequests,
	}, nil
}
