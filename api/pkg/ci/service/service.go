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

	integrations "getsturdy.com/api/pkg/integrations"

	"getsturdy.com/api/pkg/change"
	db_change "getsturdy.com/api/pkg/change/db"
	"getsturdy.com/api/pkg/ci"
	db_ci "getsturdy.com/api/pkg/ci/db"
	db_integrations "getsturdy.com/api/pkg/integrations/db"
	"getsturdy.com/api/pkg/jwt"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/statuses"
	svc_statuses "getsturdy.com/api/pkg/statuses/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	oneDay = 24 * time.Hour
)

type Service struct {
	logger           *zap.Logger
	executorProvider executor.Provider

	configRepo       db_integrations.IntegrationsRepository
	ciCommitRepo     db_ci.CommitRepository
	changeRepo       db_change.Repository
	changeCommitRepo db_change.CommitRepository

	publicApiHostname string
	statusService     *svc_statuses.Service
	jwtService        *service_jwt.Service
}

type Configuration struct {
	PublicAPIHostname string `long:"public-api-hostname" description:"Public API hostname. Used to fetch codebases from CI"`
}

func New(
	logger *zap.Logger,
	executorProvider executor.Provider,

	configRepo db_integrations.IntegrationsRepository,
	ciCommitRepo db_ci.CommitRepository,
	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,

	cfg *Configuration,
	statusService *svc_statuses.Service,
	jwtService *service_jwt.Service,
) *Service {
	return &Service{
		logger:           logger.Named("ciService"),
		executorProvider: executorProvider,

		configRepo:       configRepo,
		ciCommitRepo:     ciCommitRepo,
		changeRepo:       changeRepo,
		changeCommitRepo: changeCommitRepo,

		publicApiHostname: cfg.PublicAPIHostname,
		statusService:     statusService,
		jwtService:        jwtService,
	}
}

type sturdyJsonData struct {
	CodebaseID string `json:"codebase_id"`
	ChangeID   string `json:"change_id"`
}

//go:embed download.bash
var downloadBash string

func (svc *Service) loadSeedFiles(commit *change.ChangeCommit, seedFiles []string) (map[string][]byte, error) {
	seedFilesContents := make(map[string][]byte)
	if err := svc.executorProvider.New().Git(func(repo vcs.Repo) error {
		for _, sf := range seedFiles {
			contents, err := repo.FileContentsAtCommit(commit.CommitID, sf)
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
	}).ExecTrunk(commit.CodebaseID, "readSeedFiles"); err != nil {
		return nil, fmt.Errorf("failed to get seed files: %w", err)
	}
	return seedFilesContents, nil
}

func (svc *Service) createGit(ctx context.Context, commit *change.ChangeCommit, seedFiles []string) (string, error) {
	jwt, err := svc.jwtService.IssueToken(ctx, string(commit.ChangeID), oneDay, jwt.TokenTypeCI)
	if err != nil {
		return "", err
	}

	// Load seed files contents from trunk
	seedFilesContents, err := svc.loadSeedFiles(commit, seedFiles)
	if err != nil {
		return "", err
	}

	var commitID string
	if err := svc.executorProvider.New().
		AllowRebasingState(). // allowed because the repo might not exist yet
		Schedule(func(repoProvider provider.RepoProvider) error {
			// Create repo if not exists
			// This is a non-bare repository
			ciPath := repoProvider.ViewPath(commit.CodebaseID, "ci")

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
				repo, err = repoProvider.ViewRepo(commit.CodebaseID, "ci")
				if err != nil {
					return fmt.Errorf("failed to init repo: %w", err)
				}
			}

			data, err := json.Marshal(sturdyJsonData{
				CodebaseID: commit.CodebaseID,
				ChangeID:   string(commit.ChangeID),
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

			commitID, err = repo.AddAndCommit(fmt.Sprintf("Change %s on Sturdy", commit.ChangeID))
			if err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}

			return nil
		}).ExecView(commit.CodebaseID, "ci", "prepareContinuousIntegrationRepo"); err != nil {
		return "", err
	}

	// Record in ci commits repository
	if err := svc.ciCommitRepo.Create(ctx, &ci.Commit{
		ID:             uuid.NewString(),
		CodebaseID:     commit.CodebaseID,
		TrunkCommitID:  commit.CommitID,
		CiRepoCommitID: commitID,
		CreatedAt:      time.Now(),
	}); err != nil {
		return "", err
	}

	return commitID, nil
}

func (svc *Service) ListByCodebaseID(ctx context.Context, codebaseID string) ([]*integrations.Integration, error) {
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
	Providers *map[integrations.ProviderType]bool
}

type TriggerOption func(*TriggerOptions)

func WithProvider(providerType integrations.ProviderType) TriggerOption {
	return func(options *TriggerOptions) {
		if options.Providers == nil {
			options.Providers = &map[integrations.ProviderType]bool{}
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

// Trigger starts a contihuous integration build for the given change.
func (svc *Service) Trigger(ctx context.Context, ch *change.Change, opts ...TriggerOption) ([]*statuses.Status, error) {
	ciConfigurations, err := svc.configRepo.ListByCodebaseID(ctx, ch.CodebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list ci configs: %w", err)
	}

	commit, err := svc.changeCommitRepo.GetByChangeIDOnTrunk(ch.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get change commit: %w", err)
	}

	// todo: do not mix seed files?
	seedFiles := []string{}
	for _, c := range ciConfigurations {
		seedFiles = append(seedFiles, c.SeedFiles...)
	}

	commitID, err := svc.createGit(ctx, &commit, seedFiles)
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
	for _, configuration := range ciConfigurations {
		if options.Providers != nil {
			if !(*options.Providers)[configuration.Provider] {
				continue
			}
		}

		provider, err := integrations.Get(configuration.Provider)
		if err != nil {
			return nil, fmt.Errorf("failed to get provider: %w", err)
		}

		build, err := provider.CreateBuild(ctx, configuration.ID, commitID, title)
		if err != nil {
			return nil, fmt.Errorf("failed to trigger build: %w", err)
		}

		status := &statuses.Status{
			ID:          uuid.NewString(),
			CommitID:    commit.CommitID,
			CodebaseID:  ch.CodebaseID,
			Type:        statuses.TypePending,
			Title:       build.Name,
			Description: &build.Description,
			DetailsURL:  &build.URL,
			Timestamp:   time.Now(),
		}

		// Set status
		if err := svc.statusService.Set(ctx, status); err != nil {
			return nil, fmt.Errorf("failed to set status: %w", err)
		}

		ss = append(ss, status)
	}

	return ss, nil
}

func (svc *Service) GetTrunkCommitID(ctx context.Context, codebaseID, ciRepoCommitID string) (string, error) {
	c, err := svc.ciCommitRepo.GetByCodebaseAndCiRepoCommitID(ctx, codebaseID, ciRepoCommitID)
	if err != nil {
		return "", err
	}
	return c.TrunkCommitID, nil
}

func (svc *Service) CreateIntegration(ctx context.Context, integration *integrations.Integration) error {
	return svc.configRepo.Create(ctx, integration)
}

func (svc *Service) UpdateIntegration(ctx context.Context, integration *integrations.Integration) error {
	return svc.configRepo.Update(ctx, integration)
}
