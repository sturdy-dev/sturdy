package service

import (
	"context"
	"errors"
	"fmt"

	"getsturdy.com/api/pkg/github/enterprise/config"
)

var ErrNotSetup = errors.New("not setup")

func (svc *Service) GetAppMetadata() (*config.GitHubAppMetadata, error) {
	if svc.gitHubAppConfig == nil || svc.gitHubAppConfig.ID == 0 {
		return nil, ErrNotSetup
	}

	client, err := svc.gitHubAppClientProvider(svc.gitHubAppConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create github client: %w", err)
	}

	app, _, err := client.Get(context.Background(), "")
	if err != nil {
		return nil, fmt.Errorf("could not get app from github: %w", err)
	}

	return &config.GitHubAppMetadata{
		Name: app.GetName(),
		Slug: app.GetSlug(),
	}, nil
}
