package client

import (
	"context"
	"net/http"

	"getsturdy.com/api/pkg/github/enterprise/config"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type InstallationClientProvider func(gitHubAppConfig *config.GitHubAppConfig, installationID int64) (tokenClient *GitHubClients, appsClient AppsClient, err error)

type PersonalClientProvider func(personalOauthToken string) (personalClient *GitHubClients, err error)

type AppClientProvider func(gitHubAppConfig *config.GitHubAppConfig) (appsClient AppsClient, err error)

type GitHubClients struct {
	Repositories RepositoriesClient
	PullRequests PullRequestsClient
	Users        UsersClient
}

type RepositoriesClient interface {
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	GetByID(ctx context.Context, id int64) (*github.Repository, *github.Response, error)
	ListCollaborators(ctx context.Context, owner, repo string, opts *github.ListCollaboratorsOptions) ([]*github.User, *github.Response, error)
}

type AppsClient interface {
	CreateInstallationToken(ctx context.Context, id int64, opts *github.InstallationTokenOptions) (*github.InstallationToken, *github.Response, error)
	GetInstallation(ctx context.Context, id int64) (*github.Installation, *github.Response, error)
	Get(ctx context.Context, appSlug string) (*github.App, *github.Response, error)
}

type PullRequestsClient interface {
	List(ctx context.Context, owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	Create(ctx context.Context, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error)
	Get(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error)
	Edit(ctx context.Context, owner string, repo string, number int, pull *github.PullRequest) (*github.PullRequest, *github.Response, error)
}

type UsersClient interface {
	Get(ctx context.Context, user string) (*github.User, *github.Response, error)
}

// NewInstallationClient creates a client for installationID that's acting on behalf of the app
func NewInstallationClient(gitHubAppConfig *config.GitHubAppConfig, installationID int64) (tokenClient *GitHubClients, appsClient AppsClient, err error) {
	jwtTransport, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, gitHubAppConfig.ID, gitHubAppConfig.PrivateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	installationTokenTransport := ghinstallation.NewFromAppsTransport(jwtTransport, installationID)

	ghClient := github.NewClient(&http.Client{Transport: installationTokenTransport})
	appsGhClient := github.NewClient(&http.Client{Transport: jwtTransport})

	return &GitHubClients{
			Repositories: ghClient.Repositories,
			PullRequests: ghClient.PullRequests,
			Users:        ghClient.Users,
		},
		appsGhClient.Apps, nil
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
		Users:        client.Users,
	}, nil
}

// NewAppClient creates a client for the app (not authenticated towards any installation)
func NewAppClient(gitHubAppConfig *config.GitHubAppConfig) (appsClient AppsClient, err error) {
	jwtTransport, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, gitHubAppConfig.ID, gitHubAppConfig.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	appsGhClient := github.NewClient(&http.Client{Transport: jwtTransport})

	return appsGhClient.Apps, nil
}
