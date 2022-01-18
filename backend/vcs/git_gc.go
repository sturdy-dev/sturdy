package vcs

import (
	"bytes"
	"fmt"
	"os/exec"
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
