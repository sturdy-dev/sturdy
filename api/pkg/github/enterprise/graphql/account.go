package graphql

import (
	"context"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/db"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/users"
	"github.com/graph-gophers/graphql-go"
)

type GitHubAccountRootResolver struct {
	gitHubUserRepo db.GitHubUserRepo
	personalClient client.PersonalClientProvider
}

func NewGitHubAccountRootResolver(
	gitHubUserRepo db.GitHubUserRepo,
	personalClient client.PersonalClientProvider,
) *GitHubAccountRootResolver {
	return &GitHubAccountRootResolver{
		gitHubUserRepo: gitHubUserRepo,
		personalClient: personalClient,
	}
}

func (r *GitHubAccountRootResolver) InteralByID(_ context.Context, id users.ID) (resolvers.GitHubAccountResolver, error) {
	githubUser, err := r.gitHubUserRepo.GetByUserID(id)

	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &gitHubAccountResolver{
		githubUser:     githubUser,
		personalClient: r.personalClient,
	}, nil
}

type gitHubAccountResolver struct {
	githubUser     *github.GitHubUser
	personalClient client.PersonalClientProvider
}

func (r *gitHubAccountResolver) ID() graphql.ID {
	return graphql.ID(r.githubUser.ID)
}

func (r *gitHubAccountResolver) Login() string {
	return r.githubUser.Username
}

func (r *gitHubAccountResolver) IsValid(ctx context.Context) bool {
	personalClient, err := r.personalClient(r.githubUser.AccessToken)
	if err != nil {
		return false
	}

	_, _, err = personalClient.Users.Get(ctx, "")

	return err == nil
}
