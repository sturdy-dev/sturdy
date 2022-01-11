package db

import "mash/pkg/di"

var Module = di.NewModule(
	di.Provides(NewGitHubInstallationRepo),
	di.Provides(NewGitHubPRRepo),
	di.Provides(NewGitHubRepositoryRepo),
	di.Provides(NewGitHubUserRepo),
)
