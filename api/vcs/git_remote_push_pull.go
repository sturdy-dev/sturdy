package vcs

import (
	"errors"
	"fmt"
	"os"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	git "github.com/libgit2/git2go/v33"
)

func (r *repository) PushRemoteUrlWithRefspec(remoteUrl string, creds transport.AuthMethod, refspecs []config.RefSpec) (userError string, err error) {
	defer getMeterFunc("PushRemoteUrlWithRefspec")()

	gg, err := gogit.PlainOpen(r.path)
	if err != nil {
		return "", fmt.Errorf("failed to open gogit: %w", err)
	}

	remote, err := gg.CreateRemoteAnonymous(&config.RemoteConfig{
		Name: "anonymous",
		URLs: []string{remoteUrl},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create remote: %w", err)
	}

	err = remote.Push(&gogit.PushOptions{
		RemoteName: "anonymous",
		RefSpecs:   refspecs,
		Auth:       creds,
		Force:      true,
	})
	if errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to push: %w", err)
	}
	return "", nil
}

func (r *repository) PushNamedRemoteWithRefspec(remoteName string, creds transport.AuthMethod, refspecs []config.RefSpec) (userError string, err error) {
	defer getMeterFunc("PushNamedRemoteWithRefspec")()

	gg, err := gogit.PlainOpen(r.path)
	if err != nil {
		return "", fmt.Errorf("failed to open gogit: %w", err)
	}

	err = gg.Push(&gogit.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   refspecs,
		Auth:       creds,
		Force:      true,
	})
	if err != nil && strings.Contains(err.Error(), "protected branch hook declined") {
		return fmt.Sprintf("GitHub rejected the push as the branch is protected by branch protection rules."), fmt.Errorf("failed to push to github: stopped by branch protection rules")
	}
	if errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to push: %w", err)
	}
	return "", nil
}

func (r *repository) FetchNamedRemoteWithCreds(remoteName string, creds transport.AuthMethod, refspecs []config.RefSpec) error {
	defer getMeterFunc("FetchNamedRemoteWithCreds")()

	gg, err := gogit.PlainOpen(r.path)
	if err != nil {
		return fmt.Errorf("failed to open gogit: %w", err)
	}

	err = gg.Fetch(&gogit.FetchOptions{
		RemoteName: remoteName,
		RefSpecs:   refspecs,
		Auth:       creds,
		Progress:   os.Stderr,
		Depth:      100,
		Force:      true,
	})
	if errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	return nil
}

func (r *repository) FetchUrlRemoteWithCreds(remoteUrl string, creds transport.AuthMethod, refspecs []config.RefSpec) error {
	defer getMeterFunc("FetchUrlRemoteWithCreds")()

	gg, err := gogit.PlainOpen(r.path)
	if err != nil {
		return fmt.Errorf("failed to open gogit: %w", err)
	}

	remote, err := gg.CreateRemoteAnonymous(&config.RemoteConfig{
		Name: "anonymous",
		URLs: []string{remoteUrl},
	})
	if err != nil {
		return fmt.Errorf("failed to create remote: %w", err)
	}

	err = remote.Fetch(&gogit.FetchOptions{
		RefSpecs: refspecs,
		Auth:     creds,
		Progress: os.Stderr,
		Depth:    100,
		Force:    true,
	})
	if errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	return nil
}

func (r *repository) AddNamedRemote(name, url string) error {
	_, err := r.r.Remotes.Create(name, url)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateRef(name, commitSha string) error {
	id, err := git.NewOid(commitSha)
	if err != nil {
		return err
	}
	_, err = r.r.References.Create(name, id, true, "create-ref-"+name)
	if err != nil {
		return err
	}
	return nil
}
