package vcs

import (
	"fmt"
	"os"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"

	gh "github.com/google/go-github/v39/github"
	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

func Create(codebaseID codebases.ID) func(provider.RepoProvider) error {
	return func(trunkProvider provider.RepoProvider) error {
		path := trunkProvider.TrunkPath(codebaseID)
		if _, err := vcs.CreateBareRepoWithRootCommit(path); err != nil {
			return fmt.Errorf("failed to create trunk: %w", err)
		}
		return nil
	}
}

// If no limit is set, a default of 100 is used
// TODO(gustav): delete
func ListChanges(repo vcs.RepoGitReader, limit int) ([]*vcs.LogEntry, error) {
	if limit < 1 {
		limit = 100
	}
	changeLog, err := repo.LogHead(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get changes: %w", err)
	}

	// TODO: Do we really need this?
	var filteredLog []*vcs.LogEntry
	for _, e := range changeLog {
		if e.RawCommitMessage == "Root Commit" {
			continue
		}
		filteredLog = append(filteredLog, e)
	}

	return filteredLog, nil
}

func CloneFromGithub(logger *zap.Logger, trunkProvider provider.TrunkProvider, codebaseID codebases.ID, repo *gh.Repository, accessToken string) error {
	barePath := trunkProvider.TrunkPath(codebaseID)
	if _, err := os.Open(barePath); err != nil && os.IsNotExist(err) {
		upstream := repo.GetCloneURL()
		logger.Info("cloning from github", zap.String("upstream", upstream))
		_, err := vcs.RemoteCloneWithCreds(
			upstream,
			barePath,
			newCredentialsCallback("x-access-token", accessToken),
			true,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func newCredentialsCallback(tokenUsername, token string) git.CredentialsCallback {
	return func(url string, username string, allowedTypes git.CredType) (*git.Cred, error) {
		cred, _ := git.NewCredUserpassPlaintext(tokenUsername, token)
		return cred, nil
	}
}
