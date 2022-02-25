package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/users"
	"golang.org/x/oauth2"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (svc *Service) AddUserIDToCodebases(ctx context.Context, userID users.ID) error {
	githubUser, err := svc.gitHubUserRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get github user: %w", err)
	}
	return svc.AddUserToCodebases(ctx, githubUser)
}

func (svc *Service) AddUserToCodebases(ctx context.Context, ghUser *github.GitHubUser) error {
	githubOAuth2Client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghUser.AccessToken}))
	githubAPIClient := gh.NewClient(githubOAuth2Client)

	installations, err := svc.listAllUserInstallations(ctx, githubAPIClient)
	if err != nil {
		return fmt.Errorf("failed to lookup installations for user: %w", err)
	}

	for _, installation := range installations {
		if err = svc.addUserToInstallationCodebases(
			ctx,
			ghUser.UserID,
			githubAPIClient,
			installation.GetID(),
		); err != nil {
			svc.logger.Error("failed to set codebase access for user", zap.Error(err))
			continue
		}
	}
	return nil
}

func (svc *Service) addUserToInstallationCodebases(ctx context.Context, userID users.ID, userAuthClient *gh.Client, installationID int64) error {
	// Truth from GitHub
	repos, err := svc.userAccessibleRepoIDs(ctx, userAuthClient, installationID)
	if err != nil {
		return fmt.Errorf("failed to get user accessible repo IDs from github: %w", err)
	}

	var repoIDs []int64
	for _, r := range repos {
		repoIDs = append(repoIDs, r.id)
	}

	gitHubRepositories, err := svc.gitHubRepositoryRepo.ListByInstallationID(installationID)
	if err != nil {
		return fmt.Errorf("failed to get user accessible repos by IDs from db: %w", err)
	}

	var hasAccess []*github.GitHubRepository
	for _, ghr := range gitHubRepositories {
		if contains(repoIDs, ghr.GitHubRepositoryID) {
			hasAccess = append(hasAccess, ghr)
		}
	}

	for _, ghr := range hasAccess {
		_, err := svc.codebaseUserRepo.GetByUserAndCodebase(userID, ghr.CodebaseID)
		if errors.Is(err, sql.ErrNoRows) {
			t0 := time.Now()
			err := svc.codebaseUserRepo.Create(codebase.CodebaseUser{
				ID:         uuid.NewString(),
				UserID:     userID,
				CodebaseID: ghr.CodebaseID,
				CreatedAt:  &t0,
			})
			if err != nil {
				svc.logger.Error("failed to create codebase-user relation", zap.Error(err))
				continue
			}

			// Send events
			svc.eventsSender.Codebase(ghr.CodebaseID, events.CodebaseUpdated, ghr.CodebaseID)
		} else if err != nil {
			svc.logger.Error("failed to get codebase-user relation", zap.Error(err))
		}
	}

	return nil
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type accessibleRepo struct {
	id   int64
	name string
}

func (svc *Service) userAccessibleRepoIDs(ctx context.Context, userAuthClient *gh.Client, installationID int64) ([]accessibleRepo, error) {
	var res []accessibleRepo

	for page := 0; page < 100; page++ {
		repos, response, err := userAuthClient.Apps.ListUserRepos(ctx, installationID, &gh.ListOptions{
			Page:    page,
			PerPage: 30,
		})
		// Track rate limiting
		svc.logger.Info("github rate limit", zap.Int64("installation_id", installationID), zap.Int("limit", response.Rate.Limit), zap.Int("remaining", response.Rate.Remaining), zap.Time("reset", response.Rate.Reset.Time))
		if err != nil {
			return nil, fmt.Errorf("failed to ListUserRepos: %w", err)
		}
		for _, r := range repos.Repositories {
			res = append(res, accessibleRepo{
				id:   r.GetID(),
				name: r.GetName(),
			})
		}
		if response.LastPage <= page || len(repos.Repositories) == 0 {
			break
		}
	}

	return res, nil
}

func (svc *Service) listAllUserInstallations(ctx context.Context, userAuthClient *gh.Client) (installation []*gh.Installation, err error) {
	var installations []*gh.Installation
	page := 1
	for page != 0 {
		newInstallations, nextPage, err := listUserInstallations(ctx, userAuthClient, page)
		if err != nil {
			return nil, err
		}
		page = nextPage
		installations = append(installations, newInstallations...)
	}
	return installations, nil
}

func listUserInstallations(ctx context.Context, userAuthClient *gh.Client, page int) (installation []*gh.Installation, nextPage int, err error) {
	installations, rsp, err := userAuthClient.Apps.ListUserInstallations(ctx, &gh.ListOptions{Page: page})
	if err != nil {
		return nil, 0, err
	}
	return installations, rsp.NextPage, nil
}
