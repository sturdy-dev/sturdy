package service

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	service_buildkite "getsturdy.com/api/pkg/buildkite/service"
	"getsturdy.com/api/pkg/changes"
	"getsturdy.com/api/pkg/ci"
	db_ci "getsturdy.com/api/pkg/ci/db"
	"getsturdy.com/api/pkg/ci/service/configuration"
	"getsturdy.com/api/pkg/codebases"
	service_github "getsturdy.com/api/pkg/github/service"
	"getsturdy.com/api/pkg/integrations"
	db_integrations "getsturdy.com/api/pkg/integrations/db"
	"getsturdy.com/api/pkg/integrations/providers"
	"getsturdy.com/api/pkg/jwt"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/snapshots"
	service_snaphsotter "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/statuses"
	svc_statuses "getsturdy.com/api/pkg/statuses/service"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	oneDay = 24 * time.Hour
)

// TODO: refactor to have a more generic trigger method
type Service struct {
	logger           *zap.Logger
	executorProvider executor.Provider

	configRepo   db_integrations.IntegrationsRepository
	ciCommitRepo db_ci.CommitRepository

	buildkiteService service_buildkite.Service
	githubService    service_github.Service

	publicApiHostname string
	statusService     *svc_statuses.Service
	jwtService        *service_jwt.Service
	snapshotter       *service_snaphsotter.Service
}

func New(
	logger *zap.Logger,
	executorProvider executor.Provider,

	configRepo db_integrations.IntegrationsRepository,
	ciCommitRepo db_ci.CommitRepository,

	buildkiteService service_buildkite.Service,
	githubService service_github.Service,

	cfg *configuration.Configuration,
	statusService *svc_statuses.Service,
	jwtService *service_jwt.Service,
	snapshotter *service_snaphsotter.Service,
) *Service {
	return &Service{
		logger:           logger.Named("ciService"),
		executorProvider: executorProvider,

		configRepo:   configRepo,
		ciCommitRepo: ciCommitRepo,

		buildkiteService: buildkiteService,
		githubService:    githubService,

		publicApiHostname: cfg.PublicAPIHostname,
		statusService:     statusService,
		jwtService:        jwtService,
		snapshotter:       snapshotter,
	}
}

type sturdyJsonData struct {
	CodebaseID  codebases.ID `json:"codebase_id"`
	ChangeID    *string      `json:"change_id,omitempty"`
	WorkspaceID *string      `json:"workspace_id,omitempty"`
}

//go:embed download.bash
var downloadBash string

func (svc *Service) loadSeedFiles(commitSHA string, codebaseID codebases.ID, seedFiles []string) (map[string][]byte, error) {
	seedFilesContents := make(map[string][]byte)
	if err := svc.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
		for _, sf := range seedFiles {
			contents, err := repo.FileContentsAtCommit(commitSHA, sf)
			switch {
			case err == nil:
				seedFilesContents[sf] = contents
			case errors.Is(err, vcs.ErrFileNotFound):
				continue
			default:
				return fmt.Errorf("failed to read file contents: %w", err)
			}
		}
		return nil
	}).ExecTrunk(codebaseID, "readSeedFiles"); err != nil {
		return nil, fmt.Errorf("failed to get seed files: %w", err)
	}
	return seedFilesContents, nil
}

func (svc *Service) createSnapshotGit(ctx context.Context, snapshot *snapshots.Snapshot, seedFiles []string) (string, error) {
	jwt, err := svc.jwtService.IssueToken(ctx, snapshot.WorkspaceID, oneDay, jwt.TokenTypeCI)
	if err != nil {
		return "", err
	}

	// Load seed files contents from trunk
	seedFilesContents, err := svc.loadSeedFiles(snapshot.CommitSHA, snapshot.CodebaseID, seedFiles)
	if err != nil {
		return "", err
	}

	var commitSHA string
	if err := svc.executorProvider.New().
		AllowRebasingState(). // allowed because the repo might not exist yet
		Schedule(func(repoProvider provider.RepoProvider) error {
			// Create repo if not exists
			// This is a non-bare repository
			ciPath := repoProvider.ViewPath(snapshot.CodebaseID, "ci")

			var repo vcs.RepoWriter
			// Create if not exists
			if _, err := os.Open(ciPath); errors.Is(err, os.ErrNotExist) {
				repo, err = vcs.CreateNonBareRepoWithRootCommit(ciPath, "main")
				if err != nil {
					return fmt.Errorf("failed to init repo: %w", err)
				}
			} else if err != nil {
				return fmt.Errorf("failed to create repo: %w", err)
			} else {
				repo, err = repoProvider.ViewRepo(snapshot.CodebaseID, "ci")
				if err != nil {
					return fmt.Errorf("failed to init repo: %w", err)
				}
			}

			data, err := json.Marshal(sturdyJsonData{
				CodebaseID:  snapshot.CodebaseID,
				WorkspaceID: &snapshot.WorkspaceID,
			})
			if err != nil {
				return fmt.Errorf("failed to create metadata file: %w", err)
			}

			// Create commit for this change

			// Write seed files
			for sfPath, data := range seedFilesContents {
				filepath := path.Join(ciPath, sfPath)
				if err := os.MkdirAll(path.Dir(filepath), 0755); err != nil {
					return fmt.Errorf("failed to create directory for seed file (%s): %w", sfPath, err)
				}
				if err := os.WriteFile(filepath, data, 0o644); err != nil {
					return fmt.Errorf("failed to write seed file (%s): %w", sfPath, err)
				}
			}

			// Add metadata sturdy.json
			if err := os.WriteFile(path.Join(ciPath, "sturdy.json"), data, 0o644); err != nil {
				return fmt.Errorf("failed to write metadata file: %w", err)
			}

			replacer := strings.NewReplacer(
				"__PUBLIC_API__HOSTNAME__", svc.publicApiHostname,
				"__JWT__", jwt.Token,
			)

			generatedDownloadScript := replacer.Replace(downloadBash)

			// Write download script
			if err := os.WriteFile(path.Join(ciPath, "download"), []byte(generatedDownloadScript), 0o744); err != nil {
				return fmt.Errorf("failed to write download.bash: %w", err)
			}

			commitSHA, err = repo.AddAndCommit(fmt.Sprintf("Workspace %s on Sturdy", snapshot.WorkspaceID))
			if err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}

			return nil
		}).ExecView(snapshot.CodebaseID, "ci", "prepareContinuousIntegrationRepo"); err != nil {
		return "", err
	}

	// Record in ci commits repository
	if err := svc.ciCommitRepo.Create(ctx, &ci.Commit{
		ID:              uuid.NewString(),
		CodebaseID:      snapshot.CodebaseID,
		TrunkCommitSHA:  snapshot.CommitSHA,
		CiRepoCommitSHA: commitSHA,
		CreatedAt:       time.Now(),
	}); err != nil {
		return "", err
	}

	return commitSHA, nil
}

func (svc *Service) createChangeGit(ctx context.Context, ch *changes.Change, seedFiles []string) (string, error) {
	jwt, err := svc.jwtService.IssueToken(ctx, string(ch.ID), oneDay, jwt.TokenTypeCI)
	if err != nil {
		return "", err
	}

	// Load seed files contents from trunk
	seedFilesContents, err := svc.loadSeedFiles(*ch.CommitID, ch.CodebaseID, seedFiles)
	if err != nil {
		return "", err
	}

	var commitSHA string
	if err := svc.executorProvider.New().
		AllowRebasingState(). // allowed because the repo might not exist yet
		Schedule(func(repoProvider provider.RepoProvider) error {
			// Create repo if not exists
			// This is a non-bare repository
			ciPath := repoProvider.ViewPath(ch.CodebaseID, "ci")

			var repo vcs.RepoWriter
			// Create if not exists
			if _, err := os.Open(ciPath); errors.Is(err, os.ErrNotExist) {
				repo, err = vcs.CreateNonBareRepoWithRootCommit(ciPath, "main")
				if err != nil {
					return fmt.Errorf("failed to init repo: %w", err)
				}
			} else if err != nil {
				return fmt.Errorf("failed to create repo: %w", err)
			} else {
				repo, err = repoProvider.ViewRepo(ch.CodebaseID, "ci")
				if err != nil {
					return fmt.Errorf("failed to init repo: %w", err)
				}
			}
			changeID := ch.ID.String()
			data, err := json.Marshal(sturdyJsonData{
				CodebaseID: ch.CodebaseID,
				ChangeID:   &changeID,
			})
			if err != nil {
				return fmt.Errorf("failed to create metadata file: %w", err)
			}

			// Create commit for this change

			// Write seed files
			for sfPath, data := range seedFilesContents {
				filepath := path.Join(ciPath, sfPath)
				if err := os.MkdirAll(path.Dir(filepath), 0755); err != nil {
					return fmt.Errorf("failed to create directory for seed file (%s): %w", sfPath, err)
				}
				if err := os.WriteFile(filepath, data, 0o644); err != nil {
					return fmt.Errorf("failed to write seed file (%s): %w", sfPath, err)
				}
			}

			// Add metadata sturdy.json
			if err := os.WriteFile(path.Join(ciPath, "sturdy.json"), data, 0o644); err != nil {
				return fmt.Errorf("failed to write metadata file: %w", err)
			}

			replacer := strings.NewReplacer(
				"__PUBLIC_API__HOSTNAME__", svc.publicApiHostname,
				"__JWT__", jwt.Token,
			)

			generatedDownloadScript := replacer.Replace(downloadBash)

			// Write download script
			if err := os.WriteFile(path.Join(ciPath, "download"), []byte(generatedDownloadScript), 0o744); err != nil {
				return fmt.Errorf("failed to write download.bash: %w", err)
			}

			commitSHA, err = repo.AddAndCommit(fmt.Sprintf("Change %s on Sturdy", ch.ID))
			if err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}

			return nil
		}).ExecView(ch.CodebaseID, "ci", "prepareContinuousIntegrationRepo"); err != nil {
		return "", err
	}

	// Record in ci commits repository
	if err := svc.ciCommitRepo.Create(ctx, &ci.Commit{
		ID:              uuid.NewString(),
		CodebaseID:      ch.CodebaseID,
		TrunkCommitSHA:  *ch.CommitID,
		CiRepoCommitSHA: commitSHA,
		CreatedAt:       time.Now(),
	}); err != nil {
		return "", err
	}

	return commitSHA, nil
}

func (svc *Service) ListByCodebaseID(ctx context.Context, codebaseID codebases.ID) ([]*integrations.Integration, error) {
	return svc.configRepo.ListByCodebaseID(ctx, codebaseID)
}

func (svc *Service) GetByID(ctx context.Context, integrationID string) (*integrations.Integration, error) {
	return svc.configRepo.Get(ctx, integrationID)
}

func (svc *Service) Delete(ctx context.Context, integrationID string) error {
	cfg, err := svc.configRepo.Get(ctx, integrationID)
	if err != nil {
		return err
	}

	t := time.Now()
	cfg.DeletedAt = &t

	if err := svc.configRepo.Update(ctx, cfg); err != nil {
		return err
	}
	return nil
}

type TriggerOptions struct {
	// Which integrations to trigger. If empty, all integrations will be triggered.
	Providers *map[providers.ProviderName]bool
}

type TriggerOption func(*TriggerOptions)

func WithProvider(providerType providers.ProviderName) TriggerOption {
	return func(options *TriggerOptions) {
		if options.Providers == nil {
			options.Providers = &map[providers.ProviderName]bool{}
		}
		(*options.Providers)[providerType] = true
	}
}

func getTriggerOptions(opts ...TriggerOption) *TriggerOptions {
	triggerOptions := &TriggerOptions{}
	for _, opt := range opts {
		opt(triggerOptions)
	}
	return triggerOptions
}

func (svc *Service) TriggerWorkspace(ctx context.Context, workspace *workspaces.Workspace, opts ...TriggerOption) ([]*statuses.Status, error) {
	if workspace.LatestSnapshotID == nil {
		return nil, fmt.Errorf("workspace has no latest snapshot")
	}

	ciConfigurations, err := svc.configRepo.ListByCodebaseID(ctx, workspace.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list ci configs: %w", err)
	}

	snapshot, err := svc.snapshotter.GetByID(ctx, *workspace.LatestSnapshotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}
	// todo: do not mix seed files?
	seedFiles := []string{}
	for _, c := range ciConfigurations {
		seedFiles = append(seedFiles, c.SeedFiles...)
	}

	commitID, err := svc.createSnapshotGit(ctx, snapshot, seedFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to create git commit: %w", err)
	}

	options := getTriggerOptions(opts...)
	ss := []*statuses.Status{}

	for _, config := range ciConfigurations {
		if options.Providers != nil {
			if !(*options.Providers)[config.Provider] {
				continue
			}
		}

		switch config.Provider {
		case providers.ProviderNameBuildkite:
			build, err := svc.buildkiteService.CreateBuild(ctx, config.ID, commitID, workspace.NameOrFallback())
			if err != nil {
				return nil, fmt.Errorf("failed to trigger buildkite build: %w", err)
			}

			status := &statuses.Status{
				ID:          uuid.NewString(),
				CommitSHA:   snapshot.CommitSHA,
				CodebaseID:  snapshot.CodebaseID,
				Type:        statuses.TypePending,
				Title:       build.Name,
				Description: build.Description,
				DetailsURL:  &build.URL,
				Timestamp:   time.Now(),
			}

			// Set status
			if err := svc.statusService.Set(ctx, status); err != nil {
				return nil, fmt.Errorf("failed to set status: %w", err)
			}

			ss = append(ss, status)

		default:
			return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
		}
	}

	// if the codebase has a github integration, trigger ci via github as well
	// this does not set a status directly, statuses will be set when we receive webhooks from GitHub (if any)
	// TODO: better failure if the codebase has no github integration
	if options.Providers == nil || (*options.Providers)[providers.ProviderNameGithub] {
		if err := svc.githubService.CreateBuild(ctx, workspace.CodebaseID, snapshot.CommitSHA, "sturdy-ci-"+workspace.ID); err != nil {
			return nil, fmt.Errorf("failed to trigger github build workspace: %w", err)
		}
	}

	return ss, nil
}

// TriggerChange starts a continuous integration build for the given change.
func (svc *Service) TriggerChange(ctx context.Context, ch *changes.Change, opts ...TriggerOption) ([]*statuses.Status, error) {
	ciConfigurations, err := svc.configRepo.ListByCodebaseID(ctx, ch.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list ci configs: %w", err)
	}

	// todo: do not mix seed files?
	seedFiles := []string{}
	for _, c := range ciConfigurations {
		seedFiles = append(seedFiles, c.SeedFiles...)
	}

	commitID, err := svc.createChangeGit(ctx, ch, seedFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to create git commit: %w", err)
	}

	var title string
	if ch.Title != nil {
		title = *ch.Title
	} else {
		title = "Unnamed change on Sturdy"
	}

	options := getTriggerOptions(opts...)
	ss := []*statuses.Status{}
	for _, config := range ciConfigurations {
		if options.Providers != nil {
			if !(*options.Providers)[config.Provider] {
				continue
			}
		}

		switch config.Provider {
		case providers.ProviderNameBuildkite:
			build, err := svc.buildkiteService.CreateBuild(ctx, config.ID, commitID, title)
			if err != nil {
				return nil, fmt.Errorf("failed to trigger build: %w", err)
			}

			status := &statuses.Status{
				ID:          uuid.NewString(),
				CommitSHA:   *ch.CommitID,
				CodebaseID:  ch.CodebaseID,
				Type:        statuses.TypePending,
				Title:       build.Name,
				Description: build.Description,
				DetailsURL:  &build.URL,
				Timestamp:   time.Now(),
			}

			// Set status
			if err := svc.statusService.Set(ctx, status); err != nil {
				return nil, fmt.Errorf("failed to set status: %w", err)
			}

			ss = append(ss, status)
		default:
			return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
		}
	}

	return ss, nil
}

func (svc *Service) GetTrunkCommitSHA(ctx context.Context, codebaseID codebases.ID, ciRepoCommitID string) (string, error) {
	c, err := svc.ciCommitRepo.GetByCodebaseAndCiRepoCommitID(ctx, codebaseID, ciRepoCommitID)
	if err != nil {
		return "", err
	}
	return c.TrunkCommitSHA, nil
}

func (svc *Service) CreateIntegration(ctx context.Context, integration *integrations.Integration) error {
	return svc.configRepo.Create(ctx, integration)
}

func (svc *Service) UpdateIntegration(ctx context.Context, integration *integrations.Integration) error {
	return svc.configRepo.Update(ctx, integration)
}
