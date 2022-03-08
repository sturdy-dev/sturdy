package enterprise

import (
	"context"
	"getsturdy.com/api/pkg/github/enterprise/config"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
	"sync"
)

type gitHubAppRootResolver struct {
	conf    *config.GitHubAppConfig
	meta    *config.GitHubAppMetadata
	service *service_github.Service
}

func NewGitHubAppRootResolver(conf *config.GitHubAppConfig, meta *config.GitHubAppMetadata, service *service_github.Service) resolvers.GitHubAppRootResolver {
	return &gitHubAppRootResolver{
		conf:    conf,
		meta:    meta,
		service: service,
	}
}

func (r *gitHubAppRootResolver) GitHubApp() resolvers.GitHubApp {
	return &gitHubAppResolver{root: r}
}

type gitHubAppResolver struct {
	root *gitHubAppRootResolver
}

func (r *gitHubAppResolver) ID() graphql.ID {
	return "sturdy"
}

func (r *gitHubAppResolver) Name() string {
	return r.root.meta.Slug
}

func (r *gitHubAppResolver) ClientID() string {
	return r.root.conf.ClientID
}

func (r *gitHubAppResolver) Validation() resolvers.GithubValidationApp {
	return &GithubAppValidationResolver{service: r.root.service}
}

type GithubAppValidationResolver struct {
	doOnce             sync.Once
	service            *service_github.Service
	valid              bool
	missingPermissions []string
	missingEvents      []string
	err                error
}

func (s *GithubAppValidationResolver) fetch(ctx context.Context) {
	valid, permissions, events, err := s.service.CheckPermissions(ctx)
	s.valid = valid
	s.missingPermissions = permissions
	s.missingEvents = events
	s.err = err
}

func (s *GithubAppValidationResolver) Ok(ctx context.Context) (bool, error) {
	s.doOnce.Do(func() {
		s.fetch(ctx)
	})
	return s.valid, s.err
}

func (s *GithubAppValidationResolver) MissingPermissions(ctx context.Context) ([]string, error) {
	s.doOnce.Do(func() {
		s.fetch(ctx)
	})
	return s.missingPermissions, s.err
}

func (s *GithubAppValidationResolver) MissingEvents(ctx context.Context) ([]string, error) {
	s.doOnce.Do(func() {
		s.fetch(ctx)
	})
	return s.missingEvents, s.err
}
