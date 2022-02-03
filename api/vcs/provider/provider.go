package provider

import (
	"path"

	"getsturdy.com/api/vcs"
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
	TrunkRepo(codebaseID string) (vcs.RepoWriter, error)
	TrunkPath(codebaseID string) string
}

// ViewProvider
// Deprecated: ViewProvider in favour for executor.Executor (except when used in an executor.ExecuteFunc)
type ViewProvider interface {
	ViewRepo(codebaseID, viewID string) (vcs.RepoWriter, error)
	ViewPath(codebaseID, viewID string) string
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

func (r *repoProvider) TrunkRepo(codebaseID string) (vcs.RepoWriter, error) {
	return vcs.OpenRepoWithLFS(path.Join(r.reposBasePath, codebaseID, "trunk"), r.lfsHostname)
}

func (r *repoProvider) ViewRepo(codebaseID, viewID string) (vcs.RepoWriter, error) {
	return vcs.OpenRepoWithLFS(path.Join(r.reposBasePath, codebaseID, viewID), r.lfsHostname)
}

func (r *repoProvider) TrunkPath(codebaseID string) string {
	return path.Join(r.reposBasePath, codebaseID, "trunk")
}

func (r *repoProvider) ViewPath(codebaseID, viewID string) string {
	return path.Join(r.reposBasePath, codebaseID, viewID)
}
