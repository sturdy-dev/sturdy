package vcs

import (
	"fmt"
	"os"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func (r *repository) PushRemoteUrlWithRefspec(remoteUrl string, creds transport.AuthMethod, refspecs []config.RefSpec) (userError string, err error) {
	defer getMeterFunc("PushRemoteUrlWithRefspec")()

	gg, err := gogit.PlainOpen(r.path)
	if err != nil {
		return "", err
	}

	remote, err := gg.CreateRemoteAnonymous(&config.RemoteConfig{
		Name: "anonymous",
		URLs: []string{remoteUrl},
	})
	if err != nil {
		return "", err
	}

	err = remote.Push(&gogit.PushOptions{
		RemoteName: "anonymous",
		RefSpecs:   refspecs,
		Auth:       creds,
	})

	if err != nil {
		return "", err
	}

	return "", nil
}

func (r *repository) PushNamedRemoteWithRefspec(remoteName string, creds transport.AuthMethod, refspecs []config.RefSpec) (userError string, err error) {
	defer getMeterFunc("PushNamedRemoteWithRefspec")()

	gg, err := gogit.PlainOpen(r.path)
	if err != nil {
		return "", fmt.Errorf("failed to open gogit: %w", err)
	}

	remote, err := gg.Remote(remoteName)
	if err != nil {
		return "", fmt.Errorf("failed to get remote: %w", err)
	}

	err = remote.Push(&gogit.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   refspecs,
		Auth:       creds,
	})
	if err != nil && strings.Contains(err.Error(), "protected branch hook declined") {
		return fmt.Sprintf("GitHub rejected the push as the branch is protected by branch protection rules."), fmt.Errorf("failed to push to github: stopped by branch protection rules")
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

	remote, err := gg.Remote(remoteName)
	if err != nil {
		return fmt.Errorf("failed to get remote: %w", err)
	}

	err = remote.Fetch(&gogit.FetchOptions{
		RefSpecs: refspecs,
		Auth:     creds,
		Progress: os.Stderr,
		Depth:    100,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	return nil
}

func (r *repository) FetchUrlRemoteWithCreds(remoteUrl string, creds transport.AuthMethod, refspecs []config.RefSpec) error {
	defer getMeterFunc("FetchUrlRemoteWithCreds")()

	gg, err := gogit.PlainOpen(r.path)
	if err != nil {
		return err
	}

	remote, err := gg.CreateRemoteAnonymous(&config.RemoteConfig{
		Name: "anonymous",
		URLs: []string{remoteUrl},
	})
	if err != nil {
		return err
	}

	err = remote.Fetch(&gogit.FetchOptions{
		RefSpecs: refspecs,
		Auth:     creds,
		Progress: os.Stderr,
		Depth:    100,
	})
	if err != nil {
		return err
	}

	return nil
}
