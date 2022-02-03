package provider

import "getsturdy.com/api/pkg/configuration/flags"

type Configuration struct {
	ReposPath string               `long:"repos-path" description:"Path to the directory containing the repositories" required:"true" default:"tmp/repos"`
	LFS       *gitLFSConfiguration `flags-group:"git-lfs" namespace:"lfs"`
}

type gitLFSConfiguration struct {
	Addr flags.Addr `long:"addr" description:"Git LFS server address" required:"true" default:"localhost:8080"`
}

func FromConfigration(cfg *Configuration) RepoProvider {
	return New(cfg.ReposPath, cfg.LFS.Addr.String())
}
