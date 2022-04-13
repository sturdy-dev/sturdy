package provider

import (
	"path"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider/configuration"
)

// RepoProvider
// Deprecated: RepoProvider in favour for executor.Executor (except when used in an executor.ExecuteFunc)
type RepoProvider interface {
	TrunkProvider
	ViewProvider
}

// TrunkProvider
// Deprecated: TrunkProvider in favour for executor.Executor (except when used in an executor.ExecuteFunc)
type TrunkProvider interface {
	TrunkRepo(codebaseID codebases.ID) (vcs.RepoWriter, error)
	TrunkPath(codebaseID codebases.ID) string
}

// ViewProvider
// Deprecated: ViewProvider in favour for executor.Executor (except when used in an executor.ExecuteFunc)
type ViewProvider interface {
	ViewRepo(codebaseID codebases.ID, viewID string) (vcs.RepoWriter, error)
	ViewPath(codebaseID codebases.ID, viewID string) string
}

type repoProvider struct {
	reposBasePath string
	lfsHostname   string
}

func New(reposBasePath, lfsHostname string) RepoProvider {
	return &repoProvider{
		reposBasePath: reposBasePath,
		lfsHostname:   lfsHostname,
	}
}

func (r *repoProvider) TrunkRepo(codebaseID codebases.ID) (vcs.RepoWriter, error) {
	return vcs.OpenRepoWithLFS(path.Join(r.reposBasePath, codebaseID.String(), "trunk"), r.lfsHostname)
}

func (r *repoProvider) ViewRepo(codebaseID codebases.ID, viewID string) (vcs.RepoWriter, error) {
	return vcs.OpenRepoWithLFS(path.Join(r.reposBasePath, codebaseID.String(), viewID), r.lfsHostname)
}

func (r *repoProvider) TrunkPath(codebaseID codebases.ID) string {
	return path.Join(r.reposBasePath, codebaseID.String(), "trunk")
}

func (r *repoProvider) ViewPath(codebaseID codebases.ID, viewID string) string {
	return path.Join(r.reposBasePath, codebaseID.String(), viewID)
}

func FromConfiguration(cfg *configuration.Configuration) RepoProvider {
	return New(cfg.ReposPath, cfg.LFS.Addr.String())
}
