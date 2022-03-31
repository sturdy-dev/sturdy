package client

import (
	"context"
	"time"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/enterprise/config"
	"getsturdy.com/api/pkg/github/enterprise/db"

	gh "github.com/google/go-github/v39/github"
	"go.uber.org/zap"
)

func permissionsForInstallation(installation *github.Installation) *gh.InstallationPermissions {
	var workflows *string
	if installation.HasWorkflowsPermission {
		workflows = gh.String("write")
	}
	return &gh.InstallationPermissions{
		Contents:     gh.String("write"),
		PullRequests: gh.String("write"),
		Workflows:    workflows,
	}
}

func GetFirstAccessToken(ctx context.Context, gitHubAppConfig *config.GitHubAppConfig, installation *github.Installation, gitHubRepositoryID int64, githubClientProvider InstallationClientProvider) (*gh.InstallationToken, error) {
	// Get a new token

	_, appsClient, err := githubClientProvider(
		gitHubAppConfig,
		installation.InstallationID,
	)
	if err != nil {
		return nil, err
	}

	installToken, _, err := appsClient.CreateInstallationToken(ctx,
		installation.InstallationID,
		&gh.InstallationTokenOptions{
			RepositoryIDs: []int64{gitHubRepositoryID},
			Permissions:   permissionsForInstallation(installation),
		},
	)
	if err != nil {
		return nil, err
	}

	return installToken, nil
}

func GetAccessToken(ctx context.Context, logger *zap.Logger, gitHubAppConfig *config.GitHubAppConfig, installation *github.Installation, gitHubRepositoryID int64, repo db.GitHubRepositoryRepository, githubClientProvider InstallationClientProvider) (string, error) {
	// Check if we already have a valid token in the database
	ghr, err := repo.GetByInstallationAndGitHubRepoID(installation.InstallationID, gitHubRepositoryID)
	if err != nil {
		return "", err
	}

	logger = logger.With(zap.Int64("gitHubRepositoryID", gitHubRepositoryID))

	// Use token if it is valid for at least 10 more minutes
	if ghr.InstallationAccessToken != nil &&
		ghr.InstallationAccessTokenExpiresAt != nil &&
		ghr.InstallationAccessTokenExpiresAt.After(time.Now().Add(time.Minute*10)) {
		logger.Info("re-using existing token",
			zap.Time("expiresAt", *ghr.InstallationAccessTokenExpiresAt),
			zap.Duration("expiresIn", ghr.InstallationAccessTokenExpiresAt.Sub(time.Now())),
		)
		return *ghr.InstallationAccessToken, nil
	}

	// Get a new token
	_, appsClient, err := githubClientProvider(
		gitHubAppConfig,
		installation.InstallationID,
	)
	if err != nil {
		return "", err
	}

	installToken, _, err := appsClient.CreateInstallationToken(ctx,
		installation.InstallationID,
		&gh.InstallationTokenOptions{
			RepositoryIDs: []int64{gitHubRepositoryID},
			Permissions:   permissionsForInstallation(installation),
		},
	)
	if err != nil {
		return "", err
	}

	ghr.InstallationAccessToken = installToken.Token
	ghr.InstallationAccessTokenExpiresAt = installToken.ExpiresAt

	logger.Info("refreshed gitHub access token",
		zap.Time("expiresAt", *ghr.InstallationAccessTokenExpiresAt),
		zap.Duration("expiresIn", ghr.InstallationAccessTokenExpiresAt.Sub(time.Now())),
	)

	err = repo.Update(ghr)
	if err != nil {
		return "", err
	}

	return installToken.GetToken(), nil
}
