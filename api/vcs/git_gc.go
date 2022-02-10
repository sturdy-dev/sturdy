package vcs

import (
	"bytes"
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

func (r *repository) GitGC() error {
	cmd := exec.Command(
		"git",
		"gc",
		"--prune=2d",
		"--force",
	)
	errLog := &bytes.Buffer{}
	cmd.Dir = r.path
	cmd.Stderr = errLog
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git-gc: %w, %s", err, errLog.String())
	}
	return nil
}

func (r *repository) GitReflogExpire() error {
	cmd := exec.Command(
		"git",
		"reflog",
		"expire",
		"--expire-unreachable=now",
		"--expire=2d",
		"--all",
	)
	errLog := &bytes.Buffer{}
	cmd.Dir = r.path
	cmd.Stderr = errLog
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git-gc: %w, %s", err, errLog.String())
	}
	return nil
}

func (r *repository) GitRemotePrune(logger *zap.Logger, remoteName string) error {
	remote, err := r.r.Remotes.Lookup(remoteName)
	if err != nil {
		return err
	}

	log := logger.Named("GitRemotePrune")

	preRefspecs, err := remote.PushRefspecs()
	if err != nil {
		log.Error("failed to get push refspects pre", zap.Error(err))
		// don't fail
		return nil
	}

	remote.RefspecCount()

	if err := remote.Prune(nil); err != nil {
		log.Error("pruning failed", zap.Error(err))
		// don't fail
		return nil
	}

	postRefspecs, err := remote.PushRefspecs()
	if err != nil {
		log.Error("failed to get push refspects post", zap.Error(err))
		// don't fail
		return nil
	}

	log.Info("cleanup remote refspecs", zap.Int("pre", len(preRefspecs)), zap.Int("post", len(postRefspecs)))

	return nil
}
