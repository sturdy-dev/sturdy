package enterprise

import (
	"context"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/enterprise/db"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/users"

	"github.com/graph-gophers/graphql-go"
)

type GitHubAccountRootResolver struct {
	gitHubUserRepo db.GitHubUserRepo
}

func NewGitHubAccountRootResolver(
	gitHubUserRepo db.GitHubUserRepo,
) *GitHubAccountRootResolver {
	return &GitHubAccountRootResolver{
		gitHubUserRepo: gitHubUserRepo,
	}
}

func (r *GitHubAccountRootResolver) InteralByID(_ context.Context, id users.ID) (resolvers.GitHubAccountResolver, error) {
	githubUser, err := r.gitHubUserRepo.GetByUserID(id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &gitHubAccountResolver{
		githubUser: githubUser,
	}, nil
}

type gitHubAccountResolver struct {
	githubUser *github.GitHubUser
}

func (r *gitHubAccountResolver) ID() graphql.ID {
	return graphql.ID(r.githubUser.ID)
}

func (r *gitHubAccountResolver) Login() string {
	return r.githubUser.Username
}
