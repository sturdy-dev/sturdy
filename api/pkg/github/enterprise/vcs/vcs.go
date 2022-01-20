package vcs

import (
	"fmt"

	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

func FetchTrackedToSturdytrunk(accessToken, ref string) func(vcs.Repo) error {
	return func(repo vcs.Repo) error {
		refspec := fmt.Sprintf("+%s:refs/heads/sturdytrunk", ref)
		if err := repo.RemoteFetchWithCreds("origin", newCredentialsCallback(accessToken), []string{refspec}); err != nil {
			return fmt.Errorf("failed to perform remote fetch: %w", err)
		}

		// Make sure that sturdytrunk is the HEAD branch
		// This is the case for repositories that where empty the first time they where cloned to Sturdy
		if err := repo.SetDefaultBranch("sturdytrunk"); err != nil {
			return fmt.Errorf("could not set default branch: %w", err)
		}
		return nil
	}
}

func FetchBranchWithRefspec(accessToken, refspec string) func(vcs.Repo) error {
	return func(repo vcs.Repo) error {
		if err := repo.RemoteFetchWithCreds("origin", newCredentialsCallback(accessToken), []string{refspec}); err != nil {
			return fmt.Errorf("failed to perform remote fetch: %w", err)
		}
		return nil
	}
}

func PushTrackedToGitHub(logger *zap.Logger, repo vcs.Repo, accessToken, trackedBranchName string) (userError string, err error) {
	refspec := fmt.Sprintf("+refs/heads/sturdytrunk:refs/heads/%s", trackedBranchName)
	userError, err = repo.PushNamedRemoteWithRefspec(logger, "origin", newCredentialsCallback(accessToken), []string{refspec})
	if err != nil {
		return userError, fmt.Errorf("failed to push %s: %w", refspec, err)
	}
	return "", nil
}

func PushBranchToGithubWithForce(logger *zap.Logger, executorProvider executor.Provider, codebaseID, sturdyBranchName, remoteBranchName, accessToken string) (userError string, err error) {
	refspec := fmt.Sprintf("+refs/heads/%s:refs/heads/%s", sturdyBranchName, remoteBranchName)

	err = executorProvider.New().Git(func(r vcs.Repo) error {
		userError, err = r.PushNamedRemoteWithRefspec(logger, "origin", newCredentialsCallback(accessToken), []string{refspec})
		if err != nil {
			return fmt.Errorf("failed to push %s: %w", refspec, err)
		}
		return nil
	}).ExecTrunk(codebaseID, "PushBranchToGithubWithForce")
	if err != nil {
		return userError, err
	}
	return userError, nil
}

func PushBranchToGithubSafely(logger *zap.Logger, executorProvider executor.Provider, codebaseID, sturdyBranchName, remoteBranchName, accessToken string) (userError string, err error) {
	refspec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", sturdyBranchName, remoteBranchName)

	err = executorProvider.New().Git(func(r vcs.Repo) error {
		userError, err = r.PushNamedRemoteWithRefspec(logger, "origin", newCredentialsCallback(accessToken), []string{refspec})
		if err != nil {
			return fmt.Errorf("failed to push %s: %w", refspec, err)
		}
		return nil
	}).ExecTrunk(codebaseID, "PushBranchToGithubSafely")
	if err != nil {
		return userError, err
	}
	return userError, nil
}

func HaveTrackedBranch(executorProvider executor.Provider, codebaseID, remoteBranchName string) error {
	err := executorProvider.New().Git(func(r vcs.Repo) error {
		_, err := r.RemoteBranchCommit("origin", remoteBranchName)
		if err != nil {
			return fmt.Errorf("could not get remote branch: %w", err)
		}
		return nil
	}).ExecTrunk(codebaseID, "haveTrackedBranch")
	if err != nil {
		return err
	}
	return nil
}

func newCredentialsCallback(token string) git.CredentialsCallback {
	return func(url string, username string, allowedTypes git.CredType) (*git.Cred, error) {
		cred, _ := git.NewCredUserpassPlaintext("x-access-token", token)
		return cred, nil
	}
}

func ListImportedChanges(repo vcs.Repo) ([]*vcs.LogEntry, error) {
	entries, err := repo.LogBranch("sturdytrunk", 50)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
