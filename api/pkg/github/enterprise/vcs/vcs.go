package vcs

import (
	"fmt"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
)

func FetchTrackedToSturdytrunk(accessToken, ref string) func(vcs.RepoGitWriter) error {
	return func(repo vcs.RepoGitWriter) error {
		refspec := fmt.Sprintf("+%s:refs/heads/sturdytrunk", ref)
		if err := repo.FetchNamedRemoteWithCreds("origin", newCredentialsCallback(accessToken), []config.RefSpec{config.RefSpec(refspec)}); err != nil {
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

func FetchBranchWithRefspec(accessToken, refspec string) func(vcs.RepoGitWriter) error {
	return func(repo vcs.RepoGitWriter) error {
		if err := repo.FetchNamedRemoteWithCreds("origin", newCredentialsCallback(accessToken), []config.RefSpec{config.RefSpec(refspec)}); err != nil {
			return fmt.Errorf("failed to perform remote fetch: %w", err)
		}
		return nil
	}
}

func PushTrackedToGitHub(repo vcs.RepoGitWriter, accessToken, trackedBranchName string) (userError string, err error) {
	refspec := fmt.Sprintf("+refs/heads/sturdytrunk:refs/heads/%s", trackedBranchName)
	return PushToGitHubWithRefspec(repo, accessToken, refspec)
}

func PushToGitHubWithRefspec(repo vcs.RepoGitWriter, accessToken, refspec string) (userError string, err error) {
	userError, err = repo.PushNamedRemoteWithRefspec("origin", newCredentialsCallback(accessToken), []config.RefSpec{config.RefSpec(refspec)})
	if err != nil {
		return userError, fmt.Errorf("failed to push %s: %w", refspec, err)
	}
	return "", nil
}

func PushBranchToGithubWithForce(executorProvider executor.Provider, codebaseID codebases.ID, sturdyBranchName, remoteBranchName, accessToken string) (userError string, err error) {
	refspec := fmt.Sprintf("+refs/heads/%s:refs/heads/%s", sturdyBranchName, remoteBranchName)

	err = executorProvider.New().GitWrite(func(r vcs.RepoGitWriter) error {
		userError, err = r.PushNamedRemoteWithRefspec("origin", newCredentialsCallback(accessToken), []config.RefSpec{config.RefSpec(refspec)})
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

func PushBranchToGithubSafely(executorProvider executor.Provider, codebaseID codebases.ID, sturdyBranchName, remoteBranchName, accessToken string) (userError string, err error) {
	refspec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", sturdyBranchName, remoteBranchName)

	err = executorProvider.New().GitWrite(func(r vcs.RepoGitWriter) error {
		userError, err = r.PushNamedRemoteWithRefspec("origin", newCredentialsCallback(accessToken), []config.RefSpec{config.RefSpec(refspec)})
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

func HaveTrackedBranch(executorProvider executor.Provider, codebaseID codebases.ID, remoteBranchName string) error {
	err := executorProvider.New().GitRead(func(r vcs.RepoGitReader) error {
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

func newCredentialsCallback(token string) transport.AuthMethod {
	return &http.BasicAuth{
		Username: "x-access-token",
		Password: token,
	}
}
